package main

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"strings"
)

func executor(watcher *fsnotify.Watcher, command []string) {
	for {
		select {
		case event := <-watcher.Events:
			log.Println("Triggered event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				name := command[0]
				params := append(command[1:], event.Name)
				cmd := exec.Command(name, params...)
				log.Println("Running command:", name, strings.Join(params, " "))
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Println("Error running command:", output)
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func watch(dir string, command []string) {
	log.Println("Watching directory", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go executor(watcher, command)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("You must pass directory to monitor and command to run")
		os.Exit(1)
	}
	dir := os.Args[1]
	command := os.Args[2:]
	watch(dir, command)
}
