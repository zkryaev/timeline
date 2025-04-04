package loader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"timeline/internal/libs/envars"
	"timeline/internal/libs/fsop"

	"go.uber.org/zap"
)

type source struct {
	Name     string
	Filepath string
	Ref      string
}

func DataSourceLoader(logger *zap.Logger) error {
	sourceEnvList := loadSourceEnvList()
	sources := parseSourceEnvList(sourceEnvList)

	// есть ли файл - да, читаем из него
	// нет - пытаемся скачать по ссылке
	for _, src := range sources {
		_, err := os.Stat(src.Filepath)
		switch {
		case err == nil:

		case src.Ref != "":
			if err := loadFromRef(src); 
		default:
			logger.Warn(fmt.Sprintf("data source invalid: name=%s filepath=%s ref=%s", src.Name, src.Filepath, src.Ref))
		}
	}

	//jsoniter.ConfigFastest.NewDecoder().Decode()
}

func loadSourceEnvList() map[string]struct{} {
	srcList := os.Getenv("SRC_LIST")
	envList := strings.Split(srcList, " ")
	srcEnvList := make(map[string]struct{}, len(envList))
	for _, env := range envList {
		srcEnvList[env] = struct{}{}
	}
	return srcEnvList
}

// example: env=ref filepath
func parseSourceEnvList(sourceEnvList map[string]struct{}) []source {
	dataSrcList := make([]source, len(sourceEnvList))
	for env := range sourceEnvList {
		line := os.Getenv(env)
		parts := strings.Split(line, " ")
		filepath := envars.GetPathByEnv(parts[0])
		dataSrcList = append(dataSrcList, source{Name: env, Filepath: filepath, Ref: parts[1]})
	}
	return dataSrcList
}

func loadFromRef(src source) error {
	resp, err := http.Get(src.Ref)
	defer resp.Body.Close()
	switch {
	case err != nil:
		return fmt.Errorf("http get request failed: %w", err)
	case resp.StatusCode != http.StatusOK:
		return fmt.Errorf("wrong http get status: %s", resp.Status)
	}
	linkParts := strings.Split(src.Ref, "/")
	filename := linkParts[len(linkParts)-1]
	filepath, err := fsop.CreateDirAndFile("data/"+filename, false)
	if err != nil {
		return fmt.Errorf("failed to create dir and file: %w", err)
	}
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}

func SaveToDB() error {}
