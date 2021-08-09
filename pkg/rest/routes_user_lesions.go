package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

type extLesion struct {
	*models.Lesion
	BodyPart    *models.BodyPart `json:"body_part,omitempty"`
	LastReport  *models.Report   `json:"last_report,omitempty"`
	LastRequest *models.Request  `json:"last_request,omitempty"`
}

// LesionsList returns info about lesions owned by the user
func (a *API) LesionsList(r *gin.Context) {
	// Figure out the request options
	var (
		includeBodyParts    = r.Query("include_body_parts") == "true"
		includeLastReports  = r.Query("include_last_reports") == "true"
		includeLastRequests = r.Query("include_last_requests") == "true"
		offsetString        = r.Query("offset")
		limitString         = r.Query("limit")
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
			lesions.lid,
			lesions.laccount_id,
			lesions.lname,
			lesions.lbody_part_id,
			lesions.lbody_part_location,
			lesions.lcreated_at,
			lesions.lupdated_at,
			lesions.bid,
			lesions.bname,
			lesions.bdisplayed,
			lesions.bimage,
			lesions.border,
			lesions.bparent,
			reports.id,
			reports.request_id,
			reports.lesion_id,
			reports.photos,
			reports.status,
			reports.consultation_result,
			reports.created_at,
			reports.updated_at,
			requests.id,
			requests.account_id,
			requests.status,
			requests.answer_text,
			requests.answered_by,
			requests.answered_at,
			requests.created_at,
			requests.updated_at
		FROM (
			SELECT
				lesions.id AS lid,
				lesions.account_id AS laccount_id,
				lesions.name AS lname,
				lesions.body_part_id AS lbody_part_id,
				lesions.body_part_location AS lbody_part_location,
				lesions.created_at AS lcreated_at,
				lesions.updated_at AS lupdated_at,
				body_parts.id AS bid,
				body_parts.name AS bname,
				body_parts.displayed AS bdisplayed,
				body_parts.image AS bimage,
				body_parts.order AS border,
				body_parts.parent AS bparent
			FROM lesions
			LEFT JOIN body_parts ON body_parts.id = lesions.body_part_id
			WHERE lesions.account_id = $1
		  ORDER BY lesions.created_at DESC
		  OFFSET $2 LIMIT $3
		) AS lesions
		LEFT JOIN LATERAL (
			SELECT
				id, request_id, lesion_id, photos, status, consultation_result, created_at, updated_at
			FROM reports
			WHERE reports.lesion_id = lesions.lid
			ORDER BY created_at DESC
			LIMIT 1
		) AS reports ON true
		LEFT JOIN LATERAL (
			SELECT
				id, account_id, status, answer_text, answered_by, answered_at, created_at, updated_at
			FROM requests
			WHERE requests.id = reports.request_id
			ORDER BY created_at DESC
			LIMIT 1
		) AS requests ON true`,
		nil,
		&r.MustGet("account").(*models.Account).ID,
		offset,
		limit,
	)
	if err != nil {
		log.Print(err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to query the lesions, " + err.Error(),
		})
		return
	}

	result := []*extLesion{}
	for rows.Next() {
		item := &extLesion{
			Lesion:      &models.Lesion{},
			BodyPart:    &models.BodyPart{},
			LastReport:  &models.Report{},
			LastRequest: &models.Request{},
		}
		if err := rows.Scan(
			&item.Lesion.ID,                     // lesions.id
			&item.Lesion.AccountID,              // lesions.account_id
			&item.Lesion.Name,                   // lesions.name
			&item.Lesion.BodyPartID,             // lesions.body_part_id
			&item.Lesion.BodyPartLocation,       // lesions.body_part_location
			&item.Lesion.CreatedAt,              // lesions.created_at
			&item.Lesion.UpdatedAt,              // lesions.updated_at
			&item.BodyPart.ID,                   // body_parts.id
			&item.BodyPart.Name,                 // body_parts.name
			&item.BodyPart.Displayed,            // body_parts.displayed
			&item.BodyPart.Image,                // body_parts.image
			&item.BodyPart.Order,                // body_parts.order
			&item.BodyPart.Parent,               // body_parts.parent
			&item.LastReport.ID,                 // reports.id
			&item.LastReport.RequestID,          // reports.request_id
			&item.LastReport.LesionID,           // reports.lesion_id
			&item.LastReport.Photos,             // reports.photos
			&item.LastReport.Status,             // reports.status
			&item.LastReport.ConsultationResult, // reports.consultation_result
			&item.LastReport.CreatedAt,          // reports.created_at
			&item.LastReport.UpdatedAt,          // reports.updated_at
			&item.LastRequest.ID,                // requests.id
			&item.LastRequest.AccountID,         // requests.account_id
			&item.LastRequest.Status,            // requests.status
			&item.LastRequest.AnswerText,        // requests.answer_text
			&item.LastRequest.AnsweredBy,        // requests.answered_by
			&item.LastRequest.AnsweredAt,        // requests.answered_at
			&item.LastRequest.CreatedAt,         // requests.created_at
			&item.LastRequest.UpdatedAt,         // requests.updated_at
		); err != nil {
			log.Println(err)
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    201,
				"message": err.Error(),
			})
			return
		}

		if !includeBodyParts {
			item.BodyPart = nil
		}
		if !includeLastReports {
			item.LastReport = nil
		}
		if !includeLastRequests {
			item.LastRequest = nil
		}

		result = append(result, item)
	}

	r.JSON(http.StatusOK, result)
}

// LesionsCreate inserts a new lesion owned by the user
func (a *API) LesionsCreate(r *gin.Context) {
	var input struct {
		Name             string              `json:"name"`
		BodyPartID       *types.ExtendedUUID `json:"body_part_id"`
		BodyPartLocation string              `json:"body_part_location"`
	}
	if err := r.BindJSON(&input); err != nil || input.BodyPartID == nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input",
		})
		return
	}

	item := &models.Lesion{
		AccountID:        r.MustGet("account").(*models.Account).ID,
		Name:             input.Name,
		BodyPartID:       *input.BodyPartID,
		BodyPartLocation: input.BodyPartLocation,
	}
	if err := models.CreateLesion(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the lesion",
		})
		return
	}

	r.JSON(http.StatusCreated, item)
}

// LesionsUpdate updates a specific lesion owned by the user
func (a *API) LesionsUpdate(r *gin.Context) {
	idString := r.Param("lesion")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	if existing, err := models.GetLesion(r, a.DB, id); err != nil ||
		existing.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	var input struct {
		Name             string             `json:"name"`
		BodyPartID       types.ExtendedUUID `json:"body_part_id"`
		BodyPartLocation string             `json:"body_part_location"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input: " + err.Error(),
		})
		return
	}

	item := &models.Lesion{
		ID:               id,
		AccountID:        r.MustGet("account").(*models.Account).ID,
		Name:             input.Name,
		BodyPartID:       input.BodyPartID,
		BodyPartLocation: input.BodyPartLocation,
	}
	if err := models.UpdateLesion(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}

// LesionsDelete deletes a specific lesion owned by the user
func (a *API) LesionsDelete(r *gin.Context) {
	idString := r.Param("lesion")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	if existing, err := models.GetLesion(r, a.DB, id); err != nil ||
		existing.AccountID.Bytes != r.MustGet("account").(*models.Account).ID.Bytes {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	item, err := models.DeleteLesion(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}
