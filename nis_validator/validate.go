package nis_validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

const tagType = "json"

var ErrValidation = errors.New("input validation error")

// ValidateStructReturnFirstError Validates and returns first error
func ValidateStructReturnFirstError(structToValidate any) (string, error) {
	errMsg := ""
	v := GetNisValidator()
	err := v.Struct(structToValidate)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, validationError := range validationErrors {
			jsonFieldName := ""
			if field, ok := reflect.TypeOf(structToValidate).FieldByName(validationError.Field()); ok {
				jsonFieldName = field.Tag.Get(tagType)
			}

			if jsonFieldName == "" {
				jsonFieldName = validationError.Field()
			}

			switch validationError.ActualTag() {
			case "lte":
				errMsg = fmt.Sprintf("Should be less than or equal to %v", validationError.Param())
				//errMsg = fmt.Sprintf("exceeded value")
			case "lt":
				errMsg = fmt.Sprintf("Should be less than %v", validationError.Param())
			case "gte":
				errMsg = fmt.Sprintf("Should be greater than or equal to %v", validationError.Param())
			case "gt":
				errMsg = fmt.Sprintf("Should be greater than %v", validationError.Param())
			default:
				errMsg = validationError.ActualTag()
			}

			if errMsg != "" {
				return fmt.Sprintf("%s: %s", jsonFieldName, errMsg), ErrValidation
			}
		}
	}

	return "", nil
}

func GetNisValidator() *validator.Validate {
	validate := validator.New()
	// Slice validation
	validate.RegisterValidation("uniqueField", validateUniqueStructField)
	validate.RegisterValidation("atLeastOneField", validateAtLeastOne)
	validate.RegisterValidation("isLocaleField", validateIsLocaleField)
	//custom validation
	validate.RegisterValidation("isLocale", validateIsLocale)

	return validate
}

func validateUniqueStructField(fl validator.FieldLevel) bool {
	/*
		Eg usage
		type Data struct {
			Items []Item `validate:"uniqueField=Locale"`
		}
	*/
	field := fl.Field()
	if field.Kind() != reflect.Slice {
		panic("uniqueField validation can only be applied to a slice")
	}

	argument := fl.Param()

	values := make(map[interface{}]struct{})
	length := field.Len()

	for i := 0; i < length; i++ {
		itemValue := field.Index(i)
		argumentValue := itemValue.FieldByName(argument).Interface()

		if _, ok := values[argumentValue]; ok {
			return false
		}

		values[argumentValue] = struct{}{}
	}

	return true
}

func validateAtLeastOne(fl validator.FieldLevel) bool {
	/*	Eg usage
		type Data struct {
			Items []Item `validate:"atLeastOneField=Locale:ru"`
		}
	*/
	field := fl.Field()
	if field.Kind() != reflect.Slice {
		panic("atLeastOneField validation can only be applied to a slice")
	}

	argument := fl.Param()
	parts := strings.Split(argument, ":")

	if len(parts) != 2 {
		panic("invalid argument format for atLeastOneField validation")
	}

	fieldName := parts[0]
	requiredValue := parts[1]

	length := field.Len()

	for i := 0; i < length; i++ {
		itemValue := field.Index(i)
		fieldValue := itemValue.FieldByName(fieldName).Interface()

		if fmt.Sprintf("%v", fieldValue) == requiredValue {
			return true
		}
	}

	return false
}

func validateIsLocaleField(fl validator.FieldLevel) bool {
	/*
		Eg usage
		type Data struct {
			Items []Item `validate:"isLocaleField=Locale"`
		}
	*/
	field := fl.Field()
	if field.Kind() != reflect.Slice {
		panic("uniqueField validation can only be applied to a slice")
	}

	argument := fl.Param()

	length := field.Len()

	for i := 0; i < length; i++ {
		itemValue := field.Index(i)
		argumentValue := itemValue.FieldByName(argument).Interface()

		if langCode, ok := argumentValue.(string); !ok || !langCodes[langCode] {
			return false
		}
	}

	return true
}

