package savesync

import (
	"context"
	"errors"
	"log"
	"runtime"
	"scpsave/internal/config"
	"scpsave/internal/sem"
	"sync"
)

func SyncAll(ctx context.Context) error {
	log.Println("Starting synchronization of all games...")

	var errch = make(chan error, len(config.Value.Games))
	defer func() {
		if errch != nil {
			close(errch)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(config.Value.Games))

	sem := sem.NewSemaphore(runtime.NumCPU())
	for _, game := range config.Value.Games {
		sem.Acquire()
		go func(game *config.GameConfig) {
			defer wg.Done()
			defer sem.Release()

			if err := SyncGame(ctx, game, false); err != nil {
				errch <- err
			}
		}(game)
	}
	wg.Wait()

	close(errch)
	var errs []error
	for err := range errch {
		errs = append(errs, err)
	}
	errch = nil
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	log.Println("All games synced successfully.")
	return nil
}
