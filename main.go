package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"main/nis_validator"
	"reflect"
	"strings"
)

type Item struct {
	ID     int
	Locale string
}

type Data struct {
	ID   int  `validate:"gt=0"`
	Item Item `innerValidation:"ID>>gt=0,
								Locale>>required,eq=2" json:"item"`
}

func main() {
	v := nis_validator.GetNisValidator()
	//var data Data = Data{Item: Item{ID: 1}}
	items := []Item{{1, "en"}, {Locale: "ru"}, {3, "en"}}

	fmt.Println(v.Var(items, "isLocaleField=Locale"))
	//fmt.Println(validateEmbeddedStructFields(data))
}

// TODO trim all values for whitespaces => maybe you can use multiline validation
func validateEmbeddedStructFields(data interface{}) error {
	nagibator := validator.New()
	val := reflect.ValueOf(data)
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := val.Field(i)

		if fieldValue.Kind() != reflect.Struct {
			continue
		}

		//ID>>gt=0;Name>>lt=10
		rawRules := field.Tag.Get("innerValidation")

		if rawRules == "" {
			continue
		}

		rules := make(map[string]string)

		for _, el := range strings.Split(rawRules, ";") {
			//ID>>gt=0
			if el == "" {
				continue
			}

			keyVal := strings.Split(el, ">>")

			if len(keyVal) != 2 {
				panic(fmt.Sprintf("validation use: Item `%s`", `innerValidation:"ID>>gt=0" json:"item"`))
			}

			rules[keyVal[0]] = keyVal[1]
		}

		for f, r := range rules {
			fieldVal := fieldValue.FieldByName(f).Interface()
			err := nagibator.Var(fieldVal, r)

			if err != nil {
				vd := err.(validator.ValidationErrors)

				for _, fe := range vd {
					return fmt.Errorf("field: %s failed validation for '%s' tag", f, fe.Tag())
				}
			}
		}
	}

	return nil
}

//func main() {
//	var data Item
//
//	fmt.Println(getInnerValidation(&data))
//	PrintStructTags(data)
//
//	//fmt.Println("----------------")
//	//
//	//PrintStructTags(data)
//	//
//	//fmt.Println("++++++++++++++++")
//	//
//	//PrintStructTags(data.Item)
//}
//
//func getInnerValidation(data any) error {
//	t := reflect.TypeOf(data).Elem()
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		//tag := field.Tag.Get("validate")
//		newTag := "required" //fmt.Sprintf("%s,%s", tag, "required")
//		field.Tag = reflect.StructTag(fmt.Sprintf(`%s:"%s"`, "validate", newTag))
//	}
//
//	return nil
//}

//	func getInnerValidation(data any) error {
//		fmt.Println("inner starting")
//		t := reflect.TypeOf(data).Elem()
//		val := reflect.ValueOf(data).Elem()
//
//		for i := 0; i < t.NumField(); i++ {
//			field := t.Field(i)
//
//			if field.Type.Kind() != reflect.Struct {
//				continue
//			}
//
//			//ID>>gt=0;ID>>lt=10
//			allRules := field.Tag.Get("innerValidation")
//			fmt.Println(allRules)
//
//			//ID>>gt=0
//			for _, rule := range strings.Split(allRules, ";") {
//				fieldAndRule := strings.Split(rule, ">>")
//
//				if len(fieldAndRule) != 2 {
//					panic(`Add text in following format - 'innerValidation:"ID>>gt=0;ID>>lt=10"'`)
//				}
//
//				fieldVal := val.FieldByName(field.Name).Interface()
//				embeddedT := reflect.TypeOf(fieldVal)
//
//				for i := 0; i < embeddedT.NumField(); i++ {
//					innerField := embeddedT.Field(i)
//
//					if innerField.Name != fieldAndRule[0] {
//						continue
//					}
//
//					tag := innerField.Tag.Get("validate")
//					newTag := fmt.Sprintf("%s,%s", tag, fieldAndRule[1])
//					innerField.Tag = reflect.StructTag(fmt.Sprintf(`%s:"%s"`, "validate", newTag))
//				}
//			}
//		}
//		fmt.Println("inner ending")
//		return validator.New().Struct(data)
//	}
//func PrintStructTags(data interface{}) {
//	t := reflect.TypeOf(data)
//
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		fmt.Printf("%s => %s\n", field.Name, field.Tag)
//	}
//}
