package sorting

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

// Track ...
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

type less func(a, b *Track) bool

var byTitle less = func(a, b *Track) bool { return a.Title < b.Title }
var byArtist less = func(a, b *Track) bool { return a.Artist < b.Artist }
var byAlbum less = func(a, b *Track) bool { return a.Album < b.Album }
var byYear less = func(a, b *Track) bool { return a.Year < b.Year }
var byLength less = func(a, b *Track) bool { return a.Length < b.Length }

// MultiColumnSort ...
type MultiColumnSort struct {
	t      []*Track
	lesses []less
}

// NewMultiColumnSort ...
func NewMultiColumnSort(t []*Track) *MultiColumnSort {
	return &MultiColumnSort{t, make([]less, 0)}
}

// SelectTitle ...
func (s *MultiColumnSort) SelectTitle() {
	s.lesses = append(s.lesses, byTitle)
}

// SelectArtist ...
func (s *MultiColumnSort) SelectArtist() {
	s.lesses = append(s.lesses, byArtist)
}

// SelectAlbum ...
func (s *MultiColumnSort) SelectAlbum() {
	s.lesses = append(s.lesses, byAlbum)
}

// SelectYear ...
func (s *MultiColumnSort) SelectYear() {
	s.lesses = append(s.lesses, byYear)
}

// SelectLength ...
func (s *MultiColumnSort) SelectLength() {
	s.lesses = append(s.lesses, byLength)
}

// Len ...
func (s *MultiColumnSort) Len() int { return len(s.t) }

// Less ...
func (s *MultiColumnSort) Less(i, j int) bool {
	for i := len(s.lesses) - 1; i >= 0; i-- {
		less := s.lesses[i]
		if r := less(s.t[i], s.t[j]); r {
			return r
		}
		if r := less(s.t[j], s.t[i]); r {
			return !r
		}
	}

	return false
}

// Swap ...
func (s *MultiColumnSort) Swap(i, j int) {
	s.t[i], s.t[j] = s.t[j], s.t[i]
}

// ByTitle ...
type ByTitle []*Track

func (x ByTitle) Len() int           { return len(x) }
func (x ByTitle) Less(i, j int) bool { return x[i].Title < x[j].Title }
func (x ByTitle) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// ByArtist ...
type ByArtist []*Track

func (x ByArtist) Len() int           { return len(x) }
func (x ByArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x ByArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// ByAlbum ...
type ByAlbum []*Track

func (x ByAlbum) Len() int           { return len(x) }
func (x ByAlbum) Less(i, j int) bool { return x[i].Album < x[j].Album }
func (x ByAlbum) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// ByYear ...
type ByYear []*Track

func (x ByYear) Len() int           { return len(x) }
func (x ByYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x ByYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// ByLength ...
type ByLength []*Track

func (x ByLength) Len() int           { return len(x) }
func (x ByLength) Less(i, j int) bool { return x[i].Length < x[j].Length }
func (x ByLength) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
