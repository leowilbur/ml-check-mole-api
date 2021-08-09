package auth

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

/*
	const appClientID = `5d7sjcg4jmp5v8v3gdkpi8mvpi`
	const appClientSecret = `1mq4udn9co503gkvsm7p2a2ciav42di901qatobhmu398v5no7s1`
*/

var cognito *cognitoidentityprovider.CognitoIdentityProvider

func init() {
	awsSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	cognito = cognitoidentityprovider.New(
		awsSession, aws.NewConfig().WithRegion("ap-southeast-2"),
	)
}

// UserInfo contains all the details about an account from AWS Cognito's
// OIDC userinfo.
type UserInfo struct {
	Subject     string `json:"sub"`
	Gender      string `json:"gender"`
	BirthDate   string `json:"birthdate"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

// GetUserInfo acquires information about a user based on their token from
// AWS Cognito's OIDC endpoint.
func GetUserInfo(ctx context.Context, url string, username string) (*UserInfo, error) {
	out, err := cognito.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String("ap-southeast-2_gfSuuHw6e"),
		Username:   aws.String(username),
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get info")
	}

	userInfo := &UserInfo{}
	for _, attr := range out.UserAttributes {
		switch *attr.Name {
		case "email":
			userInfo.Email = *attr.Value
		case "phone_number":
			userInfo.PhoneNumber = *attr.Value
		case "name":
			userInfo.Name = *attr.Value
		case "gender":
			userInfo.Gender = *attr.Value
		case "sub":
			userInfo.Subject = *attr.Value
		case "birthdate":
			userInfo.BirthDate = *attr.Value
		}
	}

	return userInfo, nil

	/*
		req, err := http.NewRequest("GET", url+"/oauth2/userInfo", nil)
		if err != nil {
			return nil, errors.Wrap(err, "unable to prepare a request")
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req = req.WithContext(ctx)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get user info")
		}
		defer res.Body.Close()

		result := &UserInfo{}
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			return nil, errors.Wrap(err, "unable to decode the user info")
		}

		if result.Subject == "" {
			return nil, errors.New("Invalid token")
		}

		return result, nil
	*/
}
