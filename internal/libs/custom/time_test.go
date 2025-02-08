package custom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompareTime(t *testing.T) {
	cases := []struct {
		name string
		a    time.Time
		b    time.Time
		exp  int
	}{
		{
			name: "a earlier on minute than b",
			a:    time.Now(),
			b:    time.Now().Add(1 * time.Minute),
			exp:  -1,
		},
		{
			name: "a later on minute than b",
			a:    time.Now().Add(1 * time.Minute),
			b:    time.Now(),
			exp:  1,
		},
		{
			name: "a equal b",
			a:    time.Now(),
			b:    time.Now(),
			exp:  0,
		},
		{
			name: "a later on hour than b",
			a:    time.Now().Add(1 * time.Hour),
			b:    time.Now(),
			exp:  1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, CompareTime(tc.a, tc.b))
		})
	}
}
