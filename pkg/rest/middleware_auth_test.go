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

var _ = Describe("Auth middleware", func() {
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
		)
		BeforeEach(func() {
			ctx := context.Background()

			id, err := types.StringToUUID(uuid.NewV4().String())
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
		})

		AfterEach(func() {
			userInfoServer.Close()
			apiServer.Close()
		})

		It("should refuse access to a protected endpoint without a token", func() {
			resp, err := client.NewRequest().Get("/users/me")
			Expect(err).To(BeNil())

			Expect(resp.IsError()).To(BeTrue())
			Expect(string(resp.Body())).To(ContainSubstring("The Authorization header is missing."))
		})

		It("should refuse access to a protected endpoint with a wrong type of a token", func() {
			resp, err := client.NewRequest().
				SetHeader("Authorization", "NotBearer :)").
				Get("/users/me")
			Expect(err).To(BeNil())

			Expect(resp.IsError()).To(BeTrue())
			Expect(string(resp.Body())).To(
				ContainSubstring("The Authorization header is malformed."),
			)
		})

		It("should refuse access to a protected endpoint with an invalid JWT token", func() {
			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer notjwtforsure").
				Get("/users/me")
			Expect(err).To(BeNil())

			Expect(resp.IsError()).To(BeTrue())
			Expect(string(resp.Body())).To(ContainSubstring("token contains an invalid number of segments"))
		})

		It("should refuse access to a protected endpoint with a JWT token using an invalid method",
			func() {
				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":                   "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
					"email_verified":        "true",
					"gender":                "Male",
					"name":                  "Leo Wilbur",
					"phone_number_verified": "true",
					"phone_number":          "+48123456789",
					"email":                 "leowilburdev@gmail.com",
					"username":              "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
				}).SignedString([]byte("somehmacsecret"))
				Expect(err).To(BeNil())

				resp, err := client.NewRequest().
					SetHeader("Authorization", "Bearer "+token).
					Get("/users/me")
				Expect(err).To(BeNil())

				Expect(resp.IsError()).To(BeTrue())
				Expect(string(resp.Body())).To(ContainSubstring("Unexpected signing method: HS256"))
			})

		It("should refuse access to a protected endpoint with a JWT token using an invalid method",
			func() {
				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
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
					"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
					"auth_time": 1541595219,
					"exp":       1541598819,
					"iat":       1541595219,
					"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
					"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
				}).SignedString([]byte("somehmacsecret"))
				Expect(err).To(BeNil())

				resp, err := client.NewRequest().
					SetHeader("Authorization", "Bearer "+token).
					Get("/users/me")
				Expect(err).To(BeNil())

				Expect(resp.IsError()).To(BeTrue())
				Expect(string(resp.Body())).To(ContainSubstring("Unexpected signing method: HS256"))
			})

		It("should refuse access to a protected endpoint with an expired token", func() {
			var (
				created = time.Now().Truncate(2 * time.Hour).Unix()
				expiry  = time.Now().Truncate(time.Hour).Unix()
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
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": created,
				"exp":       expiry,
				"iat":       created,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err := unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me")
			Expect(err).To(BeNil())

			Expect(resp.IsError()).To(BeTrue())
			Expect(string(resp.Body())).To(ContainSubstring("Token is expired"))
		})

		It("should provide user with access to the endpoint given a proper token", func() {
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
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": now,
				"exp":       in1h,
				"iat":       now,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err := unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me")
			Expect(err).To(BeNil())

			Expect(resp.IsError()).To(BeFalse())
			Expect(string(resp.Body())).To(ContainSubstring("Leo Wilbur"))
		})

		It("should provide user with access to an endpoint with more specific permissions", func() {
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
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": now,
				"exp":       in1h,
				"iat":       now,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err := unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me/lesions")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("[]"))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should deny the user access to endpoints that require more permissions", func() {
			var (
				now  = time.Now().Unix()
				in1h = time.Now().Add(time.Hour).Unix()
			)

			unfinishedToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"sub":            "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
				"cognito:groups": []string{},
				"iss":            "https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_gfSuuHw6e",
				"version":        2,
				"client_id":      "5d7sjcg4jmp5v8v3gdkpi8mvpi",
				"event_id":       "25bf4170-e28c-11e8-88be-1d6003dfe6e8",
				"token_use":      "access",
				// nolint
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/requests.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.read https://test.api.checkmoleapp.demo-redisys.com/lesions.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.lesions.write https://test.api.checkmoleapp.demo-redisys.com/lesions.read https://test.api.checkmoleapp.demo-redisys.com/questions.write phone https://test.api.checkmoleapp.demo-redisys.com/body-parts.write https://test.api.checkmoleapp.demo-redisys.com/owned.requests.create https://test.api.checkmoleapp.demo-redisys.com/requests.create email",
				"auth_time": now,
				"exp":       in1h,
				"iat":       now,
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err := unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Post("/questions")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Insufficient permissions"))
			Expect(resp.IsError()).To(BeTrue())
		})
	})
})
