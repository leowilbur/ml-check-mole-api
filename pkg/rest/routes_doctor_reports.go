package rest

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/models"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type doctorReport struct {
	*models.Report
	Lesion   *models.Lesion   `json:"lesion,omitempty"`
	BodyPart *models.BodyPart `json:"body_part,omitempty"`
	Answers  []*extAnswer     `json:"answers,omitempty"`
}

// DoctorReportsList allows doctors to list reports from the whole database
func (a *API) DoctorReportsList(r *gin.Context) {
	var (
		requestIDString  = r.Query("request_id")
		lesionIDString   = r.Query("lesion_id")
		offsetString     = r.Query("offset")
		limitString      = r.Query("limit")
		tempURLs         = r.Query("temp_urls") == "true"
		includeLesions   = r.Query("include_lesions") == "true"
		includeBodyParts = r.Query("include_body_parts") == "true"
		includeAnswers   = r.Query("include_answers") == "true"
		includeQuestions = r.Query("include_questions") == "true"
		orderBy          = r.DefaultQuery("order_by", "reports.created_at")
	)

	var (
		offset int
		limit  = 50
	)

	if parsedOffset, err := strconv.Atoi(offsetString); err == nil {
		offset = parsedOffset
	}

	if parsedLimit, err := strconv.Atoi(limitString); err == nil {
		limit = parsedLimit
	}

	query := psql.Select(
		"reports.id",
		"reports.request_id",
		"reports.lesion_id",
		"reports.photos",
		"reports.status",
		"reports.consultation_result",
		"reports.created_at",
		"reports.updated_at",
		"lesions.id",
		"lesions.account_id",
		"lesions.name",
		"lesions.body_part_id",
		"lesions.body_part_location",
		"lesions.created_at",
		"lesions.updated_at",
		"body_parts.id",
		"body_parts.name",
		"body_parts.displayed",
		"body_parts.image",
		"body_parts.order",
		"body_parts.parent",
	).From("reports").
		LeftJoin("lesions ON lesions.id = reports.lesion_id").
		LeftJoin("body_parts ON body_parts.id = lesions.body_part_id").
		OrderBy(orderBy)

	if requestIDString != "" {
		requestID, err := types.StringToUUID(requestIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request ID",
			})
			return
		}
		query = query.Where("reports.request_id = ?", &requestID)
	}

	if lesionIDString != "" {
		lesionID, err := types.StringToUUID(lesionIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid lesion ID",
			})
			return
		}
		query = query.Where("reports.lesion_id = ?", &lesionID)
	}

	query = query.Offset(uint64(offset)).Limit(uint64(limit))

	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("Error while building the ListLesions query", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to build the query",
		})
		return
	}

	rows, err := a.DB.QueryEx(r, sql, nil, args...)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to list reports",
		})
		return
	}

	result := []*doctorReport{}
	for rows.Next() {
		item := &doctorReport{
			Report:   &models.Report{},
			Lesion:   &models.Lesion{},
			BodyPart: &models.BodyPart{},
		}
		if err := rows.Scan(
			&item.Report.ID,
			&item.Report.RequestID,
			&item.Report.LesionID,
			&item.Report.Photos,
			&item.Report.Status,
			&item.Report.ConsultationResult,
			&item.Report.CreatedAt,
			&item.Report.UpdatedAt,
			&item.Lesion.ID,
			&item.Lesion.AccountID,
			&item.Lesion.Name,
			&item.Lesion.BodyPartID,
			&item.Lesion.BodyPartLocation,
			&item.Lesion.CreatedAt,
			&item.Lesion.UpdatedAt,
			&item.BodyPart.ID,
			&item.BodyPart.Name,
			&item.BodyPart.Displayed,
			&item.BodyPart.Image,
			&item.BodyPart.Order,
			&item.BodyPart.Parent,
		); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to scan requests",
			})
			return
		}

		if tempURLs {
			for i, photo := range item.Report.Photos.Elements {
				var val string
				if err := photo.AssignTo(&val); err != nil {
					r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Unable to scan the photo",
					})
					return
				}

				cpl, _ := url.Parse(val)
				pathParts := strings.SplitN(cpl.Path, "/", 2)

				req, _ := a.S3.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(pathParts[0]),
					Key:    aws.String(pathParts[1]),
				})
				urlStr, err := req.Presign(60 * time.Minute)
				if err != nil {
					r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Presigned URL generation failed - " + err.Error(),
					})
					return
				}

				item.Report.Photos.Elements[i].Set(urlStr)
			}
		}

		if !includeLesions {
			item.Lesion = nil
		}

		if !includeBodyParts {
			item.BodyPart = nil
		}

		if includeAnswers {
			item.Answers = []*extAnswer{}

			subrows, err := a.DB.QueryEx(
				r,
				`SELECT
					report_answers.id,
					report_answers.report_id,
					report_answers.question_id,
					report_answers.answer,
					questions.id,
					questions.name,
					questions.type,
					questions.answers,
					questions.displayed,
					questions.order,
					questions.created_at,
					questions.updated_at
				FROM report_answers
				LEFT JOIN questions ON questions.id = report_answers.question_id
				WHERE report_answers.report_id = $1
				ORDER BY questions.order ASC`,
				nil,
				&item.ID,
			)
			if err != nil {
				log.Println("Error while querying the answers", err)
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to query the answers",
				})
				return
			}

			for subrows.Next() {
				subitem := &extAnswer{
					ReportAnswer: &models.ReportAnswer{},
					Question:     &models.Question{},
				}

				if err := subrows.Scan(
					&subitem.ReportAnswer.ID,         // report_answers.id,
					&subitem.ReportAnswer.ReportID,   // report_answers.report_id,
					&subitem.ReportAnswer.QuestionID, // report_answers.question_id,
					&subitem.ReportAnswer.Answer,     // report_answers.answer,
					&subitem.Question.ID,             // questions.id,
					&subitem.Question.Name,           // questions.name,
					&subitem.Question.Type,           // questions.type,
					&subitem.Question.Answers,        // questions.answers,
					&subitem.Question.Displayed,      // questions.displayed,
					&subitem.Question.Order,          // questions.order,
					&subitem.Question.CreatedAt,      // questions.created_at,
					&subitem.Question.UpdatedAt,      // questions.updated_at
				); err != nil {
					r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Unable to scan the answers",
					})
					return
				}

				if !includeQuestions {
					subitem.Question = nil
				}

				item.Answers = append(item.Answers, subitem)
			}
		}

		result = append(result, item)
	}

	r.JSON(http.StatusOK, result)
}

