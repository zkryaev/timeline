package envars

import (
	"os"
	"strings"
)

// get path from home dir
// example: /home/$USER/timeline/$GIVEN_PATH
func GetPathByEnv(envName string) string {
	pathToFile := os.Getenv(envName)
	if pathToFile == "" {
		return ""
	}
	currDir, _ := os.Getwd()
	projectRootDir := strings.SplitAfter(currDir, "timeline/")[0]
	path := strings.Builder{}
	path.Grow(len(projectRootDir) + len(pathToFile))
	path.WriteString(projectRootDir)
	if strings.Contains(pathToFile, path.String()) {
		return pathToFile
	}
	if !strings.HasPrefix(pathToFile, "/") {
		path.WriteString("/")
	}
	path.WriteString(pathToFile)
	return path.String()
}
