package main

import (
	"flag"
	"fmt"
	"os"

	"./lib"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Main関数
func main() {

	// 引数フラグ、ヘルプメッセージ定義
	var (
		taskOpt  = flag.String("task", "", "-task [get|query]")
		tblOpt   = flag.String("t", "", "-t テーブル名（必須）")
		regonOpt = flag.String("r", "", "-r リージョン名（必須）")
		pkOpt    = flag.String("k", "", "-k パーティションキー（必須）")
		pkvOpt   = flag.String("v", "", "-v パーティションキー値（必須）")
		itemOpt  = flag.String("i", "", "-i 表示対象項目名（task=get）")
		pktOpt   = flag.String("pt", "", "-pt パーティションキーの型（N Or S（task＝queryの場合必須）")
		profile  = flag.String("p", "", "-p AWSプロファイル（任意）")
	)
	// コマンドライン引数の取得
	flag.Parse()

	// 引数タスクチェック
	if *taskOpt == "" || (*taskOpt != "get" && *taskOpt != "query") {
		fmt.Println("コマンドライン引数に誤りがあります.タスク指定の誤りです.")
		flag.Usage()
		os.Exit(255)
	}

	// 各引数のチェック
	if *tblOpt == "" || *pkOpt == "" || *pkvOpt == "" || *regonOpt == "" {
		fmt.Println("コマンドライン引数に誤りがあります.必須項目に不足があります.")
		flag.Usage()
		os.Exit(255)
	}

	// taskの場合のチェック
	if *taskOpt == "get" && *itemOpt == "" {
		fmt.Println("コマンドライン引数に誤りがあります.getオプションに誤りがあります.")
		flag.Usage()
		os.Exit(255)
	}

	// queryの場合のチェック
	if *taskOpt == "query" && (*pktOpt != "N" && *pktOpt != "S") {
		fmt.Println("コマンドライン引数に誤りがあります.queryオプションに誤りがあります.")
		flag.Usage()
		os.Exit(255)
	}

	// AWSプロファイルの設定
	if *profile == "" {
		*profile = os.Getenv("AWS_PROFILE")
	} else {
		os.Setenv("AWS_PROFILE", *profile)
	}

	// AWS認証
	var sesclient *session.Session
	sesclient = lib.SetAwsCredential(*profile, *regonOpt)

	// 検索条件（GetItem Or Query）を判定し、構造体にパラメータをセット
	// DynamoDBにアクセスし検索結果を得る
	var rtn string
	var rtnbyte []byte
	var err error
	// GetItem
	if *taskOpt == "get" {
		var pki lib.GetItemfromKey
		pki.Tbl = *tblOpt
		pki.Pky = *pkOpt
		pki.Pkv = *pkvOpt
		pki.Itm = *itemOpt
		// DynamoDBクライアントアクセス
		rtn, err = pki.GetConfigItem(sesclient)
		fmt.Println(rtn)
		// Query
	} else if *taskOpt == "query" {
		var qfk lib.QueryfromKey
		qfk.Tbl = *tblOpt
		qfk.Pky = *pkOpt
		qfk.Pkv = *pkvOpt
		qfk.Pkvt = *pktOpt
		// DynamoDBクライアントアクセス
		rtnbyte, err = qfk.GetConfigItem(sesclient)
		fmt.Println(string(rtnbyte))
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	return

}
