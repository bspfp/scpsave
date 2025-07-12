package scp_test

import (
	"errors"
	"scpsave/internal/scp"
	"testing"
	"time"
)

func TestNoSuchFile(t *testing.T) {
	client, err := scp.NewClient("192.168.0.96:20001", "bs", `C:\Users\kwang\.ssh\id_rsa`)
	if err != nil {
		t.Fatalf("Failed to create SCP client: %v", err)
	}
	defer client.Close()

	err = client.DownloadFile(t.Context(), "/home/bs/temp/test.txt", `.\test.txt`, time.Now().UnixNano())
	if err != nil {
		if errors.Is(err, scp.ErrNoSuchFile) {
			t.Logf("Expected error: %v", err)
			return
		}
		t.Fatalf("Failed to download file: %v", err)
	}
	t.Logf("File downloaded successfully")

	_ = scp.DeleteLocalFile(`.\test.txt`)
}
