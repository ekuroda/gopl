package params

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

type data struct {
	Foo []int  `http:"f"`
	Bar string `http:"bar"`
	Baz bool
}

func TestUnpack(t *testing.T) {
	req, _ := http.NewRequest("GET", "?f=1&bar=hoge&f=2&baz=true", &bytes.Reader{})

	var d data
	if err := Unpack(req, &d); err != nil {
		t.Fatalf("failed to unpack: %v", err)
	}

	want := data{
		Foo: []int{1, 2},
		Bar: "hoge",
		Baz: true,
	}

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Unpack() got %v, want %v", d, want)
	}
	//fmt.Printf("%+v\n", d)
}

func TestPack(t *testing.T) {
	d := data{
		Foo: []int{1, 2},
		Bar: "hoge",
		Baz: true,
	}

	url, err := Pack(&d)
	if err != nil {
		t.Fatalf("failed to pack: %v", d)
	}

	req, _ := http.NewRequest("GET", url, &bytes.Reader{})

	var rd data
	if err := Unpack(req, &rd); err != nil {
		t.Fatalf("failed to unpack: %v", err)
	}

	if !reflect.DeepEqual(rd, d) {
		t.Errorf("Pack() got %v (%+v), want %+v", url, d, rd)
	}

	var x [2]int
	if _, err := Pack(&x); err == nil {
		t.Errorf("Pack(%v) got nil error, want valid error", x)

	}
}
