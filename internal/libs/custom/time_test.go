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
			a:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 00, 00, 00, time.UTC),
			b:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 1, 00, 00, time.UTC),
			exp:  -1,
		},
		{
			name: "a later on minute than b",
			a:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 1, 0, 00, time.UTC),
			b:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 00, 00, 00, time.UTC),
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
			a:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 16, 00, 00, 00, time.UTC),
			b:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 00, 00, 00, time.UTC),
			exp:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, CompareTime(tc.a, tc.b))
		})
	}
}