// DoctorReportsGet allows doctors to get any report in the database
func (a *API) DoctorReportsGet(r *gin.Context) {
	var (
		includeLesion    = r.Query("include_lesion") == "true"
		includeBodyPart  = r.Query("include_body_part") == "true"
		includeAnswers   = r.Query("include_answers") == "true"
		includeQuestions = r.Query("include_questions") == "true"
	)

	idString := r.Param("report")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	report, err := models.GetReport(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Report not found",
		})
		return
	}

	var lesion *models.Lesion
	if includeLesion {
		lesion, err = models.GetLesion(r, a.DB, report.LesionID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Lesion not found",
			})
			return
		}
	}

	var bodyPart *models.BodyPart
	if includeLesion && includeBodyPart {
		bodyPart, err = models.GetBodyPart(r, a.DB, lesion.BodyPartID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Body part not found",
			})
			return
		}
	}

	var answers []*extAnswer
	if includeAnswers {
		subrows, err := a.DB.QueryEx(
			r,
			`SELECT
					report_answers.id,
					report_answers.report_id,
					report_answers.question_id,
					report_answers.answer,
					questions.id,
					questions.name,
					questions.type,
					questions.answers,
					questions.displayed,
					questions.order,
					questions.created_at,
					questions.updated_at
				FROM report_answers
				LEFT JOIN questions ON questions.id = report_answers.question_id
				WHERE report_answers.report_id = $1
				ORDER BY questions.order ASC`,
			nil,
			&report.ID,
		)
		if err != nil {
			log.Println("Error while querying the answers", err)
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to query the answers",
			})
			return
		}

		for subrows.Next() {
			subitem := &extAnswer{
				ReportAnswer: &models.ReportAnswer{},
				Question:     &models.Question{},
			}

			if err := subrows.Scan(
				&subitem.ReportAnswer.ID,         // report_answers.id,
				&subitem.ReportAnswer.ReportID,   // report_answers.report_id,
				&subitem.ReportAnswer.QuestionID, // report_answers.question_id,
				&subitem.ReportAnswer.Answer,     // report_answers.answer,
				&subitem.Question.ID,             // questions.id,
				&subitem.Question.Name,           // questions.name,
				&subitem.Question.Type,           // questions.type,
				&subitem.Question.Answers,        // questions.answers,
				&subitem.Question.Displayed,      // questions.displayed,
				&subitem.Question.Order,          // questions.order,
				&subitem.Question.CreatedAt,      // questions.created_at,
				&subitem.Question.UpdatedAt,      // questions.updated_at
			); err != nil {
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to scan the answers",
				})
				return
			}

			if !includeQuestions {
				subitem.Question = nil
			}

			answers = append(answers, subitem)
		}
	}

	r.JSON(http.StatusOK, doctorReport{
		Report:   report,
		Lesion:   lesion,
		BodyPart: bodyPart,
		Answers:  answers,
	})
}
