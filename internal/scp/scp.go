package scp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"scpsave/internal/gzipio"
	"strings"
	"time"

	goscp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

var (
	ErrNoSuchFile = errors.New("no such file or directory")
)

type Client struct {
	scpClient *goscp.Client
}

func NewClient(serverAddr, username, privateKeyPath string) (*Client, error) {
	cfg, err := auth.PrivateKey(username, privateKeyPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return nil, fmt.Errorf("failed to create auth config: %w", err)
	}
	client := goscp.NewClient(serverAddr, &cfg)
	if err := client.Connect(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to SCP server: %w", err)
	}
	return &Client{scpClient: &client}, nil
}

func (c *Client) Close() {
	if c.scpClient != nil {
		c.scpClient.Close()
		c.scpClient = nil
	}
}

func (c *Client) UploadFile(ctx context.Context, localPath, remotePath string) error {
	localPath = filepath.Clean(localPath)
	r, err := gzipio.NewCompressReader(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrNoSuchFile, localPath)
		}
		return fmt.Errorf("failed to create compress reader for %s: %w", localPath, err)
	}

	if err := c.ensureRemoteDir(remotePath); err != nil {
		return fmt.Errorf("failed to ensure remote directory: %w", err)
	}

	remotedir := path.Dir(remotePath)
	remotename := path.Base(remotePath)
	tempfile := path.Join(remotedir, remotename+".uploading")

	if err := c.scpClient.CopyFile(ctx, r, tempfile, "0644"); err != nil {
		return fmt.Errorf("failed to upload file %s to %s: %w", localPath, remotePath, err)
	}

	_ = c.execRemote(fmt.Sprintf(`rm -f "%s"`, remotePath))
	if err := c.execRemote(fmt.Sprintf(`mv "%s" "%s"`, tempfile, remotePath)); err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", tempfile, remotePath, err)
	}
	return nil
}

func (c *Client) DownloadFile(ctx context.Context, remotePath, localPath string, modTime int64) error {
	localPath, err := filepath.Abs(localPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", localPath, err)
	}

	localPath = filepath.Clean(localPath)
	localdir := filepath.Dir(localPath)
	if err := os.MkdirAll(localdir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory %s: %w", localdir, err)
	}

	localname := filepath.Base(localPath)
	tempfile := filepath.Join(localdir, localname+".download")
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempfile) // Ignore error if file doesn't exist
		}
	}()

	err = func() error {
		w := gzipio.NewDecompressWriter(tempfile)

		if err := c.scpClient.CopyFromRemotePassThru(ctx, w, remotePath, nil); err != nil {
			if strings.Contains(err.Error(), "No such file or directory") {
				return fmt.Errorf("%w: %s", ErrNoSuchFile, remotePath)
			}
			return fmt.Errorf("failed to download file from %s to %s: %w", remotePath, tempfile, err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close decompress writer: %w", err)
		}

		return nil
	}()
	if err != nil {
		return err
	}

	_ = os.Remove(localPath) // Ignore error if file doesn't exist
	if err := os.Rename(tempfile, localPath); err != nil {
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tempfile, localPath, err)
	}

	modTimeTime := time.Unix(0, modTime)
	if err := os.Chtimes(localPath, modTimeTime, modTimeTime); err != nil {
		log.Printf("failed to set modification time for %s: %+v\n", localPath, err)
	}

	success = true
	return nil
}

func (c *Client) DeleteRemoteFile(remotePath string) error {
	if err := c.execRemote(fmt.Sprintf(`rm -f "%s"`, remotePath)); err != nil {
		return fmt.Errorf("failed to delete remote file %s: %w", remotePath, err)
	}
	return nil
}

func (c *Client) MoveRemoteFile(oldRemotePath, newRemotePath string) error {
	_ = c.execRemote(fmt.Sprintf(`rm -f "%s"`, newRemotePath))
	if err := c.ensureRemoteDir(newRemotePath); err != nil {
		return fmt.Errorf("failed to ensure remote directory for %s: %w", newRemotePath, err)
	}
	if err := c.execRemote(fmt.Sprintf(`mv "%s" "%s"`, oldRemotePath, newRemotePath)); err != nil {
		return fmt.Errorf("failed to rename remote file %s: %w", oldRemotePath, err)
	}
	return nil
}

func (c *Client) execRemote(command string) error {
	session, err := c.scpClient.SSHClient().NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()
	return session.Run(command)
}

func (c *Client) ensureRemoteDir(remotePath string) error {
	remoteDir := path.Dir(remotePath)
	if remoteDir == "." || remoteDir == "/" {
		return nil // No need to create root directory
	}

	if err := c.execRemote(fmt.Sprintf(`mkdir -p "%s"`, remoteDir)); err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", remoteDir, err)
	}

	return nil
}

type contextClientKey struct{}

func NewContextWithClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, contextClientKey{}, client)
}

func ClientFromContext(ctx context.Context) *Client {
	if client, ok := ctx.Value(contextClientKey{}).(*Client); ok {
		return client
	}
	return nil
}

func DeleteLocalFile(localPath string) error {
	localPath = filepath.Clean(localPath)
	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("failed to delete local file %s: %w", localPath, err)
	}
	return nil
}

func MoveLocalFile(oldLocalPath, newLocalPath string) error {
	_ = os.Remove(newLocalPath)
	if err := os.Rename(oldLocalPath, newLocalPath); err != nil {
		return fmt.Errorf("failed to move local file %s: %w", oldLocalPath, err)
	}
	return nil
}
