package rest

import (
	"bytes"
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
	"github.com/jackc/pgx/pgtype"
)

// ReportsCreate inserts a new report owned by the user, belonging to the lesion
func (a *API) ReportsCreate(r *gin.Context) {
	lesionString := r.Param("lesion")
	lesion, err := types.StringToUUID(lesionString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid lesion ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the lesion
	if lesion, err := models.GetLesion(r, a.DB, lesion); err != nil ||
		lesion.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	var input struct {
		RequestID *types.ExtendedUUID       `json:"request_id"`
		Photos    types.ExtendedStringArray `json:"photos"`
		Answers   []struct {
			QuestionID types.ExtendedUUID  `json:"question_id"`
			Answer     types.ExtendedJSONB `json:"answer"`
		} `json:"answers"`
		Status types.ExtendedText `json:"status"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to insert the report",
		})
		return
	}
	defer tx.Rollback()

	// First create the report
	item := &models.Report{
		LesionID: lesion,
		Photos:   input.Photos,
		Status:   input.Status,
	}

	if item.Status.Status == pgtype.Undefined {
		item.Status.Status = pgtype.Null
	}

	if input.RequestID != nil && input.RequestID.Status == pgtype.Present {
		// Verify ownership of the request
		if request, err := models.GetRequest(r, a.DB, *input.RequestID); err != nil ||
			request.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Request not found",
			})
			return
		}

		item.RequestID = *input.RequestID
	} else {
		item.RequestID = types.ExtendedUUID{
			Status: pgtype.Null,
		}
	}

	if err := models.CreateReport(r, tx, item); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to create the report",
		})
		return
	}

	// Then create the answers
	for _, answer := range input.Answers {
		subitem := &models.ReportAnswer{
			ReportID:   item.ID,
			QuestionID: answer.QuestionID,
			Answer:     answer.Answer,
		}

		if err := models.CreateReportAnswer(r, tx, subitem); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    201,
				"message": "Invalid input",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the report",
		})
	}

	r.JSON(http.StatusCreated, item)
}

type extReport struct {
	*models.Report
	Answers []*extAnswer   `json:"answers,omitempty"`
	Lesion  *models.Lesion `json:"lesion,omitempty"`
}
type extAnswer struct {
	*models.ReportAnswer
	Question *models.Question `json:"question,omitempty"`
}

// ReportsList returns info about reports owned by the user, belonging to a lesion
func (a *API) ReportsList(r *gin.Context) {
	lesionString := r.Param("lesion")
	lesion, err := types.StringToUUID(lesionString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid lesion ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the lesion
	if lesion, err := models.GetLesion(r, a.DB, lesion); err != nil ||
		lesion.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	// Figure out the request options
	var (
		includeAnswers   = r.Query("include_answers") == "true"
		includeQuestions = r.Query("include_questions") == "true"
		tempURLs         = r.Query("temp_urls") == "true"
		offsetString     = r.Query("offset")
		limitString      = r.Query("limit")
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

	rows, err := a.DB.QueryEx(
		r,
		`SELECT
			reports.id,
			reports.request_id,
			reports.lesion_id,
			reports.photos,
			reports.status,
			reports.consultation_result,
			reports.created_at,
			reports.updated_at
		FROM reports WHERE lesion_id = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3`,
		nil,
		&lesion,
		offset,
		limit,
	)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to query the reports",
		})
		return
	}

	result := []*extReport{}
	for rows.Next() {
		item := &extReport{
			Report: &models.Report{},
		}

		if err := rows.Scan(
			&item.Report.ID,                 // reports.id
			&item.Report.RequestID,          // reports.request_id
			&item.Report.LesionID,           // reports.lesion_id
			&item.Report.Photos,             // reports.photos
			&item.Report.Status,             // reports.status
			&item.Report.ConsultationResult, // reports.consultation_result
			&item.Report.CreatedAt,          // reports.created_at
			&item.Report.UpdatedAt,          // reports.updated_at
		); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to scan the records",
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

// ReportsUpdate updates a specific report owned by the user
func (a *API) ReportsUpdate(r *gin.Context) {
	lesionString := r.Param("lesion")
	lesionID, err := types.StringToUUID(lesionString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid lesion ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the lesion
	lesion, err := models.GetLesion(r, a.DB, lesionID)
	if err != nil ||
		lesion.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	reportString := r.Param("report")
	reportID, err := types.StringToUUID(reportString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid report ID supplied: " + err.Error(),
		})
		return
	}

	var input struct {
		RequestID *types.ExtendedUUID       `json:"request_id"`
		Photos    types.ExtendedStringArray `json:"photos"`
		Status    types.ExtendedText        `json:"status"`
		Answers   []struct {
			QuestionID types.ExtendedUUID  `json:"question_id"`
			Answer     types.ExtendedJSONB `json:"answer"`
		} `json:"answers"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Verify ownership of the request
	if input.RequestID != nil && input.RequestID.Status == pgtype.Present {
		request, err := models.GetRequest(r, a.DB, *input.RequestID)
		if err != nil ||
			request.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Request not found",
			})
			return
		}
	}

	// Verify ownership of the report
	report, err := models.GetReport(r, a.DB, reportID)
	if err != nil ||
		report.LesionID != lesion.ID {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Report not found",
		})
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		log.Println("Unable to enter the transaction", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the report",
		})
		return
	}
	defer tx.Rollback()

	// Now that we have sorted out the ACL, we can start properly processing
	// the request.
	if input.RequestID != nil && input.RequestID.Status == pgtype.Present {
		report.RequestID = *input.RequestID
	} else {
		report.RequestID = types.ExtendedUUID{
			Status: pgtype.Null,
		}
	}
	report.Photos = input.Photos
	report.Status = input.Status
	if err := models.UpdateReport(r, tx, report); err != nil {
		log.Println("Unable to update the report", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the report",
		})
		return
	}

	// Then go over the existing answers and run a diff
	existingAnswers, err := models.ListReportAnswersByReportID(r, a.DB, report.ID)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to acquire the existing answers",
		})
		return
	}

	var (
		oldAnswers   = map[types.ExtendedUUID]types.ExtendedJSONB{}
		newAnswers   = map[types.ExtendedUUID]struct{}{}
		oldAnswerIDs = map[types.ExtendedUUID]types.ExtendedUUID{}
	)
	for _, answer := range existingAnswers {
		oldAnswers[answer.QuestionID] = answer.Answer
		oldAnswerIDs[answer.QuestionID] = answer.ID
	}

	for _, answer := range input.Answers {
		newAnswers[answer.QuestionID] = struct{}{}

		matched, ok := oldAnswers[answer.QuestionID]
		if !ok {
			// Create an answer
			if err := models.CreateReportAnswer(r, tx, &models.ReportAnswer{
				ReportID:   reportID,
				QuestionID: answer.QuestionID,
				Answer:     answer.Answer,
			}); err != nil {
				log.Println("Error while creating an answer", err)
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to save an answer. Are you sure that such a question exists?",
				})
				return
			}

			continue
		}

		if !bytes.Equal(matched.Bytes, answer.Answer.Bytes) {
			// Update the answer
			if err := models.UpdateReportAnswer(r, tx, &models.ReportAnswer{
				ID:         oldAnswerIDs[answer.QuestionID],
				ReportID:   reportID,
				QuestionID: answer.QuestionID,
				Answer:     answer.Answer,
			}); err != nil {
				log.Println("Error while updating an answer", err)
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to update an answer. Are you sure that such a question exists?",
				})
				return
			}
		}
	}

	// 2nd pass to do deletes
	for _, answer := range existingAnswers {
		if _, ok := newAnswers[answer.QuestionID]; ok {
			continue
		}

		// Delete the answer
		if _, err := models.DeleteReportAnswer(r, tx, answer.ID); err != nil {
			log.Println("Error while deleting an answer", err)
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to delete a removed answer.",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the report",
		})
	}

	r.JSON(http.StatusOK, report)
}

// ReportsDelete deletes a specific report owned by the user
func (a *API) ReportsDelete(r *gin.Context) {
	lesionString := r.Param("lesion")
	lesion, err := types.StringToUUID(lesionString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid lesion ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the lesion
	lesionModel, err := models.GetLesion(r, a.DB, lesion)
	if err != nil || lesionModel.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	reportString := r.Param("report")
	report, err := types.StringToUUID(reportString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid report ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the request
	if reportModel, err := models.GetReport(r, a.DB, report); err != nil ||
		reportModel.LesionID != lesionModel.ID {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Report not found",
		})
		return
	}

	item, err := models.DeleteReport(r, a.DB, report)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Report not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}
