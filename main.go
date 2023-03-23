package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const ActivityTime = 1 //in minutes
const rootFolder = "/home/omega"

var TotalTimeArray []int64

func main() {
	stop := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer watcher.Close()

	err = RecursivityAdd(watcher, rootFolder)
	if err != nil {
		fmt.Println(err)
		return
	}
	SetStartingTime()
	UpdateLastModified()
	go printElapsedTime(Starttime, &lastModified, stop)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				
				if !CheckIfModifiedEachN(ActivityTime) {
					fmt.Println("You Haven't Modified Anything For", strconv.Itoa(ActivityTime), "Minutes")
					TotalTimeArray = append(TotalTimeArray, lastModified-Starttime)
					SetStartingTime()
					UpdateLastModified()
					stop <- true
					stop = make(chan bool)
					go printElapsedTime(Starttime, &lastModified, stop)
				} else {
					UpdateLastModified()
					stop <- true
					stop = make(chan bool)
					go printElapsedTime(Starttime, &lastModified, stop)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					//check if it's a file or a directory
					file, err := os.Stat(event.Name)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if file.IsDir() {
						if !IsAlreadyWatched(watcher, event.Name) {
							RecursivityAdd(watcher, event.Name)
						}
					} else {
						if !IsAlreadyWatched(watcher, filepath.Dir(event.Name)) {
							watcher.Add(filepath.Dir(event.Name))
						}
					}
					fmt.Println("File or Folder created:", event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("File removed:", event.Name)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("File renamed:", event.Name)
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("File chmod:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Watch error:", err)
			}
		}
	}()
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			total := int64(0)
			for _, time := range TotalTimeArray {
				total += time
			}
			fmt.Println("Total Time:", time.Duration(total).Minutes())

		}
	}
}
func printElapsedTime(Starttime int64, lastModified *int64, stop chan bool) {
	//start is equal to a pointer of the Starttime variable

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		start := time.Unix(0, Starttime)
		lastModif := time.Unix(0, *lastModified)
		select {
		case <-ticker.C:
			(lastModif)
			elapsed := time.Since(start)
			fmt.Printf("Elapsed time from start: %v\n", elapsed)
			elapsed = time.Since(lastModif)
			fmt.Printf("Elapsed time from last modification: %v\n", elapsed)
		case <-stop:
			return
		}
	}
}
func IsAlreadyWatched(watcher *fsnotify.Watcher, path string) bool {
	for _, watch := range watcher.WatchList() {
		if watch == path {
			return true
		}
	}
	return false
}
func RecursivityAdd(watcher *fsnotify.Watcher, root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if info.IsDir() && !IsAlreadyWatched(watcher, path) && !FolderExcluded(path) {
			err = watcher.Add(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println("Directory added:", path)
		} else if !IsAlreadyWatched(watcher, filepath.Dir(path)) && !FolderExcluded(path) {
			err = watcher.Add(filepath.Dir(path))
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println("Directory added:", filepath.Dir(path))
		}
		return nil
	})
	return err
}
func FolderExcluded(path string) bool {
	//Exclude all node_modules, hidden folders and files using contains
	//Remove all the path except the last folder

	excluded := []string{".", "node_modules", ".git", ".vscode", ".idea", ".cache", ".config", ".local", ".npm", ".npmrc", ".yarn", ".yarnrc", ".yarnrc.yml", ".yarn-integrity", ".yarn-metadata.json", ".yarn"}
	for _, ex := range excluded {
		if strings.Contains(path, ex) {
			return true
		}
	}
	return false
}

var Starttime int64
var lastModified int64

func SetStartingTime() {
	Starttime = time.Now().UnixNano()
}
func UpdateLastModified() {
	lastModified = time.Now().UnixNano()
}
func CheckIfModifiedEachN(n int64) bool {
	//convert n from minutes to nanoseconds
	n = n * 60 * 1000000000
	if Starttime == 0 {
		SetStartingTime()
		return true
	}
	if lastModified == 0 {
		UpdateLastModified()
		return true
	}
	(time.Now().UnixNano()-lastModified, "------------", n)
	if time.Now().UnixNano()-lastModified < n {
		UpdateLastModified()
		return true
	}
	return false
}
