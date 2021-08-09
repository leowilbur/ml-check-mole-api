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
	"github.com/jackc/pgx/pgtype"
)

// RequestsList returns all requests owned by the user
func (a *API) RequestsList(r *gin.Context) {
	var (
		offsetString = r.Query("offset")
		limitString  = r.Query("limit")
		skipAnswer   = r.Query("skip_answer") == "true"
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
			id, account_id, status, answer_text, answered_by, answered_at,
			created_at, updated_at
		FROM requests WHERE account_id = $1 ORDER BY created_at DESC
		OFFSET $2 LIMIT $3`,
		nil,
		&r.MustGet("account").(*models.Account).ID,
		offset, limit,
	)
	if err != nil {
		log.Println(err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to list requests",
		})
		return
	}

	result := []*models.Request{}
	for rows.Next() {
		item := &models.Request{}
		if err := rows.Scan(
			&item.ID,
			&item.AccountID,
			&item.Status,
			&item.AnswerText,
			&item.AnsweredBy,
			&item.AnsweredAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Println(err)
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to scan requests",
			})
			return
		}

		if skipAnswer {
			item.AnswerText = types.ExtendedText{
				Status: pgtype.Null,
			}
		}

		result = append(result, item)
	}

	r.JSON(http.StatusOK, result)
}

// RequestsGet returns a request by its id
func (a *API) RequestsGet(r *gin.Context) {
	var (
		includeReports   = r.Query("include_reports") == "true"
		includeLesions   = r.Query("include_lesions") == "true"
		includeAnswers   = r.Query("include_answers") == "true"
		includeQuestions = r.Query("include_questions") == "true"
		skipAnswer       = r.Query("skip_answer") == "true"
		tempURLs         = r.Query("temp_urls") == "true"
	)

	idString := r.Param("request")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the request
	request, err := models.GetRequest(r, a.DB, id)
	if err != nil ||
		request.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Request not found",
		})
		return
	}
	if skipAnswer {
		request.AnswerText = types.ExtendedText{
			Status: pgtype.Null,
		}
	}

	var reports []*extReport
	if includeReports {
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
				reports.updated_at,
				lesions.id,
				lesions.account_id,
				lesions.name,
				lesions.body_part_id,
				lesions.body_part_location,
				lesions.created_at,
				lesions.updated_at
			FROM reports
			LEFT JOIN lesions ON lesions.id = reports.lesion_id
			WHERE reports.request_id = $1
			ORDER BY reports.created_at DESC`,
			nil,
			&request.ID,
		)
		if err != nil {
			log.Println(err)
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to query the reports",
			})
			return
		}

		reports = []*extReport{}
		for rows.Next() {
			item := &extReport{
				Report: &models.Report{},
				Lesion: &models.Lesion{},
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
				&item.Lesion.ID,                 // lesions.id
				&item.Lesion.AccountID,          // lesions.account_id
				&item.Lesion.Name,               // lesions.name
				&item.Lesion.BodyPartID,         // lesions.body_part_id
				&item.Lesion.BodyPartLocation,   // lesions.body_part_location
				&item.Lesion.CreatedAt,          // lesions.created_at
				&item.Lesion.UpdatedAt,          // lesions.updated_at
			); err != nil {
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to scan the records",
				})
				return
			}

			if !includeLesions {
				item.Lesion = nil
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

			reports = append(reports, item)
		}
	}

	r.JSON(http.StatusOK, struct {
		*models.Request
		Reports []*extReport `json:"reports,omitempty"`
	}{
		Request: request,
		Reports: reports,
	})
}

// RequestsCreate creates a new request
func (a *API) RequestsCreate(r *gin.Context) {
	var input struct {
		Reports []types.ExtendedUUID `json:"reports"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer tx.Rollback()

	status := models.RequestStatus("Open")

	item := &models.Request{
		AccountID: r.MustGet("account").(*models.Account).ID,
		Status:    &status,
	}
	if err := models.CreateRequest(r, tx, item); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the request",
		})
		return
	}

	for _, id := range input.Reports {
		report, err := models.GetReport(r, tx, id)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to fetch the report: " + err.Error(),
			})
			return
		}

		lesion, err := models.GetLesion(r, tx, report.LesionID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to fetch the lesion: " + err.Error(),
			})
			return
		}

		if lesion.AccountID.Bytes != item.AccountID.Bytes {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Every report in the update has to belong to the user",
			})
			return
		}

		report.RequestID = item.ID
		report.Status = types.StringToText("Submitted")
		if err := models.UpdateReport(r, tx, report); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to save a report",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to complete the transaction",
		})
	}

	r.JSON(http.StatusCreated, item)
}

// RequestsUpdate updates an existing request
func (a *API) RequestsUpdate(r *gin.Context) {
	idString := r.Param("request")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}
	var input struct {
		Status  string               `json:"status"`
		Reports []types.ExtendedUUID `json:"reports"`
	}
	if err := r.BindJSON(&input); err != nil ||
		(input.Status != "Submitted" && input.Status != "Open") {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	// Verify ownership of the request
	request, err := models.GetRequest(r, a.DB, id)
	if err != nil ||
		request.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Request not found",
		})
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer tx.Rollback()

	status := models.RequestStatus(input.Status)
	request.Status = &status
	if err := models.UpdateRequest(r, tx, request); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    201,
			"message": "Invalid input",
		})
		return
	}

	for _, id := range input.Reports {
		report, err := models.GetReport(r, tx, id)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		lesion, err := models.GetLesion(r, tx, report.LesionID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if lesion.AccountID.Bytes != request.AccountID.Bytes {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Every report in the update has to belong to the user",
			})
			return
		}

		report.RequestID = request.ID
		report.Status = types.StringToText("Submitted")
		if err := models.UpdateReport(r, tx, report); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to save a report",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to complete the transaction",
		})
	}

	r.JSON(http.StatusOK, request)
}

// RequestsDelete deletes an existing request
func (a *API) RequestsDelete(r *gin.Context) {
	idString := r.Param("request")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	// Verify ownership of the request
	request, err := models.GetRequest(r, a.DB, id)
	if err != nil ||
		request.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Request not found",
		})
		return
	}

	item, err := models.DeleteRequest(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Request not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}
