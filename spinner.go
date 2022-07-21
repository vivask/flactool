package main

import (
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

/*var stopSpinner bool
var wait chan struct{} = make(chan struct{})

func spinner(delay time.Duration) {
	for !stopSpinner {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
	fmt.Fprint(os.Stdout, "\r \r")
	wait <- struct{}{}
}

//start animation progress
func StartSpinner() {
	stopSpinner = false
	go spinner(100 * time.Millisecond)
}

//stop animation progress
func StopSpinner() {
	stopSpinner = true
	<-wait
}*/

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
}
