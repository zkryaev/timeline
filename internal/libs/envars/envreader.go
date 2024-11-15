package envars

import (
	"os"
	"strings"
)

// get path from home dir
// example: /home/$USER/Timeline/$GIVEN_PATH
func GetPath(envName string) string {
	pathToFile := os.Getenv(envName)
	rootProjectDir, _ := os.Getwd()
	path := strings.Builder{}
	path.Grow(len(rootProjectDir) + len(pathToFile))
	path.WriteString(rootProjectDir)
	if strings.Contains(pathToFile, path.String()) {
		return pathToFile
	}
	if !strings.HasPrefix(pathToFile, "/") {
		path.WriteString("/")
	}
	path.WriteString(pathToFile)
	return path.String()
}
