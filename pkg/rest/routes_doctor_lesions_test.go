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
	uuid "github.com/satori/go.uuid"
	resty "gopkg.in/resty.v1"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/leowilbur/ml-check-mole-api/pkg/auth"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/rest"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

var _ = Describe("Doctor lesions API", func() {
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
				Name:   "First Last",
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
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write phone https://test.api.checkmoleapp.demo-redisys.com/lesion.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
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

		It("should properly list lesions in the database", func() {
			ctx := context.Background()
			now := time.Now()

			Expect(models.CreateLesion(ctx, dbConn, &models.Lesion{
				AccountID:        account1.ID,
				Name:             "Mole on my knee #1",
				BodyPartID:       bodyPart1.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			})).To(BeNil())

			Expect(models.CreateLesion(ctx, dbConn, &models.Lesion{
				AccountID:        account1.ID,
				Name:             "Mole on my knee #2",
				BodyPartID:       bodyPart2.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			})).To(BeNil())

			Expect(models.CreateLesion(ctx, dbConn, &models.Lesion{
				AccountID:        account2.ID,
				Name:             "Mole on my knee #3",
				BodyPartID:       bodyPart2.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			})).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/lesions")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #2`))
			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #3`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_body_parts": "true",
					"offset":             "1",
					"limit":              "10",
					"body_part_id":       bodyPart2.ID.String(),
				}).
				Get("/lesions")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).ToNot(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).ToNot(ContainSubstring(`Mole on my knee #2`))
			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #3`))
			Expect(string(resp.Body())).To(ContainSubstring(`Right knee`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_body_parts": "true",
					"account_id":         account1.ID.String(),
				}).
				Get("/lesions")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #2`))
			Expect(string(resp.Body())).ToNot(ContainSubstring(`Mole on my knee #3`))
			Expect(string(resp.Body())).To(ContainSubstring(`Left knee`))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should be able to get a lesion from the database", func() {
			ctx := context.Background()
			now := time.Now()

			lesion := &models.Lesion{
				AccountID:        account1.ID,
				Name:             "Mole on my knee #1",
				BodyPartID:       bodyPart1.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			}

			Expect(models.CreateLesion(ctx, dbConn, lesion)).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/lesions/" + lesion.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).ToNot(ContainSubstring(`Left knee`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_body_part": "true",
				}).
				Get("/lesions/" + lesion.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).To(ContainSubstring(`Left knee`))
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
