package main

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func executor(watcher *fsnotify.Watcher, command string) {
	for {
		select {
		case event := <-watcher.Events:
			log.Println("Triggered event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				matched, _ := regexp.MatchString("[^%]%s", command)
				var script string
				if matched {
					script = fmt.Sprintf(command, event.Name)
				} else {
					script = command
				}
				cmd := exec.Command("sh", "-c", script)
				log.Println("Running command:", script)
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Println("Error running command:", string(output))
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func watch(dir string, command string) {
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
	command := os.Args[2]
	watch(dir, command)
}
