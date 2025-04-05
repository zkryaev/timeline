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
	return GetPathFromProjectDir(pathToFile)
}

func GetPathFromProjectDir(givenPath string) string {
	currDir, _ := os.Getwd()
	projectRootDir := strings.SplitAfter(currDir, "timeline/")[0]
	path := strings.Builder{}
	path.Grow(len(projectRootDir) + len(givenPath))
	path.WriteString(projectRootDir)
	if strings.Contains(givenPath, path.String()) {
		return givenPath
	}
	if !strings.HasPrefix(givenPath, "/") {
		path.WriteString("/")
	}
	path.WriteString(givenPath)
	return path.String()
}
