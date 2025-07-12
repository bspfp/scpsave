package conio

import (
	"fmt"
	"scpsave/internal/config"
)

type ResolveMethod int

const (
	LocalToRemote ResolveMethod = iota // 업로드
	RemoteToLocal                      // 다운로드
	Abort                              // 중단
)

func ResolveConflict(game *config.GameConfig) ResolveMethod {
	mu.Lock()
	defer mu.Unlock()

	for range 3 {
		fmt.Printf("[%s] Save files conflict has occurred.\n", game.Name)
		fmt.Println("(L) Load 'Local' file")
		fmt.Println("(R) Load 'Remote' file")
		fmt.Println("(A) Abort")
		fmt.Print("Please choose which file to use: [L or R or A]: ")

		var choice string
		fmt.Scanln(&choice)
		fmt.Println()

		switch choice {
		case "l", "L":
			return LocalToRemote
		case "r", "R":
			return RemoteToLocal
		case "a", "A":
			return Abort
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

	fmt.Println("Conflict resolution failed after 3 attempts.")
	fmt.Println("Aborting resolve process...")
	return Abort
}
