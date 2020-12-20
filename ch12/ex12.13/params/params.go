package params

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Validator ...
type Validator func(v string) error

// Unpack ...
func Unpack(req *http.Request, ptr interface{}, validators map[string]Validator) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	fields := make(map[string]reflect.Value)
	formValidatorMap := make(map[string]Validator)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)

		ruleName := tag.Get("rule")
		if ruleName != "" {
			if validator, ok := validators[ruleName]; ok {
				formValidatorMap[name] = validator
			}
		}
	}

	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}
		for _, value := range values {
			if validator, ok := formValidatorMap[name]; ok {
				if err := validator(value); err != nil {
					return err
				}
			}
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}

		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported king %s", v.Type())
	}
	return nil
}
