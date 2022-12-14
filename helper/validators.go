package helper

import (
	"errors"
	"log"
	"reflect"
	"strings"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateJSON(data interface{}) error {
	translator := english.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	validate := validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

	_ = validate.RegisterTranslation(
		"required",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)

	_ = validate.RegisterTranslation(
		"email",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		},
	)

	_ = validate.RegisterTranslation(
		"password",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("password", "{0} is not strong enough", true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		},
	)

	_ = validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {
			return len(fl.Field().String()) > 6
		},
	)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	validationError := validate.Struct(data)
	if validationError != nil {
		var outcomeString string
		for _, e := range validationError.(validator.ValidationErrors) {
			currentString := e.Translate(trans) + "\n"

			outcomeString = outcomeString + currentString
		}

		return errors.New(strings.TrimSpace(outcomeString))
	}

	return nil
}
