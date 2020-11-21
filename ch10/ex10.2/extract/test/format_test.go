package test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"gopl/ch10/ex10.2/extract"
	_ "gopl/ch10/ex10.2/extract/tar"
	_ "gopl/ch10/ex10.2/extract/zip"
)

func TestZip(t *testing.T) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	archive, kind, err := extract.Decode(buf)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if kind != "zip" {
		t.Errorf("kind = %s, want %s", kind, "zip")
	}

	i := -1
	for {
		err := archive.Next()
		if err == io.EOF {
			break
		}

		i++
		if err != nil {
			t.Error(err)
			continue
		}
		buf, err := ioutil.ReadAll(archive)
		if err != nil {
			t.Errorf("failed to read %s: %v", archive.Name(), err)
			continue
		}

		s := string(buf)
		fmt.Printf("Contents of %s:\n%s\n", archive.Name(), s)
		if s != files[i].Body {
			t.Errorf("s = %s, want %s", s, files[i].Body)
		}
	}

	if i != len(files)-1 {
		t.Errorf("i = %d, want %d", i, len(files)-1)
	}

	// b := buf.Bytes()
	// readerAt := bytes.NewReader(b)
	// r, err := zip.NewReader(readerAt, int64(len(b)))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// i := -1
	// for _, f := range r.File {
	// 	i++
	// 	rc, err := f.Open()
	// 	if err != nil {
	// 		t.Fatalf("failed to open %s: %v\n", f.Name, err)
	// 		continue
	// 	}

	// 	buf, err := ioutil.ReadAll(rc)
	// 	if err != nil {
	// 		t.Errorf("failed to read %s: %v\n", f.Name, err)
	// 		rc.Close()
	// 		continue
	// 	}

	// 	s := string(buf)
	// 	fmt.Printf("Contents of %s:\n%s\n", f.Name, s)
	// 	if s != files[i].Body {
	// 		t.Errorf("s = %s, want %s", s, files[i].Body)
	// 	}
	// 	rc.Close()
	// }

	// if i != len(files)-1 {
	// 	t.Errorf("i = %d, want %d", i, len(files)-1)
	// }
}

func TestTar(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling license."},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	archive, kind, err := extract.Decode(&buf)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if kind != "tar" {
		t.Errorf("kind = %s, want %s", kind, "tar")
	}

	//tr := tar.NewReader(&buf)
	i := -1
	for {
		err := archive.Next()
		if err == io.EOF {
			break
		}

		i++
		if err != nil {
			t.Error(err)
			continue
		}
		buf, err := ioutil.ReadAll(archive)
		if err != nil {
			t.Errorf("failed to read %s: %v", archive.Name(), err)
			continue
		}

		s := string(buf)
		fmt.Printf("Contents of %s:\n%s\n", archive.Name(), s)
		if s != files[i].Body {
			t.Errorf("s = %s, want %s", s, files[i].Body)
		}
	}

	if i != len(files)-1 {
		t.Errorf("i = %d, want %d", i, len(files)-1)
	}
}
