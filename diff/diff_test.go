package diff

import (
	"log"
	"os"
	"strings"
)

func ExampleDiffStart() {
	firstInput := "same\ndifferent"
	secondInput := "same\ndiffer3nt"

	firstReader := strings.NewReader(firstInput)
	secondReader := strings.NewReader(secondInput)

	diff := &Diff{
		First:  firstReader,
		Second: secondReader,
		Output: os.Stdout,
	}

	err := diff.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// differ3nt
}

func ExampleDiffStartFirstFinishesBeforeSecond() {
	firstInput := "same\n"
	secondInput := "same\nsecond still going 1\nsecond still going 2\nsecond still going 3"

	firstReader := strings.NewReader(firstInput)
	secondReader := strings.NewReader(secondInput)

	diff := &Diff{
		First:  firstReader,
		Second: secondReader,
		Output: os.Stdout,
	}

	err := diff.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// second still going 1
	// second still going 2
	// second still going 3
}

func ExampleDiffStartSecondFinishesBeforeFirst() {
	firstInput := "same\nfirst still going 1\nfirst still going 2"
	secondInput := "same\n"

	firstReader := strings.NewReader(firstInput)
	secondReader := strings.NewReader(secondInput)

	diff := &Diff{
		First:  firstReader,
		Second: secondReader,
		Output: os.Stdout,
	}

	err := diff.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	//
}

func ExampleDiffStartVarying() {
	firstInput := "same\ndifferent\nsame"
	secondInput := "same\ndiffer3nt\nsame"

	firstReader := strings.NewReader(firstInput)
	secondReader := strings.NewReader(secondInput)

	diff := &Diff{
		First:  firstReader,
		Second: secondReader,
		Output: os.Stdout,
	}

	err := diff.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// differ3nt
}
