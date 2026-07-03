package main

type messages struct {
	Title               string
	TaskPlaceholder     string
	CategoryPlaceholder string
	SearchPlaceholder   string
	DescPlaceholder     string
	DatePlaceholder     string
	Confirm             string
	Cancel              string
	Return              string
	DeleteCategory      string
	DeleteConfirm       string
	CannotDeleteHome    string
	Help                string
	NoTasks             string
	NoVisibleTasks      string
	SortNone            string
	SortPriority        string
	SortDueDate         string
	SortName            string
}

var japaneseMessages = messages{
	Title:               "TODO リスト",
	TaskPlaceholder:     "タスク名...",
	CategoryPlaceholder: "カテゴリ名...",
	SearchPlaceholder:   "検索...",
	DescPlaceholder:     "詳細...",
	DatePlaceholder:     "期限 YYYY-MM-DD (enter=明日)...",
	Confirm:             "確定",
	Cancel:              "キャンセル",
	Return:              "戻る",
	DeleteCategory:      "カテゴリ削除",
	DeleteConfirm:       "カテゴリとそのタスクを削除しますか？",
	CannotDeleteHome:    "カテゴリは削除できません",
	Help:                " h/l:タブ • n:カテゴリ追加 • x:カテゴリ削除 • a:追加 • j/k:移動 • e:編集 • i:詳細 • p:優先度 • s:ソート • f:フィルター • /:検索 • d:削除 • u:元に戻す • q:終了",
	NoTasks:             "タスクがありません",
	NoVisibleTasks:      "表示できるタスクがありません",
	SortNone:            "なし",
	SortPriority:        "優先度",
	SortDueDate:         "期限",
	SortName:            "名前",
}

var msg = japaneseMessages
