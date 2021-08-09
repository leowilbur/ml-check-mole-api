package rest

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	fcm "github.com/appleboy/go-fcm"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgtype"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

type extRequest struct {
	*models.Request
	Account *models.Account `json:"account,omitempty"`
}

type extResponse struct {
	Data       []*extRequest `json:"data"`
	TotalCount int32         `json:"total"`
}

func tryParse(key string, input string) interface{} {
	i, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		if strings.HasSuffix(key, "_at") {
			return time.Unix(i, 0)
		}

		return i
	}

	b, err := strconv.ParseBool(input)
	if err == nil {
		return b
	}

	t, err := time.Parse(time.RFC3339, input)
	if err == nil {
		return t
	}

	return input
}

// DoctorRequestsList allows doctors to list requests from the whole database
func (a *API) DoctorRequestsList(r *gin.Context) {
	var (
		accountIDString = r.Query("account_id")
		status          = r.Query("status")
		includeAccounts = r.Query("include_accounts") == "true"
		skipAnswer      = r.Query("skip_answer") == "true"
		offsetString    = r.Query("offset")
		limitString     = r.Query("limit")
		orderBy         = r.DefaultQuery("order_by", "requests.created_at")
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
		"requests.id",
		"requests.account_id",
		"requests.status",
		"requests.answer_text",
		"requests.answered_by",
		"requests.answered_at",
		"requests.created_at",
		"requests.updated_at",
		"accounts.id",
		"accounts.name",
		"accounts.email",
		"accounts.phone",
		"accounts.gender",
		"accounts.birth_date",
		"accounts.created_at",
		"accounts.updated_at",
	).From("requests").
		LeftJoin("accounts ON accounts.id = requests.account_id").
		OrderBy(orderBy)

	query = query.Columns(`COUNT(*) OVER() AS "total_count"`)

	if accountIDString != "" {
		accountID, err := types.StringToUUID(accountIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid account ID",
			})
			return
		}
		query = query.Where("requests.account_id = ?", &accountID)
	}

	if status != "" {
		query = query.Where("requests.status = ?", status)
	}

	if filters := r.QueryMap("filters"); len(filters) > 0 {
		for key, value := range filters {
			var (
				parts   = strings.SplitN(value, ":", 2)
				desired = []interface{}{}
				symbol  string
			)
			switch parts[0] {
			case "eq":
				symbol = "= ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "like":
				symbol = "LIKE ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "ilike":
				symbol = "ILIKE ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "contains":
				symbol = "LIKE ?"
				desired = append(desired, "%"+parts[1]+"%")
			case "icontains":
				symbol = "ILIKE ?"
				desired = append(desired, "%"+parts[1]+"%")
			case "ne":
				symbol = "!= ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "gt":
				symbol = "> ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "gte":
				symbol = ">= ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "lt":
				symbol = "< ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "lte":
				symbol = "<= ?"
				desired = append(desired, tryParse(key, parts[1]))
			case "between":
				symbol = "BETWEEN ? AND ?"
				subparts := strings.SplitN(parts[1], ":", 2)
				desired = append(desired, tryParse(key, subparts[0]), tryParse(key, subparts[1]))
			}

			query = query.Where(key+" "+symbol, desired...)
		}
	}

	query = query.Offset(uint64(offset)).Limit(uint64(limit))

	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("Error while building the ListRequests query", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to build the query",
		})
		return
	}

	rows, err := a.DB.QueryEx(r, sql, nil, args...)
	if err != nil {
		log.Println("Error while executing the ListRequests query", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to list requests, " + err.Error(),
		})
		return
	}

	response := &extResponse{
		Data: []*extRequest{},
	}

	for rows.Next() {
		item := &extRequest{
			Request: &models.Request{},
			Account: &models.Account{},
		}

		if err := rows.Scan(
			&item.Request.ID,
			&item.Request.AccountID,
			&item.Request.Status,
			&item.Request.AnswerText,
			&item.Request.AnsweredBy,
			&item.Request.AnsweredAt,
			&item.Request.CreatedAt,
			&item.Request.UpdatedAt,
			&item.Account.ID,
			&item.Account.Name,
			&item.Account.Email,
			&item.Account.Phone,
			&item.Account.Gender,
			&item.Account.BirthDate,
			&item.Account.CreatedAt,
			&item.Account.UpdatedAt,
			&response.TotalCount,
		); err != nil {
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

		if !includeAccounts {
			item.Account = nil
		}

		response.Data = append(response.Data, item)
	}

	r.JSON(http.StatusOK, response)
}

