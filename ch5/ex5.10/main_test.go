package main

import "testing"

func TestTopSort(t *testing.T) {
	order := topoSort(createPrereqMap())

	prereqMap := createPrereqMap()
	for i, item := range order {
		m := prereqMap[item]
		for j := 0; j < i; j++ {
			o := order[j]
			if _, ok := m[o]; ok {
				delete(m, o)
			}
		}
		if len(m) != 0 {
			t.Errorf("item %q is before items %v; want after those", item, m)
		}
	}
}
