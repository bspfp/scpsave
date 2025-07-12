package savesync

import (
	"context"
	"errors"
	"fmt"
	"log"
	"scpsave/internal/config"
	"scpsave/internal/conio"
	"scpsave/internal/filelist"
	"scpsave/internal/scp"
	"time"
)

func SyncGame(ctx context.Context, game *config.GameConfig, skipDownloadMeta bool) error {
	log.Println("Syncing game:", game.Name)

	scpclient := scp.ClientFromContext(ctx)

	base, err := filelist.LoadFileList(game.BaseMetaFilePath())
	if err != nil {
		return fmt.Errorf("[%s] failed to load base file list: %w", game.Name, err)
	}

	mine, err := filelist.MakeFileList(game.LocalDir, game.FileRegExp)
	if err != nil {
		return fmt.Errorf("[%s] failed to make local file list: %w", game.Name, err)
	}

	var remote filelist.FileList
	if skipDownloadMeta {
		remote, err = filelist.LoadFileList(game.RemoteMetaFileLocalPath())
		if err != nil {
			return fmt.Errorf("[%s] failed to load remote file list: %w", game.Name, err)
		}
	} else {
		err = scpclient.DownloadFile(ctx, game.RemoteMetaFileRemotePath(), game.RemoteMetaFileLocalPath(), time.Now().UnixNano())
		if err == nil {
			remote, err = filelist.LoadFileList(game.RemoteMetaFileLocalPath())
			if err != nil {
				return fmt.Errorf("[%s] failed to load remote file list: %w", game.Name, err)
			}
		} else if !errors.Is(err, scp.ErrNoSuchFile) {
			return fmt.Errorf("[%s] failed to download remote file list: %w", game.Name, err)
		}
	}

	if base.Equal(mine) {
		if base.Equal(remote) {
			// 아무것도 안함
			return nil
		}

		return remoteToLocal(ctx, game, scpclient, remote, mine)
	} else {
		if base.Equal(remote) {
			return localToRemote(ctx, game, scpclient, mine, remote)
		}

		// 충돌 해결 해야 함
		switch conio.ResolveConflict(game) {
		case conio.LocalToRemote:
			return localToRemote(ctx, game, scpclient, mine, remote)

		case conio.RemoteToLocal:
			return remoteToLocal(ctx, game, scpclient, remote, mine)

		default:
			return fmt.Errorf("[%s] conflict resolution aborted by user", game.Name)
		}
	}
}
