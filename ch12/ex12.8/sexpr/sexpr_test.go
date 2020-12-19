package sexpr

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	data, err := Marshal(strangelove)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	//t.Logf("Marshal() = %s\n", data)

	reader := bytes.NewReader(data)
	decoder := NewDecoder(reader)
	var m Movie
	if err := decoder.Decode(&m); err != nil {
		t.Fatalf("failed to unmarshal %q: %v", data, err)
	}
	t.Logf("Decode() = %v\n", m)

	if !reflect.DeepEqual(m, strangelove) {
		t.Errorf("Decode %v; got %v, want %v", string(data), m, strangelove)
	}
}