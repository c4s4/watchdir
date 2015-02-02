package main

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os"
)

func fsnotify(dir, command) {
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("You must pass directory to monitor and command to run")
		os.Exit(1)
	}
	dir := os.Args[0]
	command := os.Args[1:]
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
