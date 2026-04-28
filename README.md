# TUI Todo List

Bubble Teaを使用した、モダンでインタラクティブなGo言語製CLI Todoリストアプリです。

## インストール方法

以下のコマンドを実行すると、ターミナルから golang-todolist-tui コマンドで起動できるようになります。

`ash
go install github.com/Ayu102938/golang-todolist-tui@latest
`

※ $GOPATH/bin にパスが通っている必要があります。
Windowsの場合: %USERPROFILE%\go\bin を環境変数の Path に追加してください。

## 使い方

インストール後、以下のコマンドで起動します。

`ash
golang-todolist-tui
`

## 機能
- タスク管理（追加・削除・完了トグル・編集）
- 優先度設定（Low/Mid/High）
- タスク進捗の可視化（プログレスバー）
- 期限管理
- フィルタリング機能

## 操作方法
- j/k: 移動
- enter: 完了/未完了の切り替え
- a: 新規タスク追加
- e: タスク編集
- p: 優先度切り替え
- f: フィルタリング切り替え
- d: 削除
- q: 終了
