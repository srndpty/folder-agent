// Command folder-agent は指定したローカルフォルダを監視し、
// ファイルシステムイベントを標準出力に表示する CLI ツールです。
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	path := flag.String("path", "", "監視対象のフォルダパス（必須）")
	flag.Parse()

	if err := run(logger, *path); err != nil {
		logger.Error("folder-agent を終了しました", slog.Any("error", err))
		os.Exit(1)
	}
}

// run は監視のセットアップとイベントループを行う。
// Ctrl+C（SIGINT）や SIGTERM を受け取ると graceful shutdown する。
func run(logger *slog.Logger, path string) error {
	if path == "" {
		return errors.New("--path は必須です")
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("--path はフォルダを指定してください: " + path)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return err
	}

	// SIGINT / SIGTERM を受け取るとコンテキストがキャンセルされる。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("監視を開始しました", slog.String("path", path))

	for {
		select {
		case <-ctx.Done():
			logger.Info("シャットダウンしています")
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			logEvent(event)

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			logger.Error("監視エラーが発生しました", slog.Any("error", err))
		}
	}
}

// logEvent はイベント種別を判定して標準出力に表示する。
// イベント本体は人間向けの一覧なので、ログ（stderr）ではなく標準出力へ出す。
func logEvent(event fsnotify.Event) {
	fmt.Printf("[%s] %s\n", eventKind(event), event.Name)
}

// eventKind は fsnotify のイベントを create/write/remove/rename/chmod に分類する。
func eventKind(event fsnotify.Event) string {
	switch {
	case event.Op.Has(fsnotify.Create):
		return "create"
	case event.Op.Has(fsnotify.Write):
		return "write"
	case event.Op.Has(fsnotify.Remove):
		return "remove"
	case event.Op.Has(fsnotify.Rename):
		return "rename"
	case event.Op.Has(fsnotify.Chmod):
		return "chmod"
	default:
		return "unknown"
	}
}
