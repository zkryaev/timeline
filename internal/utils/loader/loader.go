package loader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"timeline/internal/infrastructure"
	"timeline/internal/utils/envars"
	"timeline/internal/utils/fsop"
	"timeline/internal/utils/loader/objects"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type source struct {
	Name     string
	Filepath string
	Ref      string
}

type BackData struct {
	Cities objects.Cities
}

func LoadData(logger *zap.Logger, db infrastructure.BackgroundDataStore, storage *BackData) error {
	envs := loadSourceEnvsList()
	sources := parseSourceEnvsList(envs)
	logger.Info("The following sources have been fetched", zap.Strings("sources", envs))
	for i, src := range sources {
		logger.Info(fmt.Sprintf("%d. %s: filepath=%q ref=%q", i+1, src.Name, src.Filepath, src.Ref))

		if err := loadDataFromSource(storage, logger, src); err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("✓ %s has been loaded", src.Name))
	}
	logger.Info("Saving background data to DB...")
	if err := saveToDB(db, logger, storage); err != nil {
		return err
	}
	logger.Info("Successfully saved to DB")
	return nil
}

func loadSourceEnvsList() []string {
	srcList := os.Getenv("SRC_LIST")
	envList := strings.Split(srcList, " ")
	return envList
}

// example: env=ref filepath
func parseSourceEnvsList(envs []string) []*source {
	dataSrcList := make([]*source, 0, len(envs))
	for _, env := range envs {
		line := os.Getenv(env)
		parts := strings.Split(line, " ")
		filepath := envars.GetPathFromProjectDir(parts[0])
		src := &source{Name: env, Filepath: filepath, Ref: parts[1]}
		dataSrcList = append(dataSrcList, src)
	}
	return dataSrcList
}

func loadDataFromSource(store *BackData, logger *zap.Logger, src *source) error {
	visited := false
	for {
		if _, err := os.Stat(src.Filepath); err != nil {
			if visited {
				return fmt.Errorf("%s: %w", src.Filepath, err)
			}
			if err := loadFromRef(src); err != nil {
				return fmt.Errorf("%s: %w", "loadFromRef", err)
			}
			logger.Info(fmt.Sprintf("* Downloaded from the link and saved in file: %s", src.Filepath))
			visited = true
			continue
		}
		switch src.Name {
		case "CITIES_SRC":
			store.loadCities(src.Filepath)
			logger.Info("* Read from the file")
		default:
			return fmt.Errorf("%s", "undefined source name")
		}
		return nil
	}
}

func saveToDB(db infrastructure.BackgroundDataStore, logger *zap.Logger, storage *BackData) error {
	ctx := context.Background()
	if storage.Cities.Arr != nil {
		err := db.SaveCities(ctx, logger, storage.Cities)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFromRef(src *source) error {
	resp, err := http.Get(src.Ref)
	switch {
	case err != nil:
		return fmt.Errorf("http get request failed: %w", err)
	case resp == nil:
		return fmt.Errorf("http get request failed: *http.Response=nil")
	case resp.StatusCode != http.StatusOK:
		return fmt.Errorf("wrong http get status: %s", resp.Status)
	}
	defer resp.Body.Close()
	linkParts := strings.Split(src.Ref, "/")
	filename := linkParts[len(linkParts)-1]
	filepath, err := fsop.CreateDirAndFile("data/"+filename, false)
	if err != nil {
		return fmt.Errorf("failed to create dir and file: %w", err)
	}
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}
	src.Filepath = filepath
	return nil
}

func (d *BackData) loadCities(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	cities := make([]objects.City, 0, 1100)
	if err := jsoniter.ConfigFastest.NewDecoder(file).Decode(&cities); err != nil {
		return err
	}
	d.Cities = objects.New(cities)
	return nil
}
