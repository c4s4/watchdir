package main

import (
	"os/user"
	"testing"
)

func TestExpandCommand(t *testing.T) {
	e := "file event % %x"
	a := ExpandCommand("%f %e %% %x", "file", "event")
	if e != a {
		t.Error("Bad processed command")
	}
}

func TestExpandUser(t *testing.T) {
	user, _ := user.Current()
	e := "/home/" + user.Username + "/foo"
	f := "/Users/" + user.Username + "/foo"
	a := ExpandUser("~/foo")
	if e != a && f != a {
		t.Error("User directory not expanded as expected", e, "!=", a)
	}
}
