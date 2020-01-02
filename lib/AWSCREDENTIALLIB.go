package lib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SetAwsCredential AWSの認証を行う関数
// 第一引数：AWSプロファイル
// （注意）2019年時点の公式AWSライブラリ(Session)の仕様では、NewSeessionはerrorを返すとしているが実際はnilを返す。
//  よって、Session取得のエラーハンドリングは行えないことに注意
func SetAwsCredential(profile string, region string) (*session.Session, error) {
	// AWSプロファイルが設定されていればprofileに値を設定,なければAWSプロファイルにprofileの値を設定
	if profile == "" {
		profile = os.Getenv("AWS_PROFILE")
	} else {
		os.Setenv("AWS_PROFILE", profile)
	}
	// profileが設定されている場合
	if profile != "" {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
			// 第一引数はcredentialsの格納先を指定するが、空白にするとデフォルトパス（~/.aws）になる
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			fmt.Println("SetAwsCredential Error!", err)
		}

		return sess, err
	}
	// profileが設定されていない場合≒IAMロール設定されたEC2の場合
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("SetAwsCredential(NoProfile) Error!", err)
	}

	return sess, err

}
