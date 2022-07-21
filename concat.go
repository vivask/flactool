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
		if verbose {
			fmt.Println()
			fmt.Println(cmd)
			fmt.Println()
		}
		g.Go(func() error {
			/*err, out, errout := Shellout(cmd)
			if verbose {
				if err != nil {
					log.Printf("error: %v\n", err)
				}
				fmt.Println("--- stdout ---")
				fmt.Println(out)
				fmt.Println("--- stderr ---")
				fmt.Println(errout)
			}*/
			if err == nil {
				//remove source files
				for _, file := range files {
					if remove {
						r := fmt.Sprintf("%s/%s", path, file)
						err = os.Remove(r)
						if err != nil {
							return err
						}
					}
				}
				//move result to parent dir
				fmt.Printf("Src: %s, Dst: %s\n", out, getParentPath(out))
				fmt.Printf("QSrc: %s, QDst: %s\n", quotes(out), quotes(getParentPath(out)))
				err = os.Rename(quotes(out), quotes(getParentPath(out)))
				if err != nil {
					return err
				}
			}
			return err
		})
	}
	err = g.Wait()
	StopSpinner()
	return err
}

// return "src"
func quotes(src string) string {
	return fmt.Sprintf("\"%s\"", src)
}
