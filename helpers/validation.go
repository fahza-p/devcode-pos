package helpers

import (
	"devcode-pos/models"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	// Init Translator
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")

	// Init Validate
	validate = validator.New()

	// Register Translator
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Register Custom Rule
	validate.RegisterValidation("exist", exist)
	validate.RegisterValidation("enum", enum)

	// Register Custom Rule Message
	addTranslation("exist", "{0} with value {1} is not exist")
	addTranslation("enum", "{0} with value {1} is not exist")
}

// Run Valiedation for Custom Error Message
func RunValidate(input interface{}) (map[string]interface{}, error) {
	err := validate.Struct(input)
	var msg []interface{}
	var text []string

	if err != nil {
		errs := err.(validator.ValidationErrors)
		for key, val := range errs.Translate(trans) {
			structField, _ := reflect.TypeOf(input).Elem().FieldByName(key)
			jsonTag := structField.Tag.Get("json")
			text = append(text, val)
			msg = append(msg, map[string]interface{}{
				"message": val,
				"path":    []string{jsonTag},
				"type":    "any.required",
				"context": map[string]string{
					"label": jsonTag,
					"key":   jsonTag,
				},
			})
		}
	}

	jsonEncode, e := json.Marshal(text)

	if e != nil {
		err = e
	}

	result := map[string]interface{}{
		"text": "body ValidationError: " + string(jsonEncode),
		"err":  msg,
	}

	return result, err
}

// Add Message Translation
func addTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		tag := fe.Tag()
		t, err := ut.T(tag, fe.Field(), fmt.Sprintf("%v", fe.Value()))

		if err != nil {
			return fe.(error).Error()
		}

		return t
	}

	validate.RegisterTranslation(tag, trans, registerFn, transFn)
}

// Custome Rule
func exist(fl validator.FieldLevel) bool {
	var res int64
	splitParam := strings.Split(fl.Param(), ".")
	value := fl.Field().Interface()

	switch splitParam[2] {
	case "int":
		value = value.(uint32)
	case "string":
		value = value.(string)
	}

	models.DB.Table(splitParam[0]).
		Select(fmt.Sprintf("count(distinct(%s))", splitParam[1])).
		Where(splitParam[1]+" = ?", value).
		Count(&res)

	return res > 0
}

func enum(fl validator.FieldLevel) bool {
	splitParam := strings.Split(fl.Param(), ";")
	value := fl.Field().Interface()
	res := false
	for _, el := range splitParam {
		if value == el {
			res = true
		}
	}
	return res
}
