package main

// This program uses fsnotify library that defines following file events:
// CREATE, REMOVE, WRITE, RENAME and CHMOD. See sources at
// https://github.com/go-fsnotify/fsnotify/blob/master/fsnotify.go#L35

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

const (
	DEFAULT_CONFIG = "/etc/watchdir.yml"
	HELP           = `watchdir [config]
config   Configuration file (defaults to '/etc/watchdir.yml')`
)

type Configuration map[Directory]Events

type Directory string

type Events map[Event]Command

type Command string

type Event string

func (e Event) Op() fsnotify.Op {
	switch e {
	case "CREATE":
		return fsnotify.Create
	case "WRITE":
		return fsnotify.Write
	case "REMOVE":
		return fsnotify.Remove
	case "RENAME":
		return fsnotify.Rename
	case "CHMOD":
		return fsnotify.Chmod
	default:
		panic(fmt.Sprintf("Unknown event '%s'", e))
	}
}

func executor(watcher *fsnotify.Watcher, events Events) {
	r := regexp.MustCompile("%(f|e)")
	for {
		select {
		case event := <-watcher.Events:
			log.Println("Triggered event:", event)
			for e, command := range events {
				if event.Op&e.Op() == e.Op() {
					cmd := r.ReplaceAllStringFunc(string(command), func(s string) string {
						switch s {
						case "%e":
							return string(e)
						case "%f":
							return event.Name
						default:
							return s
						}
					})
					c := exec.Command("sh", "-c", cmd)
					log.Println("Running command:", cmd)
					output, err := c.CombinedOutput()
					if err != nil {
						log.Println("ERROR running command:", string(output))
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("ERROR:", err)
		}
	}
}

func watch(dir Directory, events Events) {
	log.Println("Watching directory", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go executor(watcher, events)
	err = watcher.Add(string(dir))
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func loadConfig(file string) Configuration {
	config := make(Configuration)
	source, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("ERROR loading configuration file:", err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatal("ERROR parsing configuration file:", err)
	}
	return config
}

func main() {
	configFile := DEFAULT_CONFIG
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Println("ERROR: you may pass only one configuration file on command line")
		fmt.Println(HELP)
		os.Exit(1)
	}
	configuration := loadConfig(configFile)
	for dir, events := range configuration {
		watch(dir, events)
	}
}
