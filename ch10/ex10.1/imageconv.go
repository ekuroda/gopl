package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

var outputFormat = flag.String("o", "", "gif, png or jpeg")

func main() {
	flag.Parse()

	if *outputFormat == "" {
		fmt.Fprintln(os.Stderr, "usage: imageconv -o=gif|png|jpg")
		os.Exit(1)
	}

	if err := convert(*outputFormat, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func convert(format string, in io.Reader, out io.Writer) (err error) {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Input format =", kind)

	switch format {
	case "gif":
		err = gif.Encode(out, img, nil)
	case "png":
		err = png.Encode(out, img)
	case "jpeg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	default:
		err = fmt.Errorf("unknown output format %s", format)
	}

	return
}
