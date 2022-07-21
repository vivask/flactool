package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func FileToFlac(shntool, input string, verbose bool) {
	task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, filepath.Dir(input))
	err, out, errout := Shellout(task)
	execVerbose(err, out, errout, verbose)
}

func DirToFlac(shntool, dir string, parallel uint, concat, remove, verbose bool) (err error) {
	list, err := getFilesFromDir(dir, ".ape", ".wav")
	if len(list) == 0 {
		if !concat {
			return fmt.Errorf("ape or wav files not found")
		}
	} else {

		pathes, keys := prepareFiles(list, true)

		StartSpinner()
		g, _ := errgroup.WithContext(context.Background())
		g.SetLimit(int(parallel))
		for _, path := range keys {
			path := path
			for _, file := range pathes[path] {
				input := fmt.Sprintf("%s/%s", path, file)
				g.Go(func() error {
					task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, path)
					err, stdout, errout := Shellout(task)
					execVerbose(err, stdout, errout, verbose)

					if err == nil {
						if remove {
							err = os.Remove(input)
						}
					}
					return err
				})
			}
		}
		err = g.Wait()
		StopSpinner()
		return err
	}
	return nil
}
