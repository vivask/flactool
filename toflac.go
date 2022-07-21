package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func FileToFlac(shntool, input string, verbose bool) {
	task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, input, filepath.Dir(input))
	err, out, errout := Shellout(task)
	if verbose {
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		fmt.Println("--- stdout ---")
		fmt.Println(out)
		fmt.Println("--- stderr ---")
		fmt.Println(errout)
	}
}

func DirToFlac(shntool, dir string, parallel uint, outnum, concat, rename, remove, verbose bool) (err error) {
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
			for i, file := range pathes[path] {
				i := i + 1
				input := fmt.Sprintf("%s/%s", path, file)
				newName := input
				if rename {
					newName = fmt.Sprintf("%s/%04d.ape", path, i+1)
					err := os.Rename(input, newName)
					if err != nil {
						return fmt.Errorf("rename error: %w", err)
					}
				}
				g.Go(func() error {
					task := fmt.Sprintf("%s conv -o flac \"%s\" -d \"%s\"", shntool, newName, path)
					err, out, errout := Shellout(task)
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
						if remove {
							err = os.Remove(newName)
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
