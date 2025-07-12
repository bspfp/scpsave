package gamewatcher_test

import (
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func TestWatch(t *testing.T) {
	procs, err := process.Processes()
	if err != nil {
		t.Fatal(err)
	}
	for _, proc := range procs {
		name, err := proc.Name()
		if err != nil {
			t.Log("no name")
		} else {
			t.Log(name)
		}
		exe, err := proc.Exe()
		if err != nil {
			t.Log("no exe")
		} else {
			t.Log(exe)
		}
	}
}