// DoctorRequestsGet allows doctors to get any request in the database
func (a *API) DoctorRequestsGet(r *gin.Context) {
	var (
		includeAccount = r.Query("include_account") == "true"
	)

	idString := r.Param("request")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	request, err := models.GetRequest(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Request not found",
		})
		return
	}

	var account *models.Account
	if includeAccount {
		account, err = models.GetAccount(r, a.DB, request.AccountID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}
	}

	r.JSON(http.StatusOK, struct {
		*models.Request
		Account *models.Account `json:"account,omitempty"`
	}{
		Request: request,
		Account: account,
	})
}

// DoctorRequestsRespond allows doctors to respond to any request in the database
func (a *API) DoctorRequestsRespond(r *gin.Context) {
	idString := r.Param("request")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	var input struct {
		Status     string                  `json:"status"`
		AnswerText types.ExtendedText      `json:"answer_text"`
		AnsweredBy types.ExtendedText      `json:"answered_by"`
		AnsweredAt types.ExtendedTimestamp `json:"answered_at"`
		NotifyMsg  *string                 `json:"notify_msg"`
		Reports    []struct {
			ID                 string             `json:"id"`
			ConsultationResult types.ExtendedText `json:"consultation_result"`
		} `json:"reports"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input, " + err.Error(),
		})
		return
	}

	// Verify ownership of the request
	request, err := models.GetRequest(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Request not found",
		})
		return
	}

	status := models.RequestStatus(input.Status)

	request.Status = &status

	if input.AnswerText.Status == pgtype.Undefined {
		input.AnswerText.Status = pgtype.Null
	}
	request.AnswerText = input.AnswerText
	if input.AnsweredBy.Status == pgtype.Undefined {
		input.AnsweredBy.Status = pgtype.Null
	}
	request.AnsweredBy = input.AnsweredBy
	if input.AnsweredAt.Status == pgtype.Undefined {
		input.AnsweredAt.Status = pgtype.Null
	} else {
		input.AnsweredAt.Time = input.AnsweredAt.Time.UTC()
	}
	request.AnsweredAt = input.AnsweredAt

	tx, err := a.DB.Begin()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer tx.Rollback()

	for _, report := range input.Reports {
		reportID, err := types.StringToUUID(report.ID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid report ID supplied: " + err.Error(),
			})
			return
		}

		existing, err := models.GetReport(r, tx, reportID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if existing.RequestID.Bytes != request.ID.Bytes {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Every report in the update has to belong to the request",
			})
			return
		}

		existing.ConsultationResult = report.ConsultationResult
		if err := models.UpdateReport(r, tx, existing); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to save a report",
			})
			return
		}
	}

	if err := models.UpdateRequest(r, tx, request); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the request: " + err.Error(),
		})
		return
	}

	if err := tx.Commit(); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to complete the transaction",
		})
		return
	}

	if input.NotifyMsg != nil {
		if _, err := a.Config.FCMClient.SendWithRetry(&fcm.Message{
			To: request.AccountID.String(),
			Data: map[string]interface{}{
				"message": *input.NotifyMsg,
			},
		}, 3); err != nil {
			log.Printf(
				"Unable to send a FCM notification to %s: %s",
				request.AccountID.String(),
				err.Error(),
			)
		}
	}

	r.JSON(http.StatusOK, request)
}
