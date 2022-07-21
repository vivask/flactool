package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

//concatenation of flac, ape, wav files by directories with conversion to flac
func ConcatFlacs(sox, dir string, parallel uint, rename, remove, verbose, worked bool) error {
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
		input := ""
		for i, file := range files {
			newName := file
			if rename {
				newName = fmt.Sprintf("%04d.flac", i+1)
				if newName != file {
					err := os.Rename(fmt.Sprintf("%s/%s", path, file), fmt.Sprintf("%s/%s", path, newName))
					if err != nil {
						return fmt.Errorf("rename error: %w", err)
					}
				}
			}
			input = fmt.Sprintf("%s \"%s/%s\" ", input, path, newName)
		}

		cmd = fmt.Sprintf("%s %s\"%s/%s.flac\"", cmd, input, path, getLastDir(path))
		if verbose {
			fmt.Println()
			fmt.Println(cmd)
			fmt.Println()
		}
		g.Go(func() error {
			err, out, errout := Shellout(cmd)
			if verbose {
				if err != nil {
					log.Printf("error: %v\n", err)
				}
				fmt.Println("--- stdout ---")
				fmt.Println(out)
				fmt.Println("--- stderr ---")
				fmt.Println(errout)
			}
			if err == nil {
				for i, file := range files {
					newName := file
					if rename {
						newName = fmt.Sprintf("%04d.flac", i+1)
					}
					if remove {
						r := fmt.Sprintf("%s/%s", path, newName)
						err = os.Remove(r)
						if err != nil {
							return err
						}
					}
				}
			}
			return err
		})
	}
	err = g.Wait()
	StopSpinner()
	return err

}

func getLastDir(path string) string {
	split := strings.Split(path, "/")
	cnt := len(split)
	if cnt > 0 {
		return split[len(split)-1]
	}
	return path
}
