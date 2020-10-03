package sorting

import (
	"html/template"
	"io"
	"log"
	"time"
)

const tracksTemplateText = `
<h2>Tracks</h2>
<table>
  <tr style='text-align:left'>
    <th><a href='/?s=Title'>Title</a></th>
    <th><a href='/?s=Artist'>Artist</a></th>
    <th><a href='/?s=Album'>Album</a></th>
    <th><a href='/?s=Year'>Year</a></th>
	<th><a href='/?s=Length'>Length</a></th>
  </tr>
  {{range .}}
  <tr>
    <td>{{.Title}}</td>
    <td>{{.Artist}}</td>
    <td>{{.Album}}</td>
    <td>{{.Year}}</td>
    <td>{{.Length}}</td>
  </tr>
  {{end}}
</table>
`

var tracksTemplate *template.Template = template.Must(template.New("tracks").Parse(tracksTemplateText))

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

// PrintTracks ...
func PrintTracks(w io.Writer) error {
	if err := tracksTemplate.Execute(w, tracks); err != nil {
		log.Printf("failed to write tracks template: %s", err)
		return err
	}
	return nil
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
func NewMultiColumnSort() *MultiColumnSort {
	return &MultiColumnSort{tracks, make([]less, 0)}
}

// SelectTitle ...
func (s *MultiColumnSort) SelectTitle() {
	s.lesses = append(s.lesses, byTitle)
	if len(s.lesses) > 5 {
		s.lesses = s.lesses[1:]
	}
}

// SelectArtist ...
func (s *MultiColumnSort) SelectArtist() {
	s.lesses = append(s.lesses, byArtist)
	if len(s.lesses) > 5 {
		s.lesses = s.lesses[1:]
	}
}

// SelectAlbum ...
func (s *MultiColumnSort) SelectAlbum() {
	s.lesses = append(s.lesses, byAlbum)
	if len(s.lesses) > 5 {
		s.lesses = s.lesses[1:]
	}
}

// SelectYear ..
func (s *MultiColumnSort) SelectYear() {
	s.lesses = append(s.lesses, byYear)
	if len(s.lesses) > 5 {
		s.lesses = s.lesses[1:]
	}
}

// SelectLength ...
func (s *MultiColumnSort) SelectLength() {
	s.lesses = append(s.lesses, byLength)
	if len(s.lesses) > 5 {
		s.lesses = s.lesses[1:]
	}
}

// Len ...
func (s *MultiColumnSort) Len() int { return len(s.t) }

// Less ...
func (s *MultiColumnSort) Less(i, j int) bool {
	for li := len(s.lesses) - 1; li >= 0; li-- {
		less := s.lesses[li]
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
