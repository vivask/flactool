package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

//create new spinner
func NewSpinner() {
	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}

//start animation progress
func StartSpinner() {
	s.Start()
}

//stop animation progress
func StopSpinner() {
	s.Stop()
	fmt.Println()
}
