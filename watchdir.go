package main

// This program uses fsnotify library that defines following file events:
// CREATE, REMOVE, WRITE, RENAME and CHMOD. See sources at
// https://github.com/go-fsnotify/fsnotify/blob/master/fsnotify.go#L35

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/kylelemons/go-gypsy/yaml"
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

type Events map[fsnotify.Op]string

type Configuration map[string]Events

func eventToString(e fsnotify.Op) string {
	switch e {
	case fsnotify.Create:
		return "CREATE"
	case fsnotify.Write:
		return "WRITE"
	case fsnotify.Remove:
		return "REMOVE"
	case fsnotify.Rename:
		return "RENAME"
	case fsnotify.Chmod:
		return "CHMOD"
	default:
		return ""
	}
}

func nodeToEvents(n yaml.Node) Events {
	e := make(Events)
	for event, node := range nodeToMap(n) {
		command := nodeToString(node)
		switch event {
		case "CREATE":
			e[fsnotify.Create] = command
		case "WRITE":
			e[fsnotify.Write] = command
		case "REMOVE":
			e[fsnotify.Remove] = command
		case "RENAME":
			e[fsnotify.Rename] = command
		case "CHMOD":
			e[fsnotify.Chmod] = command
		default:
			log.Fatal(fmt.Sprintf("ERROR: unknown event '%s'", event))
		}
	}
	return e
}

func executor(watcher *fsnotify.Watcher, events Events) {
	r := regexp.MustCompile("%(f|e)")
	for {
		select {
		case event := <-watcher.Events:
			log.Println("Triggered event:", event)
			for e, command := range events {
				if event.Op&e == e {
					cmd := r.ReplaceAllStringFunc(command, func(s string) string {
						switch s {
						case "%e":
							return eventToString(e)
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

func watch(dir string, events Events) {
	log.Println("Watching directory", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go executor(watcher, events)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func nodeToMap(node yaml.Node) yaml.Map {
	m, ok := node.(yaml.Map)
	if !ok {
		log.Fatal(fmt.Sprintf("ERROR parsing configuration file: %v is not of type map", node))
	}
	return m
}

func nodeToString(node yaml.Node) string {
	s, ok := node.(yaml.Scalar)
	if !ok {
		log.Fatal(fmt.Sprintf("ERROR parsing configuration file: %v is not of type string", node))
	}
	return s.String()
}

func loadConfig(file string) Configuration {
	config := make(Configuration)
	doc, err := yaml.ReadFile(file)
	if err != nil {
		log.Fatal("ERROR parsing configuration file:", err)
	}
	for d, e := range nodeToMap(doc.Root) {
		config[d] = nodeToEvents(e)
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
