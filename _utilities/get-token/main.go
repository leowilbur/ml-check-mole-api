package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var (
	client = flag.String("client", "5d7sjcg4jmp5v8v3gdkpi8mvpi", "client id")
	secret = flag.String("secret", "1mq4udn9co503gkvsm7p2a2ciav42di901qatobhmu398v5no7s1", "client secret")
)

func computeSecretHash(username string) string {
	hash := hmac.New(sha256.New, []byte(*secret))
	hash.Write([]byte(username))
	hash.Write([]byte(*client))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func main() {
	flag.Parse()

	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-southeast-2"),
	})

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username [leowilburdev@gmail.com]: ")
	rawUsername, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	trimmedUsername := strings.TrimSpace(rawUsername)
	if trimmedUsername == "" {
		trimmedUsername = "leowilburdev@gmail.com"
	}

	fmt.Print("Password [X4s4aFxGLVNVThAi]: ")
	rawPassword, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	trimmedPassword := strings.TrimSpace(rawPassword)
	if trimmedPassword == "" {
		trimmedPassword = "X4s4aFxGLVNVThAi"
	}

	authResp, err := svc.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(trimmedUsername),
			"PASSWORD":    aws.String(trimmedPassword),
			"SECRET_HASH": aws.String(computeSecretHash(trimmedUsername)),
		},
		ClientId: aws.String(*client),
	})
	if err != nil {
		panic(err)
	}

	if authResp.AuthenticationResult != nil {
		fmt.Printf("Logged in, the token will expire in %d:\n", *authResp.AuthenticationResult.ExpiresIn)
		fmt.Println(*authResp.AuthenticationResult.AccessToken)
		return
	}

	if authResp.ChallengeName == nil {
		fmt.Printf("Authentication failed. Expected a challenge, got %+v\n", authResp)
	}

	if *authResp.ChallengeName != "SMS_MFA" {
		fmt.Printf("Authentication failed. SMS MFA challenge expected, got %s.\n", *authResp.ChallengeName)
	}

	fmt.Print("Please enter the SMS code: ")
	rawCode, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	trimmedCode := strings.TrimSpace(rawCode)

	// SMS_MFA_CODE
	challResp, err := svc.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: aws.String("SMS_MFA"),
		ChallengeResponses: map[string]*string{
			"SMS_MFA_CODE": aws.String(trimmedCode),
			"USERNAME":     aws.String(trimmedUsername),
			"SECRET_HASH":  aws.String(computeSecretHash("leowilburdev@gmail.com")),
		},
		ClientId: aws.String(*client),
		Session:  authResp.Session,
	})
	if err != nil {
		panic(err)
	}

	if challResp.AuthenticationResult != nil {
		fmt.Printf("Logged in, the token will expire in %d:\n", *challResp.AuthenticationResult.ExpiresIn)
		fmt.Println(*challResp.AuthenticationResult.AccessToken)
		return
	}

	log.Printf("%+v", challResp)
}
