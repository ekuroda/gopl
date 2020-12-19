package params

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Unpack ...
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}
		for _, value := range values {
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

// Pack ...
func Pack(data interface{}) (string, error) {
	v := reflect.ValueOf(data).Elem()
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("data %v must be a struct", data)
	}

	urlVals := &url.Values{}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fi := v.Type().Field(i)
		tag := fi.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fi.Name)
		}

		switch f.Kind() {
		case reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				switch f.Index(j).Kind() {
				case reflect.Bool, reflect.Int, reflect.String:
					urlVals.Add(name, fmt.Sprintf("%v", f.Index(j)))
				}
			}
		case reflect.Bool, reflect.Int, reflect.String:
			urlVals.Add(name, fmt.Sprintf("%v", v.Field(i)))
		}
	}

	u := url.URL{RawQuery: urlVals.Encode()}
	return u.String(), nil
}
