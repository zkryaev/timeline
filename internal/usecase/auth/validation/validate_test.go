package validation

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsCodeExpired(t *testing.T) {
	cases := []struct {
		name string
		in   time.Time
		exp  bool
	}{
		{
			name: "fresh",
			in:   time.Now().UTC(),
			exp:  false,
		},
		{
			name: "expired",
			in:   time.Time{},
			exp:  true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, IsCodeExpired(tc.in))
		})
	}
}

func TestIsAccountExpired(t *testing.T) {
	cases := []struct {
		name string
		in   time.Time
		exp  bool
	}{
		{
			name: "fresh",
			in:   time.Now(),
			exp:  false,
		},
		{
			name: "expired",
			in:   time.Time{},
			exp:  true,
		},
		{
			name: "almost expired",
			in:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Add(-23*time.Hour).Hour(), 59, 59, 59, time.UTC),
			exp:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, IsAccountExpired(tc.in))
		})
	}
}

func TestValidateTokenClaims(t *testing.T) {
	cases := []struct {
		name   string
		claims jwt.Claims
		exp    error
	}{
		{
			name: "correct claims",
			claims: jwt.MapClaims{
				"id":     1.0,
				"is_org": false,
				"type":   "access",
			},
			exp: nil,
		},
		{
			name: "id type is incorrect",
			claims: jwt.MapClaims{
				"id":     1,
				"is_org": false,
				"type":   "access",
			},
			exp: ErrWrongClaims,
		},
		{
			name: "is_org type is incorrect",
			claims: jwt.MapClaims{
				"id":     1.0,
				"is_org": "false",
				"type":   "access",
			},
			exp: ErrWrongClaims,
		},
		{
			name: "type field type is incorrect",
			claims: jwt.MapClaims{
				"id":     1.0,
				"is_org": false,
				"type":   666999,
			},
			exp: ErrWrongClaims,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, ValidateTokenClaims(tc.claims))
		})
	}
}
