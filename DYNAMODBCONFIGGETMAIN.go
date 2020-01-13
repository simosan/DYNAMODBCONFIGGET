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
		pktOpt   = flag.String("pt", "", "-pt パーティションキーの型（N Or S 必須）")
		profile  = flag.String("p", "", "-p AWSプロファイル（任意）")
		role     = flag.String("rl", "", "-rl IAMロール名（任意）") 
	)

	// コマンドライン引数の取得
	flag.Parse()

	// 引数チェック用の連想配列宣言（map）
	m := make(map[string]string)
	m["taskOpt"] = *taskOpt
	m["tblOpt"] = *tblOpt
	m["regonOpt"] = *regonOpt
	m["pkOpt"] = *pkOpt
	m["pkvOpt"] = *pkvOpt
	m["itemOpt"] = *itemOpt
	m["pktOpt"] = *pktOpt
	m["profile"] = *profile
	m["role"] = *role

	// 引数チェック
	if lib.CheckParam(m) != 0 {
		flag.Usage()
		os.Exit(255)
	}

	// errorハンドラー宣言
	var err error	

	// AWS認証セッション取得
	var sesclient *session.Session sesclient, err = lib.SetAwsCredential(*profile, *regonOpt, *role)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	// 検索条件（GetItem Or Query）を判定し、構造体にパラメータをセット
	// DynamoDBにアクセスし検索結果を得る
	var rtn string
	var rtnbyte []byte
	// GetItem
	if *taskOpt == "get" {
		var pki lib.GetItemfromKey
		pki.Tbl = *tblOpt
		pki.Pky = *pkOpt
		pki.Pkv = *pkvOpt
		pki.Pkvt = *pktOpt
		pki.Itm = *itemOpt
		// DynamoDBクライアントアクセス
		rtn, err = pki.GetConfigItem(sesclient)
		if err != nil {
			os.Exit(255)
		} else {
			fmt.Println(rtn)
		}
		// Query
	} else if *taskOpt == "query" {
		var qfk lib.QueryfromKey
		qfk.Tbl = *tblOpt
		qfk.Pky = *pkOpt
		qfk.Pkv = *pkvOpt
		qfk.Pkvt = *pktOpt
		// DynamoDBクライアントアクセス
		rtnbyte, err = qfk.GetConfigItem(sesclient)
		if err != nil {
			os.Exit(255)
		} else {
			fmt.Println(string(rtnbyte))
		}
	}

	return

}
