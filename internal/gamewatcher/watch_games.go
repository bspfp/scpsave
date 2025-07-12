package gamewatcher

import (
	"context"
	"log"
	"scpsave/internal/config"
	"scpsave/internal/savesync"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

const refreshInterval = 5 * time.Second

type gameState int

const (
	gameStateNotRunning gameState = iota
	gameStateRunning
)

func StartWatchGames(ctx context.Context) {
	if config.Value.WatchTargetCount < 1 {
		log.Println("No game detected for execution. Shutting down.")
		return
	}

	gameStates := make(map[string]gameState, len(config.Value.Games))

	log.Println("Starting game execution detection.")
	for {
		time.Sleep(refreshInterval)

		procnames := getProcessNames(ctx)
		if len(procnames) == 0 {
			continue
		}

		for _, game := range config.Value.Games {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if game.ProgramName == "" {
				continue
			}

			var newstate gameState
			if _, exists := procnames[game.ProgramName]; exists {
				newstate = gameStateRunning
			} else {
				newstate = gameStateNotRunning
			}
			if gameStates[game.Name] == newstate {
				continue
			}

			gameStates[game.Name] = newstate
			if newstate == gameStateRunning {
				log.Printf("[%s] game started\n", game.Name)
			} else {
				log.Printf("[%s] game stopped\n", game.Name)

				if err := savesync.SyncGame(ctx, game, true); err == nil {
					log.Printf("[%s] synced game\n", game.Name)
				} else {
					log.Printf("[%s] failed to sync game: %+v\n", game.Name, err)
				}
			}
		}
	}
}

func getProcessNames(ctx context.Context) map[string]struct{} {
	procs, err := process.Processes()
	if err != nil {
		log.Printf("failed to get process list: %+v\n", err)
		select {
		case <-ctx.Done():
			return nil
		default:
			return nil
		}
	}

	procnames := make(map[string]struct{}, len(procs)*2)
	for _, proc := range procs {
		if name, err := proc.Name(); err == nil && name != "" {
			procnames[strings.ToLower(name)] = struct{}{}
		}
		if name, err := proc.Exe(); err == nil && name != "" {
			procnames[strings.ToLower(name)] = struct{}{}
		}
	}
	return procnames
}
