package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var bits = flag.Int("b", 256, "digest bits (256 or 512)")

func main() {
	flag.Parse()

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	digest, err := createHash(buf, *bits)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", digest)
}

func createHash(input []byte, bits int) ([]byte, error) {
	switch bits {
	case 256:
		digest := sha256.Sum256(input)
		return digest[:], nil
	case 512:
		digest := sha512.Sum512(input)
		return digest[:], nil
	default:
		return nil, fmt.Errorf("bits must be 256 or 512. bits=%d", bits)
	}
}
