package word

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

func randomPalindrome(rng *rand.Rand, puncts, spaces []rune) string {
	n := rng.Intn(25)
	punctCount := rng.Intn(4)
	spaceCount := rng.Intn(4)
	exlen := n + punctCount + spaceCount
	indices := make([]int, exlen)

	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}

	for i := 0; i < exlen; i++ {
		indices[i] = i
	}
	shuffle(indices, rng)

	punctMap := make(map[int]rune)
	for i := 0; i < punctCount; i++ {
		punctMap[indices[i]] = puncts[rng.Intn(len(puncts))]
	}

	spaceIndices := indices[punctCount:]
	spaceMap := make(map[int]rune)
	for i := 0; i < spaceCount; i++ {
		spaceMap[spaceIndices[i]] = spaces[rng.Intn(len(spaces))]
	}

	expanded := make([]rune, exlen)
	runesIndex := 0
	for i := 0; i < exlen; i++ {
		if p, ok := punctMap[i]; ok {
			expanded[i] = p
		} else if s, ok := spaceMap[i]; ok {
			expanded[i] = s
		} else {
			expanded[i] = runes[runesIndex]
			runesIndex++
		}
	}

	//fmt.Printf("indices=%v, puncts=%v, spaces=%v, expanded=%q\n", indices, punctMap, spaceMap, expanded)
	return string(expanded)
}

func shuffle(nums []int, rng *rand.Rand) {
	n := len(nums)
	for i := n - 1; i >= 0; i-- {
		j := rng.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func TestRandomPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	puncts := make([]rune, 0)
	for i := 0x00; i <= 0x7f; i++ {
		r := rune(i)
		if unicode.IsPunct(r) {
			puncts = append(puncts, r)
		}
	}

	spaces := []rune{'\t', '\n', '\v', '\f', '\r', ' ', '\u0085', '\u00A0'}

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng, puncts, spaces)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}
