package main

import (
	"errors"
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"sync"
)

var VERSION = "UNKNOWN"
var HELP = `watchdir [config(s)]
config Configuration file(s) (defaults to '~/.watchdir.yml' or '/etc/watchdir.yml')`
var USER_CONFIG = "~/.watchdir.yml"
var SYS_CONFIG = "/etc/watchdir.yml"
var REGEXP = regexp.MustCompile("%(f|e|%)")

// Configuration is a map that gives Events for a directory
type Configuration map[string]Events

// Events is a map that gives a command for an event
type Events map[string]string

// EventCode return a code of a given event.
func EventCode(s string) fsnotify.Op {
	switch s {
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
		panic(fmt.Sprintf("Unknown event '%s'", s))
	}
}

// ExpandCommand replace '%f' with file name and '%e' with event in a given
// command.
func ExpandCommand(command, file, event string) string {
	return REGEXP.ReplaceAllStringFunc(command, func(s string) string {
		switch s {
		case "%f":
			return file
		case "%e":
			return event
		case "%%":
			return "%"
		default:
			return s
		}
	})
}

// Execute a given command, if an error occurs, it is logged.
func ExecuteCommand(command, file, event string) {
	expanded := ExpandCommand(command, file, event)
	log.Println("Running command:", expanded)
	c := exec.Command("sh", "-c", expanded)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Println("Error running command:", string(output))
	}
}

// WatchDirectory fires commands when specific events are triggered in a given
// directory.
func WatchDirectory(directory string, events Events, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	log.Println("Watching directory", directory)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Println("Directory", directory, "not found")
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer watcher.Close()
	err = watcher.Add(directory)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for {
		select {
		case event := <-watcher.Events:
			for e, command := range events {
				if event.Op&EventCode(e) == EventCode(e) {
					log.Println("Triggered event", event)
					ExecuteCommand(command, event.Name, e)
				}
			}
		case err := <-watcher.Errors:
			log.Println(err.Error())
		}
	}
}

// LoadConfiguration loads configuration from file.
func LoadConfiguration(file string) (Configuration, error) {
	log.Println("Loading configuration file", file)
	configuration := make(Configuration)
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return configuration, err
	}
	err = yaml.Unmarshal(source, &configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, nil
}

// ExpandUser expand file names starting with '~/' by the user directory.
// Thus '~/.watchdir.yml' would be expanded to '/home/foo/.watchdir.yml' for
// user foo.
func ExpandUser(file string) string {
	if file[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		return strings.Replace(file, "~", dir, 1)
	} else {
		return file
	}
}

// ConfigFile return user expanded configuration file:
// - Configuration file passed on command line if any.
// - User sonfiguration file '~/.watchdir.yml' if it exists.
// - System configuration file '/etc/watchdir.yml' if it exists.
// - Panic if none is found.
func ConfigFile() (string, error) {
	if len(os.Args) == 2 {
		return ExpandUser(os.Args[1]), nil
	} else if len(os.Args) < 2 {
		if _, err := os.Stat(USER_CONFIG); err == nil {
			return ExpandUser(USER_CONFIG), nil
		}
		if _, err := os.Stat(SYS_CONFIG); err == nil {
			return SYS_CONFIG, nil
		}
		return "", errors.New("Configuration file not found")
	} else {
		return "", errors.New("You can pass one configuration file on command line")
	}
}

func main() {
	log.Println("Starting watchdir version", VERSION)
	configurationFile, err := ConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	configuration, err := LoadConfiguration(configurationFile)
	if err != nil {
		log.Fatal(err)
	}
	waitGroup := &sync.WaitGroup{}
	for directory, events := range configuration {
		waitGroup.Add(1)
		go WatchDirectory(ExpandUser(directory), events, waitGroup)
	}
	waitGroup.Wait()
}
