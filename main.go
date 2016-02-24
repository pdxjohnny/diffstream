package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pdxjohnny/diffstream/diff"
)

func main() {
	// Print usage
	if len(os.Args) < 2 {
		fmt.Println("Usage: tail -f somefile |", os.Args[0], "comparefile")
		return
	}

	// Open the file to check against
	firstReader, err := os.Open(os.Args[1])
	defer firstReader.Close()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	// The second is the input so we can stream lots of new data in
	secondReader := os.Stdin

	// Set up the diff
	diff := &diff.Diff{
		First:  firstReader,
		Second: secondReader,
		Output: os.Stdout,
	}

	// Start the diff
	diff.Start()

	// Now write the rest of the second to output
	io.Copy(diff.Output, diff.Second)
}
