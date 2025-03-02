package validation

import (
	"time"

	"github.com/go-playground/validator"
)

const (
	TimeOnlyHM = "15:04"
)

func validTime(fl validator.FieldLevel) bool {
	inputTime := fl.Field().String()

	// Парсим строку в time.Time
	_, err := time.Parse(TimeOnlyHM, inputTime)
	// Если ошибка парсинга, то возвращаем false (значит, время некорректное)
	return err == nil
}

func validDate(fl validator.FieldLevel) bool {
	inputDate := fl.Field().String()

	_, err := time.Parse(time.DateOnly, inputDate)

	return err == nil
}

func NewCustomValidator() (*validator.Validate, error) {
	validate := validator.New()

	validationList := []struct {
		field string
		check func(fl validator.FieldLevel) bool
	}{
		{
			field: "time",
			check: validTime,
		},
		{
			field: "date",
			check: validDate,
		},
	}

	for _, v := range validationList {
		if err := validate.RegisterValidation(v.field, v.check); err != nil {
			return nil, err
		}
	}

	return validate, nil
}
