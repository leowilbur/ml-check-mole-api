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
	resty "gopkg.in/resty.v1"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/leowilbur/ml-check-mole-api/pkg/auth"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/rest"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

var _ = Describe("User reports API", func() {
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
			bodyPart       *models.BodyPart
			lesion         *models.Lesion
			request        *models.Request
			question1      *models.Question
			question2      *models.Question
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
				now  = time.Now()
				in1h = time.Now().Add(time.Hour)
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
				"scope":     "https://test.api.checkmoleapp.demo-redisys.com/reports.respond aws.cognito.signin.user.admin https://test.api.checkmoleapp.demo-redisys.com/owned.reports.read https://test.api.checkmoleapp.demo-redisys.com/reports.write openid profile https://test.api.checkmoleapp.demo-redisys.com/owned.reports.write https://test.api.checkmoleapp.demo-redisys.com/reports.read https://test.api.checkmoleapp.demo-redisys.com/reports.write phone https://test.api.checkmoleapp.demo-redisys.com/report.write https://test.api.checkmoleapp.demo-redisys.com/owned.reports.create https://test.api.checkmoleapp.demo-redisys.com/reports.create email",
				"auth_time": now.Unix(),
				"exp":       in1h.Unix(),
				"iat":       now.Unix(),
				"jti":       "4a5ac6c1-f001-4d90-8c23-65ff781b5a95",
				"username":  "e7e42dd0-870e-443e-929d-b1ed7446ddd0",
			})

			unfinishedToken.Header["kid"] = "KWZ5ZSZZ"

			token, err = unfinishedToken.SignedString(pair.Key.(*rsa.PrivateKey))
			Expect(err).To(BeNil())

			bodyPart = &models.BodyPart{
				Name:      "Left knee",
				Displayed: true,
				Image:     "https://example.org/test.jpg",
				Order:     3,
			}
			Expect(models.CreateBodyPart(ctx, dbConn, bodyPart)).To(BeNil())

			lesion = &models.Lesion{
				AccountID:        account.ID,
				Name:             "Mole on my knee #1",
				BodyPartID:       bodyPart.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
			}
			Expect(models.CreateLesion(ctx, dbConn, lesion)).To(BeNil())

			request = &models.Request{
				AccountID: account.ID,
				Status:    &models.StatusDraft,
				AnswerText: types.ExtendedText{
					Status: pgtype.Null,
				},
			}
			Expect(models.CreateRequest(ctx, dbConn, request)).To(BeNil())

			question1 = &models.Question{
				Name: "Do you smoke?",
				Type: "select",
				Answers: types.ExtendedJSONB{
					Status: pgtype.Present,
					Bytes:  []byte(`{"hello": "world1"}`),
				},
				Displayed: true,
				Order:     3,
			}
			Expect(models.CreateQuestion(ctx, dbConn, question1)).To(BeNil())

			question2 = &models.Question{
				Name: "Do you really smoke?",
				Type: "select",
				Answers: types.ExtendedJSONB{
					Status: pgtype.Present,
					Bytes:  []byte(`{"hello": "world1"}`),
				},
				Displayed: true,
				Order:     4,
			}
			Expect(models.CreateQuestion(ctx, dbConn, question2)).To(BeNil())
		})

		AfterEach(func() {
			userInfoServer.Close()
			apiServer.Close()
		})

		It("should properly list reports in the database", func() {
			ctx := context.Background()
			now := time.Now()

			array := &types.ExtendedStringArray{}
			array.Set([]string{"hello", "world"})

			Expect(models.CreateReport(ctx, dbConn, &models.Report{
				RequestID: request.ID,
				LesionID:  lesion.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			})).To(BeNil())

			array.Set([]string{"hello2", "world2"})
			Expect(models.CreateReport(ctx, dbConn, &models.Report{
				RequestID: request.ID,
				LesionID:  lesion.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			})).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me/lesions/" + lesion.ID.String() + "/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`"hello","world"`))
			Expect(string(resp.Body())).To(ContainSubstring(`"hello2","world2"`))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should prevent the user from creating a report with invalid request id", func() {
			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"request_id": account.ID.String(),
					"photos": []string{
						"photo1",
						"photo2",
					},
					"answers": []map[string]interface{}{
						{
							"question_id": question1.ID.String(),
							"answer":      "some_answer1",
						},
						{
							"question_id": question2.ID.String(),
							"answer":      "some_answer2",
						},
					},
				}).
				Post("/users/me/lesions/" + lesion.ID.String() + "/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("Request not found"))
			Expect(resp.IsError()).To(BeTrue())
		})

		It("should let the user from creating a report with a valid request ID", func() {
			result := &models.Report{}
			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"request_id": request.ID.String(),
					"photos": []string{
						"photo1",
						"photo2",
					},
					"answers": []map[string]interface{}{
						{
							"question_id": question1.ID.String(),
							"answer":      "some_answer1",
						},
						{
							"question_id": question2.ID.String(),
							"answer":      "some_answer2",
						},
					},
				}).
				SetResult(&result).
				Post("/users/me/lesions/" + lesion.ID.String() + "/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(request.ID.String()))
			Expect(string(resp.Body())).To(ContainSubstring("photo1"))
			Expect(string(resp.Body())).To(ContainSubstring("photo2"))
			Expect(result.ID).ToNot(BeZero())
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should let the user create, update and delete a new report", func() {
			result := &models.Report{}
			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"photos": []string{
						"photo1",
						"photo2",
					},
					"answers": []map[string]interface{}{
						{
							"question_id": question1.ID.String(),
							"answer":      "some_answer1",
						},
					},
				}).
				SetResult(&result).
				Post("/users/me/lesions/" + lesion.ID.String() + "/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("photo1"))
			Expect(string(resp.Body())).To(ContainSubstring("photo2"))
			Expect(result.ID).ToNot(BeZero())
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me/lesions/" + lesion.ID.String() + "/reports?include_answers=true")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("photo1"))
			Expect(string(resp.Body())).To(ContainSubstring("some_answer1"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"request_id": request.ID.String(),
					"photos": []string{
						"photo3",
						"photo4",
					},
					"answers": []map[string]interface{}{
						{
							"question_id": question1.ID.String(),
							"answer":      "some_answer3",
						},
						{
							"question_id": question2.ID.String(),
							"answer":      "actually weed lol",
						},
					},
				}).
				Put("/users/me/lesions/" + lesion.ID.String() + "/reports/" + result.ID.String())
			Expect(err).To(BeNil())
			Expect(string(resp.Body())).To(ContainSubstring("photo3"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me/lesions/" + lesion.ID.String() + "/reports?include_answers=true")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).ToNot(ContainSubstring("some_answer1"))
			Expect(string(resp.Body())).To(ContainSubstring("some_answer3"))
			Expect(string(resp.Body())).To(ContainSubstring("actually weed lol"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetBody(map[string]interface{}{
					"request_id": request.ID.String(),
					"photos": []string{
						"photo3",
						"photo4",
					},
					"answers": []map[string]interface{}{
						{
							"question_id": question2.ID.String(),
							"answer":      "actually weed lol",
						},
					},
				}).
				Put("/users/me/lesions/" + lesion.ID.String() + "/reports/" + result.ID.String())
			Expect(err).To(BeNil())
			Expect(string(resp.Body())).ToNot(ContainSubstring("some_answer3"))
			Expect(string(resp.Body())).To(ContainSubstring("photo3"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Get("/users/me/lesions/" + lesion.ID.String() + "/reports?include_answers=true")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("actually weed lol"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				Delete("/users/me/lesions/" + lesion.ID.String() + "/reports/" + result.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring("photo3"))
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
