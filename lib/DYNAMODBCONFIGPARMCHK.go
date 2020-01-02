package lib

import (
	"fmt"
)

// CheckParam コマンドライン引数のバリデーションを行う
// 第一引数: コマンドライン引数の連想配列(map)
func CheckParam(mp map[string]string) int {

	var rtnval int = 0

	// 引数タスクチェック
	if mp["taskOpt"] == "" || (mp["taskOpt"] != "get" && mp["taskOpt"] != "query") {
		fmt.Println("コマンドライン引数に誤りがあります.タスク指定の誤りです.")
		rtnval = 255
	}

	// 各引数のチェック
	if mp["tblOpt"] == "" || mp["pkOpt"] == "" || mp["pkvOpt"] == "" || mp["regonOpt"] == "" {
		fmt.Println("コマンドライン引数に誤りがあります.必須項目に不足があります.")
		rtnval = 255
	}

	// taskの場合のチェック
	if mp["taskOpt"] == "get" && mp["itemOpt"] == "" {
		fmt.Println("コマンドライン引数に誤りがあります.getオプションに誤りがあります.")
		rtnval = 255
	}

	// 型のチェック
	if mp["pktOpt"] != "N" && mp["pktOpt"] != "S" {
		fmt.Println("コマンドライン引数に誤りがあります.パーティションキーの型指定に誤りがあります.")
		rtnval = 255
	}

	// DynamoDBの名前ルールに準拠しているか（テーブル、インデックスが255文字以内か）
	if len(mp["tblOpt"]) > 255 || len(mp["pkOpt"]) > 255 {
		fmt.Println("コマンドライン引数に誤りがあります.テーブル名、インデックス名が長すぎます")
		rtnval = 255
	}

	return rtnval
}
