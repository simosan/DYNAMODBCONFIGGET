package lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetItemfromKey パーティションキーのみで項目を取得するための構造体
type GetItemfromKey struct {
	Tbl  string // テーブル名
	Pky  string // パーティションキー
	Pkv  string // パーティションキーの値
	Pkvt string // パーティションキーの型（S:String,N:Numberだけに対応）
	Itm  string // 項目列名
}

// GetConfigItem メソッド（GetItem検索）
// 第1引数：AWSセッション
func (pk *GetItemfromKey) GetConfigItem(ses *session.Session) (string, error) {
	var params *dynamodb.GetItemInput
	// DynamoDBクライアント生成
	svc := dynamodb.New(ses)
	// パーティションキーの型が文字列型の場合
	if pk.Pkvt == "S" {
		params = &dynamodb.GetItemInput{
			TableName: aws.String(pk.Tbl),
			Key: map[string]*dynamodb.AttributeValue{
				// パーティションキー名
				pk.Pky: {
					// パーティションキーの値（検索したい値のキー)
					S: aws.String(pk.Pkv),
				},
			},
			AttributesToGet: []*string{
				// 欲しいデータの項目列
				aws.String(pk.Itm),
			},
			// 常に最新を取得するかどうか
			ConsistentRead: aws.Bool(true),

			//返ってくるデータの種類
			ReturnConsumedCapacity: aws.String("NONE"),
		}
		// パーティションキーの型が数値型の場
	} else if pk.Pkvt == "N" {
		params = &dynamodb.GetItemInput{
			TableName: aws.String(pk.Tbl),
			Key: map[string]*dynamodb.AttributeValue{
				// パーティションキー名
				pk.Pky: {
					// パーティションキーの値（検索したい値のキー)
					N: aws.String(pk.Pkv),
				},
			},
			AttributesToGet: []*string{
				// 欲しいデータの項目列
				aws.String(pk.Itm),
			},
			// 常に最新を取得するかどうか
			ConsistentRead: aws.Bool(true),

			//返ってくるデータの種類
			ReturnConsumedCapacity: aws.String("NONE"),
		}
	}

	resp, err := svc.GetItem(params)
	if err != nil {
		fmt.Println("DynamoDB GetItem Error! in (pk *GetItemfromKey) GetConfigItem")
	}

	return *resp.Item[pk.Itm].S, err
}