func validateIsLocale(fl validator.FieldLevel) bool {
	/*
		Eg usage
		type Data struct {
			Items []Item `validate:"isLocaleField=Locale"`
		}
	*/
	field := fl.Field()
	langCode, ok := field.Interface().(string)

	if !ok {
		return false
	}

	if !langCodes[langCode] {
		return false
	}

	return true
}

var langCodes = map[string]bool{
	"aa": true,
	"ab": true,
	"ae": true,
	"af": true,
	"ak": true,
	"am": true,
	"an": true,
	"ar": true,
	"as": true,
	"av": true,
	"ay": true,
	"az": true,
	"ba": true,
	"be": true,
	"bg": true,
	"bi": true,
	"bm": true,
	"bn": true,
	"bo": true,
	"br": true,
	"bs": true,
	"ca": true,
	"ce": true,
	"ch": true,
	"co": true,
	"cr": true,
	"cs": true,
	"cv": true,
	"cy": true,
	"da": true,
	"de": true,
	"dv": true,
	"dz": true,
	"ee": true,
	"el": true,
	"en": true,
	"eo": true,
	"es": true,
	"et": true,
	"eu": true,
	"fa": true,
	"ff": true,
	"fi": true,
	"fj": true,
	"fo": true,
	"fr": true,
	"fy": true,
	"ga": true,
	"gd": true,
	"gl": true,
	"gn": true,
	"gu": true,
	"gv": true,
	"ha": true,
	"he": true,
	"hi": true,
	"ho": true,
	"hr": true,
	"ht": true,
	"hu": true,
	"hy": true,
	"hz": true,
	"ia": true,
	"id": true,
	"ie": true,
	"ig": true,
	"ii": true,
	"ik": true,
	"io": true,
	"is": true,
	"it": true,
	"iu": true,
	"ja": true,
	"jv": true,
	"ka": true,
	"kg": true,
	"ki": true,
	"kj": true,
	"kk": true,
	"kl": true,
	"km": true,
	"kn": true,
	"ko": true,
	"kr": true,
	"ks": true,
	"ku": true,
	"kv": true,
	"kw": true,
	"ky": true,
	"la": true,
	"lb": true,
	"lg": true,
	"li": true,
	"ln": true,
	"lo": true,
	"lt": true,
	"lu": true,
	"lv": true,
	"mg": true,
	"mh": true,
	"mi": true,
	"mk": true,
	"ml": true,
	"mn": true,
	"mr": true,
	"ms": true,
	"mt": true,
	"my": true,
	"na": true,
	"nb": true,
	"nd": true,
	"ne": true,
	"ng": true,
	"nl": true,
	"nn": true,
	"no": true,
	"nr": true,
	"nv": true,
	"ny": true,
	"oc": true,
	"oj": true,
	"om": true,
	"or": true,
	"os": true,
	"pa": true,
	"pi": true,
	"pl": true,
	"ps": true,
	"pt": true,
	"qu": true,
	"rm": true,
	"rn": true,
	"ro": true,
	"ru": true,
	"rw": true,
	"sa": true,
	"sc": true,
	"sd": true,
	"se": true,
	"sg": true,
	"si": true,
	"sk": true,
	"sl": true,
	"sm": true,
	"sn": true,
	"so": true,
	"sq": true,
	"sr": true,
	"ss": true,
	"st": true,
	"su": true,
	"sv": true,
	"sw": true,
	"ta": true,
	"te": true,
	"tg": true,
	"th": true,
	"ti": true,
	"tk": true,
	"tl": true,
	"tn": true,
	"to": true,
	"tr": true,
	"ts": true,
	"tt": true,
	"tw": true,
	"ty": true,
	"ug": true,
	"uk": true,
	"ur": true,
	"uz": true,
	"ve": true,
	"vi": true,
	"vo": true,
	"wa": true,
	"wo": true,
	"xh": true,
	"yi": true,
	"yo": true,
	"za": true,
	"zh": true,
	"zu": true,
}
