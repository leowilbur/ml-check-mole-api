package rest_test

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	resty "gopkg.in/resty.v1"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/leowilbur/ml-check-mole-api/pkg/auth"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/rest"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

var _ = Describe("Body parts API", func() {
	var (
		dbConn = &pgx.ConnPool{}
		dbName string
	)

	BeforeEach(prepareDB(dbConn, &dbName))

	AfterEach(cleanupDB(dbConn, &dbName))

	Context("given an account and an API", func() {
		var (
			account        *models.Account
			jwk            *jose.JSONWebKeySet
			pair           *jose.JSONWebKey
			api            *rest.API
			userInfoServer *httptest.Server
			apiServer      *httptest.Server
			client         *resty.Client
			token          string
		)
		BeforeEach(func() {
			ctx := context.Background()

			id, err := types.StringToUUID("e7e42dd0-870e-443e-929d-b1ed7446ddd0")
			Expect(err).To(BeNil())

			account = &models.Account{
				ID:     id,
				Name:   "First Last",
				Email:  "hello@world.com",
				Phone:  "+48123456789",
				Gender: "Male",
			}
			Expect(models.UpsertAccount(ctx, dbConn, account)).To(BeNil())

			jwk, err = auth.CognitoJWK(jwkPublicSet)
			Expect(err).To(BeNil())

			pair = &jose.JSONWebKey{}
			Expect(json.Unmarshal([]byte(jwkKeyPair), pair)).To(BeNil())

			userInfoServer = httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{
						"sub": "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
						"email_verified": "true",
						"gender": "Male",
						"name": "Leo Wilbur",
						"phone_number_verified": "true",
						"phone_number": "+48123456789",
						"email": "leowilburdev@gmail.com",
						"username": "e7e42dd0-870e-443e-929d-b1ed7446ddd0"
					}`))
				}),
			)

			api, err = rest.New(dbConn, jwk, &rest.Config{
				AuthURL: userInfoServer.URL,
			})
			Expect(err).To(BeNil())

			apiServer = httptest.NewServer(api)

			client = resty.New().SetHostURL(apiServer.URL)

			var (
				now  = time.Now().Unix()
				in1h = time.Now().Add(time.Hour).Unix()
			)

			unfinishedToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"sub": "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
				"cognito:groups": []string{
					"Doctors",
					"Administrators",
				},
				"iss":       "https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_gfSuuHw6e",
				"version":   2,
				"client_id": "5d7sjcg4jmp5v8v3gdkpi8mvpi",
				"event_id":  "25bf4170-e28c-11e8-88be-1d6003dfe6e8",
				"token_use": "access",
				// nolint
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": now,
				"exp":       in1h,
				"iat":       now,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err = unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			userInfoServer.Close()
			apiServer.Close()
		})

		It("should properly list body parts in the database", func() {
			ctx := context.Background()

			Expect(models.CreateBodyPart(ctx, dbConn, &models.BodyPart{
				Name:      "Left knee",
				Displayed: true,
				Image:     "https://example.org/test.jpg",
				Order:     3,
			})).To(BeNil())

			Expect(models.CreateBodyPart(ctx, dbConn, &models.BodyPart{
				Name:      "Right knee",
				Displayed: true,
				Image:     "https://example.org/test.jpg",
				Order:     4,
			})).To(BeNil())

			Expect(models.CreateBodyPart(ctx, dbConn, &models.BodyPart{
				Name:      "Middle knee",
				Displayed: false,
				Image:     "https://example.org/test.jpg",
				Order:     1,
			})).To(BeNil())

			resp, err := client.NewRequest().
				Get("/body-parts")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Left knee"))
			Expect(string(resp.Body())).To(ContainSubstring("Right knee"))
			Expect(string(resp.Body())).To(ContainSubstring("Middle knee"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/body-parts?displayed=true")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Left knee"))
			Expect(string(resp.Body())).To(ContainSubstring("Right knee"))
			Expect(string(resp.Body())).ToNot(ContainSubstring("Middle knee"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				Get("/body-parts?displayed=false")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).ToNot(ContainSubstring("Left knee"))
			Expect(string(resp.Body())).ToNot(ContainSubstring("Right knee"))
			Expect(string(resp.Body())).To(ContainSubstring("Middle knee"))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should let an admin create, update and delete a new body part", func() {
			result := &models.BodyPart{}
			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"name":      "Left knee",
					"displayed": true,
					"image":     "https://someurl.org/lol.jpg",
					"order":     1,
				}).
				SetResult(&result).
				Post("/body-parts")
			Expect(err).To(BeNil())

			Expect(result.ID).ToNot(BeZero())
			Expect(string(resp.Body())).To(ContainSubstring("Left knee"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				Get("/body-parts")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Left knee"))
			Expect(resp.IsError()).To(BeFalse())

			result.Name = "Right knee"
			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(result).
				Put("/body-parts/" + result.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Right knee"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Delete("/body-parts/" + result.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Right knee"))
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
