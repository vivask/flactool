package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

func concatFlacs(sox, dir string, parallel uint, rename, remove, verbose, worked bool) error {
	list, err := getFilesFromDir(dir, ".flac")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	if len(list) == 0 {
		return fmt.Errorf("flac files not found")
	}

	pathes := prepareFiles(list)
	for path, files := range pathes {
		fmt.Println(path)
		for _, file := range files {
			fmt.Println(file)
		}
		fmt.Println()
	}

	t := "y\n"
	if !worked {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Concatenate this files?")
		fmt.Println("[Yes/no]?")
		t, _ = reader.ReadString('\n')
		t = strings.ToLower(t)
	}
	if t == "y\n" || t == "yes\n" || t == "\n" {
		StartSpinner()
		g, _ := errgroup.WithContext(context.Background())
		g.SetLimit(int(parallel))
		cdNum := 1
		for path, files := range pathes {
			if len(files) < 2 {
				continue
			}
			path := path
			files := files
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
			cmd = fmt.Sprintf("%s %s\"%s/CD%d.flac\"", cmd, input, path, cdNum)
			if verbose {
				fmt.Println()
				fmt.Println(cmd)
				fmt.Println()
			}
			cdNum++
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

	return nil
}
