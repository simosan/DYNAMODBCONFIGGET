package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SetAwsCredential AWSの認証を行う関数（iamユーザのクレデンシャルのみでのリソースアクセスは行えない）
// 第一引数：AWSプロファイル
func SetAwsCredential(profile string, region string, role string) (*session.Session, error) {

	var sess *session.Session
	var err error
	var arns string

	// AWSプロファイルが設定されていればprofileに値を設定,なければAWSプロファイルにprofileの値を設定
	if profile == "" {
		profile = os.Getenv("AWS_PROFILE")
	} else {
		os.Setenv("AWS_PROFILE", profile)
	}

	// IAMロールを指定している場合、~/.aws/config（role_arn)を読み込み、IAMロールのarnを取得
	if role != "" {
		arns, err = getIamRoleArn(role)
		if err != nil {
			return nil, err
		}
	}

	// IAMロールにスイッチして認証
	if profile != "" && role != "" {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
			// 第一引数はcredentialsの格納先を指定するが、空白にするとデフォルトパス（~/.aws）になる
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			fmt.Println("SetAwsCredential Error!", err)
			return nil, err
		}

		// iamロールにスイッチ
		screds := stscreds.NewCredentials(sess, arns)
		sconfig := aws.Config{Region: sess.Config.Region, Credentials: screds}
		sess = session.New(&sconfig)

		// EC2に適用されたIAMロールで認証
	} else if profile == "" {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			fmt.Println("SetAwsCredential(NoProfile) Error!", err)
			return nil, err
		}

	}

	return sess, err

}

// getIamRoleArn configに設定しているrole_arn（指定role）を読み込み、roleのarnを取得する
// 第１引数：　IAMロール名
func getIamRoleArn(role string) (string, error) {
	// role_arn文字列のパターン
	re1 := regexp.MustCompile("role_arn")
	// IAMロール名文字列のパターン
	re2 := regexp.MustCompile(role)

	// Linux/WindowsともにHOMEDIRのパスを取得できる仕様
	homedir, err := user.Current()
	f, err := os.Open(homedir.HomeDir + "/.aws/config")
	if err != nil {
		fmt.Println("getIamRoleArn Error! configが読み込めませんでした.")
		defer f.Close()
		return "", err
	}
	defer f.Close()

	// configを１行ずつ読みこむ
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if re1.MatchString(scanner.Text()) && re2.MatchString(scanner.Text()) {
			//空白でSplit(role_arn = arn:aws:iam~)
			arnstr := strings.Fields(scanner.Text())[2]

			return arnstr, err
		}
	}

	return "", err
}
