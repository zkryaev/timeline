package envars

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathByEnv(t *testing.T) {
	currDir, _ := os.Getwd()
	projectRootDir := strings.SplitAfter(currDir, "timeline/")[0]

	cases := map[string]struct {
		env, val, exp string
	}{
		"existed_env": {
			env: "TEST",
			val: "testing",
			exp: projectRootDir + "/testing",
		},
		"unexist_env": {
			env: "NONAME",
			val: "",
			exp: "",
		},
	}

	require.NoError(t, os.Setenv(cases["existed_env"].env, cases["existed_env"].val))

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, c.exp, GetPathByEnv(c.env))
		})
	}
}
