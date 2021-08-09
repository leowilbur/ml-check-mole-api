package models_test

import (
	"context"

	"github.com/dchest/uniuri"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

var _ = Describe("Database model", func() {
	var (
		dbConn *pgx.Conn
		dbName string
	)

	BeforeEach(func() {
		setupConfig, err := pgx.ParseURI(baseDBURI + "/molepatrol?sslmode=disable")
		Expect(err).To(BeNil())
		setupConn, err := pgx.Connect(setupConfig)
		Expect(err).To(BeNil())

		dbName = uniuri.NewLenChars(16, []byte("abcdefghijklmnopqrstuvwxyz"))
		_, err = setupConn.Exec("CREATE DATABASE " + dbName)
		Expect(err).To(BeNil())
		Expect(setupConn.Close()).To(BeNil())

		testConfig, err := pgx.ParseURI(baseDBURI + "/" + dbName + "?sslmode=disable")
		Expect(err).To(BeNil())
		dbConn, err = pgx.Connect(testConfig)
		Expect(err).To(BeNil())

		m, err := migrate.New(
			"file://../../migrations",
			baseDBURI+"/"+dbName+"?sslmode=disable",
		)
		Expect(err).To(BeNil())
		Expect(m.Up()).To(BeNil())
		srcErr, dbErr := m.Close()
		Expect(srcErr).To(BeNil())
		Expect(dbErr).To(BeNil())
	})

	AfterEach(func() {
		dbConn.Close()

		setupConfig, err := pgx.ParseURI(baseDBURI + "/molepatrol?sslmode=disable")
		Expect(err).To(BeNil())
		setupConn, err := pgx.Connect(setupConfig)
		Expect(err).To(BeNil())

		_, err = setupConn.Exec("DROP DATABASE " + dbName)
		Expect(err).To(BeNil())
		Expect(setupConn.Close()).To(BeNil())
	})

	It("Account should upsert and select properly", func() {
		ctx := context.Background()

		id, err := types.StringToUUID(uuid.NewV4().String())
		Expect(err).To(BeNil())

		input := &models.Account{
			ID:     id,
			Name:   "First Last",
			Email:  "hello@world.com",
			Phone:  "+48123456789",
			Gender: "Male",
		}
		Expect(models.UpsertAccount(ctx, dbConn, input)).To(BeNil())

		account1, err := models.GetAccount(ctx, dbConn, id)
		Expect(err).To(BeNil())

		Expect(account1.ID).To(Equal(input.ID))
		Expect(account1.Name).To(Equal(input.Name))
		Expect(account1.Email).To(Equal(input.Email))
		Expect(account1.Phone).To(Equal(input.Phone))
		Expect(account1.Gender).To(Equal(input.Gender))

		// Do the 2nd upsert
		Expect(models.UpsertAccount(ctx, dbConn, input)).To(BeNil())
		account2, err := models.GetAccount(ctx, dbConn, id)
		Expect(err).To(BeNil())

		Expect(account2.CreatedAt.Time).To(Equal(account1.CreatedAt.Time))
		Expect(account2.UpdatedAt.Time.After(account1.UpdatedAt.Time)).To(BeTrue())
	})

	It("BodyPart should correctly perform the CRUD flow", func() {
		ctx := context.Background()

		input := &models.BodyPart{
			Name:      "Left knee",
			Displayed: true,
			Image:     "https://example.org/test.jpg",
			Order:     3,
		}
		Expect(models.CreateBodyPart(ctx, dbConn, input)).To(BeNil())

		bp, err := models.GetBodyPart(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		Expect(bp.ID).To(Equal(input.ID))
		Expect(bp.Name).To(Equal(input.Name))
		Expect(bp.Displayed).To(Equal(input.Displayed))
		Expect(bp.Image).To(Equal(input.Image))
		Expect(bp.Order).To(Equal(input.Order))

		bp.Order = 4
		Expect(models.UpdateBodyPart(ctx, dbConn, bp)).To(BeNil())
		Expect(bp.Order).To(Equal(4))

		_, err = models.DeleteBodyPart(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())
	})

	It("Lesion should correctly perform the CRUD flow", func() {
		ctx := context.Background()

		// First insert the account
		accountID, err := types.StringToUUID(uuid.NewV4().String())
		Expect(err).To(BeNil())
		Expect(models.UpsertAccount(ctx, dbConn, &models.Account{
			ID:     accountID,
			Name:   "First Last",
			Email:  "hello@world.com",
			Phone:  "+48123456789",
			Gender: "Male",
		})).To(BeNil())

		// Then the body part
		bodyPart := &models.BodyPart{
			Name:      "Left knee",
			Displayed: true,
			Image:     "https://example.org/test.jpg",
			Order:     3,
		}
		Expect(models.CreateBodyPart(ctx, dbConn, bodyPart)).To(BeNil())

		// Then insert the lesion
		input := &models.Lesion{
			AccountID:        accountID,
			Name:             "Mole on my knee",
			BodyPartID:       bodyPart.ID,
			BodyPartLocation: "https://example.org/test.jpg",
		}
		Expect(models.CreateLesion(ctx, dbConn, input)).To(BeNil())

		lesion, err := models.GetLesion(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		Expect(lesion.ID).To(Equal(input.ID))
		Expect(lesion.AccountID).To(Equal(input.AccountID))
		Expect(lesion.Name).To(Equal(input.Name))
		Expect(lesion.BodyPartID).To(Equal(input.BodyPartID))
		Expect(lesion.BodyPartLocation).To(Equal(input.BodyPartLocation))

		lesion.Name = "Mole on left knee"
		Expect(models.UpdateLesion(ctx, dbConn, lesion)).To(BeNil())
		Expect(lesion.Name).To(Equal("Mole on left knee"))

		_, err = models.DeleteLesion(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())
	})

	It("Question should correctly perform the CRUD flow", func() {
		ctx := context.Background()

		input := &models.Question{
			Name: "Do you smoke?",
			Type: "select",
			Answers: types.ExtendedJSONB{
				Status: pgtype.Present,
				Bytes:  []byte(`{"hello":"world"}`),
			},
			Displayed: true,
			Order:     3,
		}
		Expect(models.CreateQuestion(ctx, dbConn, input)).To(BeNil())

		bp, err := models.GetQuestion(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		Expect(bp.ID).To(Equal(input.ID))
		Expect(bp.Name).To(Equal(input.Name))
		Expect(bp.Type).To(Equal(input.Type))
		Expect(bp.Answers).To(Equal(input.Answers))
		Expect(bp.Displayed).To(Equal(input.Displayed))
		Expect(bp.Order).To(Equal(input.Order))

		bp.Order = 4
		Expect(models.UpdateQuestion(ctx, dbConn, bp)).To(BeNil())
		Expect(bp.Order).To(Equal(int64(4)))

		_, err = models.DeleteQuestion(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())
	})

	It("Request should correctly perform the CRUD flow", func() {
		ctx := context.Background()

		// First insert the account
		accountID, err := types.StringToUUID(uuid.NewV4().String())
		Expect(err).To(BeNil())
		Expect(models.UpsertAccount(ctx, dbConn, &models.Account{
			ID:     accountID,
			Name:   "First Last",
			Email:  "hello@world.com",
			Phone:  "+48123456789",
			Gender: "Male",
		})).To(BeNil())

		// Then the body part
		bodyPart := &models.BodyPart{
			Name:      "Left knee",
			Displayed: true,
			Image:     "https://example.org/test.jpg",
			Order:     3,
		}
		Expect(models.CreateBodyPart(ctx, dbConn, bodyPart)).To(BeNil())

		// Then insert the lesion
		lesion := &models.Lesion{
			AccountID:        accountID,
			Name:             "Mole on my knee",
			BodyPartID:       bodyPart.ID,
			BodyPartLocation: "https://example.org/test.jpg",
		}
		Expect(models.CreateLesion(ctx, dbConn, lesion)).To(BeNil())

		// Then insert the request
		input := &models.Request{
			AccountID: accountID,
			Status:    &models.StatusSubmitted,
		}
		Expect(models.CreateRequest(ctx, dbConn, input)).To(BeNil())

		report, err := models.GetRequest(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		Expect(report.ID).To(Equal(input.ID))
		Expect(report.AccountID).To(Equal(input.AccountID))
		Expect(report.Status).To(Equal(input.Status))

		report.Status = &models.StatusAnswered
		Expect(models.UpdateRequest(ctx, dbConn, report)).To(BeNil())
		Expect(*report.Status).To(Equal(models.StatusAnswered))

		_, err = models.DeleteRequest(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())
	})

	It("Report and report answers should correctly perform the CRUD flow", func() {
		ctx := context.Background()

		// First insert the account
		accountID, err := types.StringToUUID(uuid.NewV4().String())
		Expect(err).To(BeNil())
		Expect(models.UpsertAccount(ctx, dbConn, &models.Account{
			ID:     accountID,
			Name:   "First Last",
			Email:  "hello@world.com",
			Phone:  "+48123456789",
			Gender: "Male",
		})).To(BeNil())

		// Then the body part
		bodyPart := &models.BodyPart{
			Name:      "Left knee",
			Displayed: true,
			Image:     "https://example.org/test.jpg",
			Order:     3,
		}
		Expect(models.CreateBodyPart(ctx, dbConn, bodyPart)).To(BeNil())

		// Then insert the lesion
		lesion := &models.Lesion{
			AccountID:        accountID,
			Name:             "Mole on my knee",
			BodyPartID:       bodyPart.ID,
			BodyPartLocation: "https://example.org/test.jpg",
		}
		Expect(models.CreateLesion(ctx, dbConn, lesion)).To(BeNil())

		// Then insert the request
		request := &models.Request{
			AccountID: accountID,
			Status:    &models.StatusSubmitted,
		}
		Expect(models.CreateRequest(ctx, dbConn, request)).To(BeNil())

		// Then the question
		question := &models.Question{
			Name: "Do you smoke?",
			Type: "select",
			Answers: types.ExtendedJSONB{
				Status: pgtype.Present,
				Bytes:  []byte(`{"hello":"world"}`),
			},
			Displayed: true,
			Order:     3,
		}
		Expect(models.CreateQuestion(ctx, dbConn, question)).To(BeNil())

		// And the report
		strarr := &types.ExtendedStringArray{}
		strarr.Set([]string{"hello", "world"})
		input := &models.Report{
			RequestID: request.ID,
			LesionID:  lesion.ID,
			Photos:    *strarr,
		}
		Expect(models.CreateReport(ctx, dbConn, input)).To(BeNil())

		report, err := models.GetReport(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		Expect(report.ID).To(Equal(input.ID))
		Expect(report.RequestID).To(Equal(input.RequestID))
		Expect(report.LesionID).To(Equal(input.LesionID))
		Expect(report.Photos).To(Equal(input.Photos))

		strarr.Set([]string{"world", "hello"})
		report.Photos = *strarr
		Expect(models.UpdateReport(ctx, dbConn, report)).To(BeNil())
		Expect(report.Photos).To(Equal(*strarr))

		// Tackle the report answers
		answer1 := &models.ReportAnswer{
			ReportID:   report.ID,
			QuestionID: question.ID,
			Answer: types.ExtendedJSONB{
				Status: pgtype.Present,
				Bytes:  []byte(`{"hello": "world"}`),
			},
		}
		Expect(models.CreateReportAnswer(ctx, dbConn, answer1)).To(BeNil())

		answer2 := &models.ReportAnswer{
			ReportID:   report.ID,
			QuestionID: question.ID,
			Answer: types.ExtendedJSONB{
				Status: pgtype.Present,
				Bytes:  []byte(`{"hello": "world"}`),
			},
		}
		Expect(models.CreateReportAnswer(ctx, dbConn, answer2)).To(BeNil())

		answer1b, err := models.GetReportAnswer(ctx, dbConn, answer1.ID)
		Expect(err).To(BeNil())

		Expect(answer1b.ID).To(Equal(answer1.ID))
		Expect(answer1b.ReportID).To(Equal(answer1.ReportID))
		Expect(answer1b.QuestionID).To(Equal(answer1.QuestionID))
		Expect(answer1b.Answer).To(Equal(answer1.Answer))

		newJSON := types.ExtendedJSONB{
			Status: pgtype.Present,
			Bytes:  []byte(`{"world": "hello"}`),
		}
		answer1b.Answer = newJSON
		Expect(models.UpdateReportAnswer(ctx, dbConn, answer1b)).To(BeNil())
		Expect(answer1b.Answer).To(Equal(newJSON))

		// Both answers should be visible at once
		answers, err := models.ListReportAnswersByReportID(ctx, dbConn, report.ID)
		Expect(err).To(BeNil())
		Expect(len(answers) >= 2).To(BeTrue())

		_, err = models.DeleteReportAnswer(ctx, dbConn, answer1.ID)
		Expect(err).To(BeNil())

		_, err = models.DeleteReport(ctx, dbConn, input.ID)
		Expect(err).To(BeNil())

		_, err = models.GetReportAnswer(ctx, dbConn, answer2.ID)
		Expect(err).ToNot(BeNil())
	})
})
