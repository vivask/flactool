package cmd

import (
	"fmt"
	"os"
	"time"
)

var stopSpinner bool
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
}
