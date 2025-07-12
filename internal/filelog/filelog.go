package filelog

import (
	"fmt"
	"io"
	"log"
	"os"
)

func SetFileLog() (func(), error) {
	f, err := os.Create("./scpsave.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags)
	return func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close log file: %v", err)
		}
	}, nil
}
