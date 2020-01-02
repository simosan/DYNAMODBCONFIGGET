package lib

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// QueryfromKey 検索条件を項目とする構造体
type QueryfromKey struct {
	Tbl  string // テーブル名
	Pky  string // パーティションキー
	Pkv  string // パーティションキーの値
	Pkvt string // パーティションキーの型（S:String,N:Numberだけに対応）
}

// GetConfigItem メソッド （Query検索,列は全件出力）
// 第1引数：AWSセッション
func (qk *QueryfromKey) GetConfigItem(ses *session.Session) ([]byte, error) {
	var input *dynamodb.QueryInput
	// DynamoDBクライアント生成
	svc := dynamodb.New(ses)
	// パーティションキーの型が文字列型の場合
	if qk.Pkvt == "S" {
		input = &dynamodb.QueryInput{
			TableName: aws.String(qk.Tbl),
			ExpressionAttributeNames: map[string]*string{
				"#ID": aws.String(qk.Pky),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				// 検索対象パーティションキーの値(文字列型の場合)
				":id": {
					S: aws.String(qk.Pkv),
				},
			},
			//　検索条件（条件とキーが完全一致）
			KeyConditionExpression: aws.String("#ID = :id"),
		}
		// パーティションキーの型が数値型の場合
	} else if qk.Pkvt == "N" {
		input = &dynamodb.QueryInput{
			TableName: aws.String(qk.Tbl),
			ExpressionAttributeNames: map[string]*string{
				// 予約語とぶつからないように＃でプレースホルダ
				// パーティションキー
				"#ID": aws.String(qk.Pky),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				// 検索対象パーティションキーの値(数値型の場合)
				":id": {
					N: aws.String(qk.Pkv),
				},
			},
			//　検索条件（条件とキーが完全一致）
			KeyConditionExpression: aws.String("#ID = :id"),
		}
	}

	// クエリ実行
	result, err := svc.Query(input)
	if err != nil {
		fmt.Println("Query Error! in (qk *GetItemfromKey) GetConfigItem ", err)
		return nil, err
	}

	// Dynamodb型のデータ形式をJson形式に変換
	j, _ := json.Marshal(result)

	return j, err
}
