package envars

import (
	"fmt"
	"os"
	"strings"
)

// get path from home dir
// example: /home/$USER/Timeline/$GIVEN_PATH
func GetPath(envName string) string {
	user := os.Getenv("USER")
	givenPath := os.Getenv(envName)
	base := strings.Builder{}
	base.Grow(len("/home") + len(user) + len("/Timeline"))
	base.WriteString("/home/")
	base.WriteString(user)
	base.WriteString("/Timeline")
	if strings.Contains(givenPath, base.String()) {
		return givenPath
	}
	return fmt.Sprintf("%s/%s", base.String(), givenPath)
}
