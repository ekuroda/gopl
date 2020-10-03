package sorting

import (
	"sort"
	"testing"
)

func TestMultiColumnSort(t *testing.T) {
	ts := make([]*Track, len(tracks))
	copy(ts, tracks)

	s := NewMultiColumnSort(ts)
	s.SelectTitle()
	s.SelectYear()
	sort.Sort(s)

	printTracks(ts)

	if ts[0].Title != "Go" || ts[1].Title != "Go" || ts[0].Year != 1992 || ts[1].Year != 2012 {
		t.Errorf("sort by title and year: got first, second = (%q, %d), (%q, %d), want (%q, %d), (%q, %d)",
			ts[0].Title, ts[0].Year, ts[1].Title, ts[1].Year, "Go", 1992, "Go", 2012)
	}

	copy(ts, tracks)
	sort.Stable(ByYear(ts))
	sort.Stable(ByTitle(ts))
	printTracks(ts)
}
