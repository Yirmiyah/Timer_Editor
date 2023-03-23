package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer watcher.Close()

	err = filepath.Walk("/home/omega", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println("Directory added:", path)
		} else {
			err = watcher.Add(filepath.Dir(path))
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println("Directory added:", filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("File modified:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watch error:", err)
		}
	}
}
