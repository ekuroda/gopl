package word

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

func randomNonPalindrome(rng *rand.Rand) string {
	// 0 or 1文字は回文なので2文字以上の文字列生成
	n := 2 + rng.Intn(23)
	runes := make([]rune, 0)
	letters := make([]rune, 0)
	letterIndice := make([]int, 0)
	for i := 0; len(letters) < n; i++ {
		r := rune(rng.Intn(0x1000))
		runes = append(runes, r)
		if unicode.IsLetter(r) {
			letters = append(letters, r)
			letterIndice = append(letterIndice, i)
		}
	}

	index := rng.Intn(n / 2)
	for {
		// 元々回文でないかも知れないが、
		// 適当に選んだ先頭,最後からの文字が同じなら変えてしまえば回文でない
		if letters[index] == letters[n-1-index] {
			r := rune(rng.Intn(0x1000))
			letters[index] = r
			runes[letterIndice[index]] = r
		} else {
			break
		}
	}
	return string(runes)
}

func TestRandomPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}

func TestRandomNonPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomNonPalindrome(rng)
		if IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = true", p)
		}
	}
}
