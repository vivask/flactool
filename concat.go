package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

//concatenation of flac, ape, wav files by directories with conversion to flac
func ConcatFlacs(sox, dir string, parallel uint, remove, verbose bool) error {
	list, err := getFilesFromDir(dir, ".flac")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	if len(list) == 0 {
		return fmt.Errorf("flac files not found")
	}

	pathes, keys := prepareFiles(list, false)

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
			//err, out, errout := Shellout(cmd)
			//execVerbose(err, out, errout, verbose)
			if err == nil {
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
				//move result to parent dir
				/*cmd = fmt.Sprintf("mv %s %s", quotes(out), quotes(getParentPath(out)))
				cmdVerbose(cmd, verbose)
				err, out, errout := Shellout(cmd)
				execVerbose(err, out, errout, verbose)*/
			}
			return err
		})
	}
	err = g.Wait()
	if err != nil {
		return err
	}
	//move result to parent dir
	for _, path := range keys {
		out := fmt.Sprintf("%s/%s.flac", path, getLastDir(path))
		cmd := fmt.Sprintf("mv %s %s", quotes(out), quotes(getParentPath(out)))
		cmdVerbose(cmd, verbose)
		err, out, errout := Shellout(cmd)
		execVerbose(err, out, errout, verbose)
	}
	return nil
}

// return "src"
func quotes(src string) string {
	return fmt.Sprintf("\"%s\"", src)
}
