package diff

import (
	"bufio"
	"fmt"
	"io"
)

// Diff takes two inputs and if a line in the second matches a line in the
// first it is not output. If a line in the second does not match a line in the
// first it is sent to output
type Diff struct {
	First  io.Reader
	Second io.Reader
	Output io.Writer
}

// Start reading from First and Second and put unique output in Output
func (diff *Diff) Start() error {
	// These are the lines from each reader that will be compared
	var firstCurrentLine string
	var secondCurrentLine string

	// Pass the line from the reader to the compitor
	firstRecvLine := make(chan string, 1)
	secondRecvLine := make(chan string, 1)

	// When you are done compareing send thorugh this to ask for the next line
	firstGoOn := make(chan bool, 1)
	secondGoOn := make(chan bool, 1)

	// We need to know when we have hit the end of each reader
	firstDone := make(chan error, 1)
	secondDone := make(chan error, 1)

	// Create scanners to read the lines
	firstScanner := bufio.NewScanner(diff.First)
	secondScanner := bufio.NewScanner(diff.Second)

	go diff.OnLine(firstScanner, firstRecvLine, firstGoOn, firstDone)
	go diff.OnLine(secondScanner, secondRecvLine, secondGoOn, secondDone)

	// Get both lines and compare them until the first runs out
	hasBoth := make(chan bool, 1)
	hasFirst := false
	hasSecond := false

	// If either finihes reading then return
	for {
		select {
		case err := <-firstDone:
			// If first is done print the output of what we read from second
			if hasSecond {
				fmt.Fprintln(diff.Output, secondCurrentLine)
			}
			// Tell second to stop scanning
			secondGoOn <- false
			<-secondDone
			// Send the rest of second to output
			for secondScanner.Scan() {
				fmt.Fprintln(diff.Output, secondScanner.Text())
			}
			return err
		case err := <-secondDone:
			// If second is done stop scanning first
			firstGoOn <- false
			<-firstDone
			return err
		case firstCurrentLine = <-firstRecvLine:
			hasFirst = true
			if hasFirst && hasSecond {
				hasBoth <- true
			}
		case secondCurrentLine = <-secondRecvLine:
			hasSecond = true
			if hasFirst && hasSecond {
				hasBoth <- true
			}
		case <-hasBoth:
			if secondCurrentLine != firstCurrentLine {
				fmt.Fprintln(diff.Output, secondCurrentLine)
			}

			firstGoOn <- true
			secondGoOn <- true

			hasFirst = false
			hasSecond = false
		}
	}
}

func (diff *Diff) OnLine(scanner *bufio.Scanner, sendLine chan string, goOn chan bool, done chan error) {
	// So long as you can scan keep doing it
	for scanner.Scan() {
		// Every time you scan send the line back
		sendLine <- scanner.Text()
		// So long as you've been asked to goOn keep scaning
		keepScanning := <-goOn
		if !keepScanning {
			break
		}
	}
	// When you can scan no more let us know if it was because of an error
	done <- scanner.Err()
}
