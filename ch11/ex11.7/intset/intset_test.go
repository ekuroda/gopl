package intset

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkIntSetHas(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	setnum := 100
	addnum := 1000
	intSets := make([]*IntSet, setnum, setnum)
	for i := range intSets {
		is := new(IntSet)
		intSets[i] = is
		for j := 0; j < addnum; j++ {
			is.Add(rng.Intn(max))
		}
	}

	benchIntSetHas(b, rng, max, intSets)
}

func benchIntSetHas(b *testing.B, rng *rand.Rand, rndmax int, intSets []*IntSet) {
	l := len(intSets)
	for i := 0; i < b.N; i++ {
		intSets[i%l].Has(rng.Intn(rndmax))
	}
}

func BenchmarkMapIntSetHas(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	setnum := 100
	addnum := 1000
	intSets := make([]*MapIntSet, setnum, setnum)
	for i := range intSets {
		is := NewMapIntSet()
		intSets[i] = is
		for j := 0; j < addnum; j++ {
			is.Add(rng.Intn(max))
		}
	}

	benchMapIntSetHas(b, rng, max, intSets)
}

func benchMapIntSetHas(b *testing.B, rng *rand.Rand, rndmax int, intSets []*MapIntSet) {
	l := len(intSets)
	for i := 0; i < b.N; i++ {
		intSets[i%l].Has(rng.Intn(rndmax))
	}
}

func BenchmarkIntSetAdd(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	benchIntSetAdd(b, rng, max, new(IntSet))
}

func benchIntSetAdd(b *testing.B, rng *rand.Rand, rndmax int, intSet *IntSet) {
	for i := 0; i < b.N; i++ {
		intSet.Add(rng.Intn(rndmax))
	}
}

func BenchmarkMapIntSetAdd(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	benchMapIntSetAdd(b, rng, max, NewMapIntSet())
}

func benchMapIntSetAdd(b *testing.B, rng *rand.Rand, rndmax int, intSet *MapIntSet) {
	for i := 0; i < b.N; i++ {
		intSet.Add(rng.Intn(rndmax))
	}
}

func BenchmarkIntSetUnionWith(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	setnum := 100
	addnum := 1000
	intSets := make([]*IntSet, setnum, setnum)
	for i := range intSets {
		is := new(IntSet)
		intSets[i] = is
		for j := 0; j < addnum; j++ {
			is.Add(rng.Intn(max))
		}
	}

	benchIntSetUnionWith(b, rng, max, new(IntSet), intSets)
}

func benchIntSetUnionWith(b *testing.B, rng *rand.Rand, rndmax int, intSet *IntSet, intSets []*IntSet) {
	l := len(intSets)
	for i := 0; i < b.N; i++ {
		intSet.UnionWith(intSets[rng.Intn(l)])
	}
}

func BenchmarkMapIntSetUnionWith(b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	max := 10000
	setnum := 100
	addnum := 1000
	intSets := make([]*MapIntSet, setnum, setnum)
	for i := range intSets {
		is := NewMapIntSet()
		intSets[i] = is
		for j := 0; j < addnum; j++ {
			is.Add(rng.Intn(max))
		}
	}

	benchMapIntSetUnionWith(b, rng, max, NewMapIntSet(), intSets)
}

func benchMapIntSetUnionWith(b *testing.B, rng *rand.Rand, rndmax int, intSet *MapIntSet, intSets []*MapIntSet) {
	l := len(intSets)
	for i := 0; i < b.N; i++ {
		intSet.UnionWith(intSets[rng.Intn(l)])
	}
}
