package fsop

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

// Create dir and file with names in givenPath
// If isTimestamp true, filename will contain current timestamp.
// Returns filepath
//
//	Example:
//	 in: folder/file -> DirName=folder, FileName=file
//	 out: home/user/project/folder/file
func CreateDirAndFile(givenPath string, isTimestamp bool) (string, error) {
	var timestamp string
	if isTimestamp {
		timestamp = time.Now().Format("15:04:05_2006-01-02_")
	}
	pathparts := strings.SplitAfter(givenPath, "/")
	fileName := pathparts[len(strings.SplitAfter(givenPath, "/"))-1]
	dirName := pathparts[len(strings.SplitAfter(givenPath, "/"))-2]
	pathDir := strings.TrimSuffix(givenPath, fileName)
	if _, err := os.Stat(pathDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(pathDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("couldn't create %s dir: %s", dirName, err.Error())
		}
	}
	filepath := pathDir + timestamp + fileName
	if _, err := os.Create(filepath); err != nil && errors.Is(err, fs.ErrExist) {
		return "", fmt.Errorf("couldn't create %s file: %s", fileName, err.Error())
	}
	return filepath, nil
}
