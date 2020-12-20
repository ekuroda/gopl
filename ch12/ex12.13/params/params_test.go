package params

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"
)

type data struct {
	Foo []int  `http:"f" rule:"gt0"`
	Bar string `http:"bar" rule:"email"`
	Baz bool
}

var emailPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")

var emailValidator = func(v string) error {
	if emailPattern.MatchString(v) {
		return nil
	}
	return fmt.Errorf("invalid email: %v", v)
}

var gt0Validator = func(v string) error {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fmt.Errorf("not int: %v", i)
	}

	if i > 0 {
		return nil
	}
	return fmt.Errorf("not gt 0: %v", i)
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		url string
		ok  bool
	}{
		{"?f=1&bar=foo@example.com", true},
		{"?f=1&bar=hoge&f=2&baz=true", false},
		{"?f=a&bar=foo@example.com", false},
		{"?f=-2&bar=foo@example.com", false},
	}

	validators := make(map[string]Validator)
	validators["email"] = emailValidator
	validators["gt0"] = gt0Validator

	for _, test := range tests {
		req, _ := http.NewRequest("GET", test.url, &bytes.Reader{})
		var d data
		err := Unpack(req, &d, validators)
		if err == nil {
			if !test.ok {
				t.Errorf("Unpack() url: %v got nil error, want error", test.url)
			}
		} else {
			if test.ok {
				t.Errorf("Unpack() url: %v got %v, want nil error", test.url, err)
			}
		}
	}
}
