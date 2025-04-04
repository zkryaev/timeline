package envars

import (
	"os"
	"strings"
)

// Retrieve path from env and add to path project work dir
//
//	Example:
//	in:  $ENV_PATH
//	out: /home/$USER/timeline/$ENV_PATH
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
