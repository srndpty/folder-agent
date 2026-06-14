# folder-agent

ローカルフォルダを監視し、ファイルシステムのイベント（作成・書き込み・削除・リネーム・属性変更）を標準出力に表示する CLI ツールです。

## 必要環境

- Go 1.21 以上

## ビルド

```sh
go build -o folder-agent ./cmd/folder-agent
```

## 実行方法

`--path` で監視対象のフォルダを指定します（必須）。

```sh
# ビルド済みバイナリを実行
./folder-agent --path /path/to/watch

# あるいは go run で直接実行
go run ./cmd/folder-agent --path /path/to/watch
```

実行すると、対象フォルダ内のイベントが標準出力に表示されます。

```
[create] /path/to/watch/example.txt
[write]  /path/to/watch/example.txt
[remove] /path/to/watch/example.txt
```

- イベント本体は **標準出力 (stdout)** に出力されます。
- 起動・終了・エラーなどのログは **標準エラー出力 (stderr)** に `log/slog` 形式で出力されます。

## 終了方法

`Ctrl+C`（SIGINT）または SIGTERM を送ると graceful shutdown します。

## 制限事項（現時点の MVP）

以下は未対応です。

- 再帰的なサブフォルダ監視
- 永続化（SQLite など）
- HTTP API
- Windows サービス化
- OpenTelemetry
- 設定ファイル（config.yaml）
