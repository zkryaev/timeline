package validation

import (
	"time"

	"github.com/go-playground/validator"
)

func validTime(fl validator.FieldLevel) bool {
	const timeFormat = "15:04"
	inputTime := fl.Field().String()

	// Парсим строку в time.Time
	_, err := time.Parse(timeFormat, inputTime)
	// Если ошибка парсинга, то возвращаем false (значит, время некорректное)
	return err == nil
}

func NewCustomValidator() *validator.Validate {
	validate := validator.New()

	// Регистрация кастомных функций валидации полей
	validate.RegisterValidation("time", validTime)

	return validate
}
