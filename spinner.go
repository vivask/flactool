package main

import (
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

//create new spinner
func NewSpinner() {
	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("magenta", "bold")
}

//start animation progress
func StartSpinner() {
	s.Start()
}

//stop animation progress
func StopSpinner() {
	s.Stop()
}
