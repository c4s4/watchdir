package main

import "testing"

func TestProcessCommand(t *testing.T) {
	c := "%f %e %% %x"
	a := processCommand(c, "file", "event")
	e := "file event % %x"
	if a != e {
		t.Error("Bad processed command")
	}
}
