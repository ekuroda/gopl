package sexpr

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestToken(t *testing.T) {
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
	tokens := make([]Token, 0)
	for {
		tok, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				t.Fatalf("faied to get token %q: %v", data, err)
			}
			break
		}
		tokens = append(tokens, tok)
	}

	wants := []Token{
		StartList{},
		StartList{},
		Symbol("Title"), String("Dr. Strangelove"),
		EndList{},
		StartList{},
		Symbol("Subtitle"), String("How I Learned to Stop Worrying and Love the Bomb"),
		EndList{},
		StartList{},
		Symbol("Year"), Int(1964),
		EndList{},
	}

	for i, want := range wants {
		if !reflect.DeepEqual(tokens[i], want) {
			t.Errorf("Token(); got %T %v, want %T %v", tokens[i], tokens[i], want, want)
		}
	}
}
