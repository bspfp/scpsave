package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"scpsave/internal/config"
	"scpsave/internal/filelog"
	"scpsave/internal/gamewatcher"
	"scpsave/internal/savesync"
	"scpsave/internal/scp"
	"syscall"
)

var (
	flagCreateSampleConfig = flag.Bool("c", false, "Create a sample config file and exit")
)

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	closelog, err := filelog.SetFileLog()
	if err != nil {
		log.Fatalf("Failed to set file log: %+v\n", err)
	}
	defer closelog()

	if *flagCreateSampleConfig {
		if err := config.MakeSampleConfig(); err != nil {
			log.Fatalf("Failed to create sample config: %+v\n", err)
		}
		os.Exit(0)
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %+v\n", err)
	}

	scpclient, err := scp.NewClient(config.Value.ServerAddress, config.Value.Username, config.Value.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to create SCP client: %+v\n", err)
	}
	defer scpclient.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx = scp.NewContextWithClient(ctx, scpclient)

	if err := savesync.SyncAll(ctx); err != nil {
		log.Fatalf("Failed to sync saves: %+v\n", err)
	}

	gamewatcher.StartWatchGames(ctx)

	log.Println("Exiting...")
}
