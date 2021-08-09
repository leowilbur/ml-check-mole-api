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

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/auth"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/models"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/rest"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
)

var _ = Describe("Doctor reports API", func() {
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
			question1      *models.Question
			question2      *models.Question
			lesion1        *models.Lesion
			lesion2        *models.Lesion
			request1       *models.Request
			request2       *models.Request
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
				now  = time.Now()
				in1h = time.Now().Add(time.Hour)
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
				"auth_time": now.Unix(),
				"exp":       in1h.Unix(),
				"iat":       now.Unix(),
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

			lesion1 = &models.Lesion{
				AccountID:        account1.ID,
				Name:             "Mole on my knee #1",
				BodyPartID:       bodyPart1.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			}
			Expect(models.CreateLesion(ctx, dbConn, lesion1)).To(BeNil())

			lesion2 = &models.Lesion{
				AccountID:        account2.ID,
				Name:             "Mole on my knee #2",
				BodyPartID:       bodyPart2.ID,
				BodyPartLocation: "https://example.org/someurl.jpg",
				CreatedAt:        types.TimeToTimestamp(now),
				UpdatedAt:        types.TimeToTimestamp(now),
			}
			Expect(models.CreateLesion(ctx, dbConn, lesion2)).To(BeNil())

			request1 = &models.Request{
				AccountID: account1.ID,
				Status:    &models.StatusDraft,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			}
			Expect(models.CreateRequest(ctx, dbConn, request1)).To(BeNil())

			request2 = &models.Request{
				AccountID: account2.ID,
				Status:    &models.StatusSubmitted,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			}
			Expect(models.CreateRequest(ctx, dbConn, request2)).To(BeNil())

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

			report1 := &models.Report{
				RequestID: request1.ID,
				LesionID:  lesion1.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			}
			Expect(models.CreateReport(ctx, dbConn, report1)).To(BeNil())

			Expect(models.CreateReport(ctx, dbConn, &models.Report{
				RequestID: request1.ID,
				LesionID:  lesion2.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			})).To(BeNil())

			array.Set([]string{"world", "hello"})
			Expect(models.CreateReport(ctx, dbConn, &models.Report{
				RequestID: request2.ID,
				LesionID:  lesion2.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			})).To(BeNil())

			answer := &models.ReportAnswer{
				ReportID:   report1.ID,
				QuestionID: question1.ID,
				Answer: types.ExtendedJSONB{
					Status: pgtype.Present,
					Bytes:  []byte(`{"value":"answer1"}`),
				},
			}
			Expect(models.CreateReportAnswer(ctx, dbConn, answer)).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"offset": "1",
				}).
				Get("/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(request1.ID.String()))
			Expect(string(resp.Body())).To(ContainSubstring(`["world","hello"]`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_answers":   "true",
					"include_questions": "true",
					"limit":             "10",
					"request_id":        request1.ID.String(),
				}).
				Get("/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(request1.ID.String()))
			Expect(string(resp.Body())).ToNot(ContainSubstring(request2.ID.String()))
			Expect(string(resp.Body())).To(ContainSubstring("Do you smoke?"))
			Expect(string(resp.Body())).To(ContainSubstring("answer1"))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_lesions":    "true",
					"include_body_parts": "true",
					"limit":              "10",
					"lesion_id":          lesion1.ID.String(),
				}).
				Get("/reports")
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(lesion1.ID.String()))
			Expect(string(resp.Body())).ToNot(ContainSubstring(lesion2.ID.String()))
			Expect(string(resp.Body())).To(ContainSubstring("Left knee"))
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should be able to get a report from the database", func() {
			ctx := context.Background()
			now := time.Now()

			array := &types.ExtendedStringArray{}
			array.Set([]string{"hello", "world"})

			report := &models.Report{
				RequestID: request1.ID,
				LesionID:  lesion1.ID,
				Photos:    *array,
				CreatedAt: types.TimeToTimestamp(now),
				UpdatedAt: types.TimeToTimestamp(now),
			}
			Expect(models.CreateReport(ctx, dbConn, report)).To(BeNil())

			answer := &models.ReportAnswer{
				ReportID:   report.ID,
				QuestionID: question1.ID,
				Answer: types.ExtendedJSONB{
					Status: pgtype.Present,
					Bytes:  []byte(`{"value":"answer1"}`),
				},
			}
			Expect(models.CreateReportAnswer(ctx, dbConn, answer)).To(BeNil())

			resp, err := client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_lesion":    "true",
					"include_body_part": "true",
				}).
				Get("/reports/" + report.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`["hello","world"]`))
			Expect(string(resp.Body())).To(ContainSubstring(`Mole on my knee #1`))
			Expect(string(resp.Body())).To(ContainSubstring(`Left knee`))
			Expect(resp.IsError()).To(BeFalse())

			resp, err = client.NewRequest().
				SetHeader("Authorization", "Bearer "+token).
				SetQueryParams(map[string]string{
					"include_answers":   "true",
					"include_questions": "true",
				}).
				Get("/reports/" + report.ID.String())
			Expect(err).To(BeNil())

			Expect(string(resp.Body())).To(ContainSubstring(`answer1`))
			Expect(string(resp.Body())).To(ContainSubstring(`Do you smoke?`))
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
