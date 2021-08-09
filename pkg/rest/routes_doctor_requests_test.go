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
	"github.com/jackc/pgx/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	resty "gopkg.in/resty.v1"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/leowilbur/ml-check-mole-api/pkg/auth"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/rest"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

var _ = Describe("Doctor requests API", func() {
	var (
		dbConn = &pgx.ConnPool{}
		dbName string
	)

	BeforeEach(prepareDB(dbConn, &dbName))

	AfterEach(cleanupDB(dbConn, &dbName))

	Context("given an account and an API", func() {
		var (
			account1       *models.Account
			account2       *models.Account
			jwk            *jose.JSONWebKeySet
			pair           *jose.JSONWebKey
			api            *rest.API
			userInfoServer *httptest.Server
			apiServer      *httptest.Server
			client         *resty.Client
			token          string
			bodyPart1      *models.BodyPart
			bodyPart2      *models.BodyPart
		)
		BeforeEach(func() {
			ctx := context.Background()

			id, err := types.StringToUUID("e7e42dd0-870e-443e-929d-b1ed7446ddd0")
			Expect(err).To(BeNil())

			account1 = &models.Account{
				ID:     id,
				Name:   "First Last",
				Email:  "hello@world.com",
				Phone:  "+48123456789",
				Gender: "Male",
			}
			Expect(models.UpsertAccount(ctx, dbConn, account1)).To(BeNil())

			id2, err := types.StringToUUID(uuid.NewV4().String())
			Expect(err).To(BeNil())

			account2 = &models.Account{
				ID:     id2,
				Name:   "Some Other",
				Email:  "hello@world.com",
				Phone:  "+48123456789",
				Gender: "Male",
			}
			Expect(models.UpsertAccount(ctx, dbConn, account2)).To(BeNil())

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
				},
				"iss":       "https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_gfSuuHw6e",
				"version":   2,
				"client_id": "5d7sjcg4jmp5v8v3gdkpi8mvpi",
				"event_id":  "25bf4170-e28c-11e8-88be-1d6003dfe6e8",
				"token_use": "access",
				// nolint
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.requests.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write phone https://test.api.checkmoleapp.demo-redisys.com/lesion.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": now,
				"exp":       in1h,
				"iat":       now,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err = unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			bodyPart1 = &models.BodyPart{
				Name:      "Left knee",
				Displayed: true,
				Image:     "https://example.org/test.jpg",
				Order:     3,
			}
			Expect(models.CreateBodyPart(ctx, dbConn, bodyPart1)).To(BeNil())

			bodyPart2 = &models.BodyPart{
				Name:      "Right knee",
				Displayed: true,
				Image:     "https://example.org/test.jpg",
				Order:     4,
			}
			Expect(models.CreateBodyPart(ctx, dbConn, bodyPart2)).To(BeNil())
		})

		AfterEach(func() {
			userInfoServer.Close()
			apiServer.Close()
		})

		It("should properly list requests in the database", func() {
			ctx := context.Background()

			Expect(models.CreateRequest(ctx, dbConn, &models.Request{
				AccountID: account1.ID,
				Status:    &models.StatusOpen,
			})).To(BeNil())

			Expect(models.CreateRequest(ctx, dbConn, &models.Request{
				AccountID: account1.ID,
				Status:    &models.StatusSubmitted,
			})).To(BeNil())

			Expect(models.CreateRequest(ctx, dbConn, &models.Request{
				AccountID: account2.ID,
				Status:    &models.StatusAnswered,
			})).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/requests")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Open`))
			Expect(string(resp.Body())).To(ContainSubstring(`Submitted`))
			Expect(string(resp.Body())).To(ContainSubstring(`Answered`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_accounts": "true",
					"offset":           "0",
					"limit":            "10",
					"account_id":       account1.ID.String(),
				}).
				Get("/requests")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Open`))
			Expect(string(resp.Body())).To(ContainSubstring(`Submitted`))
			Expect(string(resp.Body())).To(ContainSubstring(`First Last`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_accounts": "true",
					"status":           "Answered",
				}).
				Get("/requests")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).ToNot(ContainSubstring(`Open`))
			Expect(string(resp.Body())).ToNot(ContainSubstring(`Submitted`))
			Expect(string(resp.Body())).To(ContainSubstring(`Answered`))
			Expect(string(resp.Body())).To(ContainSubstring(`Some Other`))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should be able to get a request from the database and update it", func() {
			ctx := context.Background()

			request := &models.Request{
				AccountID: account1.ID,
				Status:    &models.StatusOpen,
				AnswerText: types.ExtendedText{
					Status: pgtype.Null,
				},
			}

			Expect(models.CreateRequest(ctx, dbConn, request)).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/requests/" + request.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Open`))
			Expect(string(resp.Body())).ToNot(ContainSubstring(`First Last`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_account": "true",
				}).
				Get("/requests/" + request.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Open`))
			Expect(string(resp.Body())).To(ContainSubstring(`First Last`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"status":      "Answered",
					"answer_text": "cool beans",
					"answered_by": "cool doctor",
				}).
				Put("/requests/" + request.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Answered`))
			Expect(string(resp.Body())).To(ContainSubstring(`cool beans`))
			Expect(string(resp.Body())).To(ContainSubstring(`cool doctor`))
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
