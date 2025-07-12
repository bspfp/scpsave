package conio_test

import (
	"os"
	"scpsave/internal/config"
	"scpsave/internal/conio"
	"testing"
)

func TestConsoleIO(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	w.WriteString("h\nA\n")
	w.Close()
	os.Stdin = r

	switch conio.ResolveConflict(&config.GameConfig{Name: "Test Game"}) {
	case conio.LocalToRemote:
		t.Log("Local to Remote")
	case conio.RemoteToLocal:
		t.Log("Remote to Local")
	case conio.Abort:
		t.Log("Abort")
	}
}
