package savesync

import (
	"context"
	"fmt"
	"log"
	"scpsave/internal/config"
	"scpsave/internal/filelist"
	"scpsave/internal/scp"
)

func localToRemote(
	ctx context.Context,
	game *config.GameConfig,
	scpclient *scp.Client,
	mine filelist.FileList,
	remote filelist.FileList,
) error {
	updated, removed := mine.Diff(remote)

	log.Printf("[%s] start uploading...\n", game.Name)

	uploaded := make([]string, 0, len(updated)*2)
	for relPath := range updated {
		remoteUpload := game.RemoteFileUploadPath(relPath)
		log.Printf("[%s] uploading file %s...", game.Name, relPath)
		if err := scpclient.UploadFile(ctx, game.LocalFilePath(relPath), remoteUpload); err != nil {
			return fmt.Errorf("[%s] failed to upload file %s: %w", game.Name, relPath, err)
		}
		uploaded = append(uploaded, remoteUpload, game.RemoteFilePath(relPath))
	}

	log.Printf("[%s] uploading metadata\n", game.Name)
	remoteMetaLocal := game.RemoteMetaFileLocalPath()
	if err := mine.Save(remoteMetaLocal); err != nil {
		return fmt.Errorf("[%s] failed to save metadata for %s: %w", game.Name, remoteMetaLocal, err)
	}
	remoteUpload := game.RemoteMetaFileUploadPath()
	if err := scpclient.UploadFile(ctx, remoteMetaLocal, remoteUpload); err != nil {
		return fmt.Errorf("[%s] failed to upload metadata for %s: %w", game.Name, remoteMetaLocal, err)
	}
	uploaded = append(uploaded, remoteUpload, game.RemoteMetaFileRemotePath())

	for i := 0; i < len(uploaded); i += 2 {
		if err := scpclient.MoveRemoteFile(uploaded[i], uploaded[i+1]); err != nil {
			return fmt.Errorf("[%s] failed to move remote file %s: %w", game.Name, uploaded[i], err)
		}
	}

	for relPath := range removed {
		log.Printf("[%s] deleting remote file %s\n", game.Name, relPath)
		if err := scpclient.DeleteRemoteFile(game.RemoteFilePath(relPath)); err != nil {
			return fmt.Errorf("[%s] failed to delete remote file %s: %w", game.Name, relPath, err)
		}
	}

	baseFilePath := game.BaseMetaFilePath()
	_ = scp.DeleteLocalFile(baseFilePath)
	if err := scp.MoveLocalFile(remoteMetaLocal, baseFilePath); err != nil {
		return fmt.Errorf("[%s] failed to update remote meta file in working: %w", game.Name, err)
	}

	return nil
}
