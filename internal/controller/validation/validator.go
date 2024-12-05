package validation

import (
	"time"

	"github.com/go-playground/validator"
)

const (
	dateFormat = "2001-01-01"
	timeFormat = "15:04"
)

func validTime(fl validator.FieldLevel) bool {
	inputTime := fl.Field().String()

	// Парсим строку в time.Time
	_, err := time.Parse(timeFormat, inputTime)
	// Если ошибка парсинга, то возвращаем false (значит, время некорректное)
	return err == nil
}

func validDate(fl validator.FieldLevel) bool {
	inputDate := fl.Field().String()

	_, err := time.Parse(dateFormat, inputDate)

	return err == nil
}

func NewCustomValidator() *validator.Validate {
	validate := validator.New()

	// Регистрация кастомных функций валидации полей
	validate.RegisterValidation("time", validTime)
	validate.RegisterValidation("date", validDate)

	return validate
}
