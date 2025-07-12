package savesync

import (
	"context"
	"fmt"
	"log"
	"scpsave/internal/config"
	"scpsave/internal/filelist"
	"scpsave/internal/scp"
)

func remoteToLocal(
	ctx context.Context,
	game *config.GameConfig,
	scpclient *scp.Client,
	remote filelist.FileList,
	mine filelist.FileList,
) error {
	updated, removed := remote.Diff(mine)

	log.Printf("[%s] start downloading...\n", game.Name)

	downloaded := make([]string, 0, len(updated)*2)
	for relPath, metadata := range updated {
		localDownload := game.LocalFileDownloadPath(relPath)
		log.Printf("[%s] downloading file %s\n", game.Name, relPath)
		if err := scpclient.DownloadFile(ctx, game.RemoteFilePath(relPath), localDownload, metadata.ModifiedTime); err != nil {
			return fmt.Errorf("[%s] failed to download file %s: %w", game.Name, relPath, err)
		}
		downloaded = append(downloaded, localDownload, game.LocalFilePath(relPath))
	}

	for i := 0; i < len(downloaded); i += 2 {
		if err := scp.MoveLocalFile(downloaded[i], downloaded[i+1]); err != nil {
			return fmt.Errorf("[%s] failed to move local file %s: %w", game.Name, downloaded[i], err)
		}
	}

	for relPath := range removed {
		log.Printf("[%s] deleting local file %s\n", game.Name, relPath)
		if err := scp.DeleteLocalFile(game.LocalFilePath(relPath)); err != nil {
			return fmt.Errorf("[%s] failed to delete local file %s: %w", game.Name, relPath, err)
		}
	}

	baseFilePath := game.BaseMetaFilePath()
	if err := remote.Save(baseFilePath); err != nil {
		return fmt.Errorf("[%s] failed to save metadata for %s: %w", game.Name, baseFilePath, err)
	}

	return nil
}
