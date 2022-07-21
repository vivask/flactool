package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

//concatenation of flac files by directories
func ConcatFlacs(sox, dir string, dirs []string, parallel uint, remove, verbose bool) error {

	var pathes FileList
	var keys []string
	var err error
	if len(dirs) == 0 {
		list, err := getFilesFromDir(dir, ".flac")
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		if len(list) == 0 {
			return fmt.Errorf("flac files not found")
		}
		pathes, keys = prepareFiles(list, false)
	} else {
		pathes, keys, err = getFilesFromDirs(dirs, ".flac")
		if err != nil {
			return err
		}
	}

	for _, path := range pathes {
		fmt.Println(path)
	}

	StartSpinner()
	defer StopSpinner()
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(int(parallel))
	for _, path := range keys {
		path := path
		files := pathes[path]
		cmd := fmt.Sprintf("%s -S", sox)

		//create list input files
		input := ""
		for _, file := range files {
			input = fmt.Sprintf("%s \"%s/%s\" ", input, path, file)
		}

		out := fmt.Sprintf("%s/%s.flac", path, getLastDir(path))
		cmd = fmt.Sprintf("%s %s%s", cmd, input, quotes(out))
		cmdVerbose(cmd, verbose)
		g.Go(func() error {
			//concat shntool
			err, stdout, errout := Shellout(cmd)
			execVerbose(err, stdout, errout, verbose)
			if err != nil {
				return err
			}

			//move result to parent dir
			cmd = fmt.Sprintf("mv %s %s", quotes(out), quotes(getParentPath(out)))
			cmdVerbose(cmd, verbose)
			err, stdout, errout = Shellout(cmd)
			execVerbose(err, stdout, errout, verbose)
			if err != nil {
				return err
			}

			//remove source files
			if remove {
				for _, file := range files {
					r := fmt.Sprintf("%s/%s", path, file)
					err = os.Remove(r)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})
	}
	return g.Wait()
}

// return "src"
func quotes(src string) string {
	return fmt.Sprintf("\"%s\"", src)
}
