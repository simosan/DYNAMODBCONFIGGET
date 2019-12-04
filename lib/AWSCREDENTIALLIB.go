package lib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SetAwsCredential AWSの認証を行う関数
// 第一引数：AWSプロファイル
func SetAwsCredential(profile string, region string) *session.Session {
	// profileが設定されている場合
	if profile != "" {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
			// 第一引数はcredentialsの格納先を指定するが、空白にするとデフォルトパス（~/.aws）になる
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			fmt.Println("SetAwsCredential Error!")
		}

		return sess
	}
	// profileが設定されていない場合≒AWS_PROFILE環境変数 Or IAMロール設定されたEC2の場合
	sess := session.Must(session.NewSession())
	return sess

}
