package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"timeline/internal/entity/dto/authdto"

	"github.com/stretchr/testify/require"
)

func TestDecodeAndValidateSuccess(t *testing.T) {
	correctJSON := `{
		"email":"test@email.ru",
		"password":"nevergonnagiveyouup"
	}
	`
	dst := authdto.LoginReq{}
	r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(correctJSON))
	require.NoError(t, DecodeAndValidate(r, &dst))
}

func TestDecodeJsonError(t *testing.T) {
	invalidJSON := `{
		"email":"never"
		"password":gonna,
	}
	`
	r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(invalidJSON))
	require.Error(t, DecodeAndValidate(r, nil))
}

func TestValidateError(t *testing.T) {
	dst := authdto.LoginReq{
		Credentials: authdto.Credentials{Email: "testemail.ru", Password: "never"},
	}
	require.Error(t, Validate(&dst))
}
