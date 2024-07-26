package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	initial = 0x0b704ce
	polynom = 0x1864cfb
)

type digest struct {
	sum uint32
}

func (d *digest) Write(p []byte) (n int, err error) {
	for _, v := range p {
		d.sum ^= uint32(v) << 16
		for i := 0; i < 8; i++ {
			d.sum <<= 1
			if d.sum&0x1000000 != 0 {
				d.sum ^= polynom
			}
		}
	}
	return len(p), nil
}

func (d *digest) Sum32() uint32 {
	return d.sum & 0xffffff
}

func (d *digest) Reset() {
	d.sum = initial
}

func New() *digest {
	d := &digest{}
	d.Reset()
	return d
}

func main() {
	base64Flag := flag.Bool("b", false, "Output in base64 OpenPGP format")
	flag.BoolVar(base64Flag, "base64", false, "Output in base64 OpenPGP format")
	helpFlag := flag.Bool("h", false, "Show help")
	flag.BoolVar(helpFlag, "help", false, "Show help")

	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: crc24 [-b|--base64] < file")
		fmt.Println("Ctrl+D (Unix) or Ctrl+Z (Windows) to finish input from stdin.")
		fmt.Println("  -b, --base64   Output the checksum in base64 OpenPGP format")
		fmt.Println("  -h, --help     Show this help message")
		os.Exit(0)
	}

	d := New()
	reader := bufio.NewReader(os.Stdin)

	_, err := io.Copy(d, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	checksum := d.Sum32()

	if *base64Flag {
		checksumBytes := []byte{byte(checksum >> 16), byte(checksum >> 8), byte(checksum)}
		encoded := base64.StdEncoding.EncodeToString(checksumBytes)
		fmt.Println("="+encoded)
	} else {
		// Output in hexadecimal format
		fmt.Printf("%06X\n", checksum)
	}
}