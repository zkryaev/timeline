package validation

import (
	"fmt"
	"net/http"
)

// Принимает список ожидаемых параметров. Проверяет их наличие в запросе
func IsQueryValid(r *http.Request, params map[string]bool) error {
	for param, required := range params {
		// Если параметр не обязательный то ошибки не будет
		if !r.URL.Query().Has(param) && required {
			return fmt.Errorf("%s required but didn't provide", param)
		}
	}
	return nil
}
