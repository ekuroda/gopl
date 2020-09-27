package counter

import (
	"bufio"
	"bytes"
)

// WordCounter ...
type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	var count int
	for {
		advance, _, err := bufio.ScanWords(p, true)
		if err != nil {
			return 0, err
		}
		if advance == 0 {
			break
		}
		p = p[advance:]
		count++
	}
	*c += WordCounter(count)
	return count, nil
}

// LineCounter ...
type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	var count int
	for {
		advance, _, err := bufio.ScanLines(p, true)
		if err != nil {
			return 0, err
		}
		if advance == 0 {
			break
		}
		p = p[advance:]
		count++
	}
	*c += LineCounter(count)
	return count, nil
}

// WordCounter2 ...
type WordCounter2 int

func (c *WordCounter2) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	var count int
	for scanner.Scan() {
		count++
	}
	*c += WordCounter2(count)
	return count, nil
}

// LineCounter2 ...
type LineCounter2 int

func (c *LineCounter2) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines)
	var count int
	for scanner.Scan() {
		count++
	}
	*c += LineCounter2(count)
	return count, nil
}
