package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

//split flac, ape, wav files according to cue by directories
func SplitApeOrFlac(shntool, cuetag, dir string, parallel uint, remove, verbose bool) error {
	files, err := getSplitFilesFromDir(dir, ".flac", ".ape", ".wav")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("files not found")
	}

	for _, file := range files {
		fmt.Println(file)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Split this files?")
	fmt.Println("[Yes/no]?")
	t, _ := reader.ReadString('\n')
	t = strings.ToLower(t)
	if t == "y\n" || t == "yes\n" || t == "\n" {
		StartSpinner()
		g, _ := errgroup.WithContext(context.Background())
		g.SetLimit(int(parallel))
		for _, file := range files {
			file := file
			cue := replaceExtToCue(file)
			out := filepath.Dir(file)
			cmd := shntool + " split -f \"" + cue + "\" -o flac -t \"%n %t\" " + "\"" + file + "\" -d " + "\"" + out + "\""
			cmdVerbose(cmd, verbose)
			g.Go(func() error {
				err, out, errout := Shellout(cmd)
				execVerbose(err, out, errout, verbose)
				if err == nil {
					if remove {
						err = os.Remove(file)
						if err != nil {
							return err
						}
					}
					cmd := fmt.Sprintf("%s \"%s\" \"%s/*%s\"", cuetag, cue, filepath.Dir(file), filepath.Ext(file))
					cmdVerbose(cmd, verbose)
					err, out, errout := Shellout(cmd)
					execVerbose(err, out, errout, verbose)
				}
				return err
			})
		}
		err = g.Wait()
		StopSpinner()
		return err
	}

	return nil
}
