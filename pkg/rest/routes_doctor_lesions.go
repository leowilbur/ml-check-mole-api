package rest

import (
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/models"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	"github.com/gin-gonic/gin"

	sq "github.com/Masterminds/squirrel"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

// DoctorLesionsList allows doctors to list lesions from the whole database
func (a *API) DoctorLesionsList(r *gin.Context) {
	var (
		includeBodyParts = r.Query("include_body_parts") == "true"
		accountIDString  = r.Query("account_id")
		bodyPartIDString = r.Query("body_part_id")
		offsetString     = r.Query("offset")
		limitString      = r.Query("limit")
		orderBy          = r.DefaultQuery("order_by", "lesions.created_at")
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
	).From("lesions").
		LeftJoin("body_parts ON body_parts.id = lesions.body_part_id").
		OrderBy(orderBy)

	if accountIDString != "" {
		accountID, err := types.StringToUUID(accountIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid account ID",
			})
			return
		}
		query = query.Where("lesions.account_id = ?", &accountID)
	}

	if bodyPartIDString != "" {
		bodyPartID, err := types.StringToUUID(bodyPartIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid body part ID",
			})
			return
		}
		query = query.Where("lesions.body_part_id = ?", &bodyPartID)
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
		log.Println("Error while listing lesions", err)
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to list lesions",
		})
		return
	}

	result := []*extLesion{}
	for rows.Next() {
		item := &extLesion{
			Lesion:   &models.Lesion{},
			BodyPart: &models.BodyPart{},
		}
		if err := rows.Scan(
			&item.Lesion.ID,               // lesions.id
			&item.Lesion.AccountID,        // lesions.account_id
			&item.Lesion.Name,             // lesions.name
			&item.Lesion.BodyPartID,       // lesions.body_part_id
			&item.Lesion.BodyPartLocation, // lesions.body_part_location
			&item.Lesion.CreatedAt,        // lesions.created_at
			&item.Lesion.UpdatedAt,        // lesions.updated_at
			&item.BodyPart.ID,             // body_parts.id
			&item.BodyPart.Name,           // body_parts.name
			&item.BodyPart.Displayed,      // body_parts.displayed
			&item.BodyPart.Image,          // body_parts.image
			&item.BodyPart.Order,          // body_parts.order
			&item.BodyPart.Parent,         // body_parts.parent
		); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to scan requests",
			})
			return
		}

		if !includeBodyParts {
			item.BodyPart = nil
		}

		result = append(result, item)
	}

	r.JSON(http.StatusOK, result)
}

// DoctorLesionsGet allows doctors to get any lesion in the database
func (a *API) DoctorLesionsGet(r *gin.Context) {
	var (
		includeBodyPart = r.Query("include_body_part") == "true"
	)

	idString := r.Param("lesion")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	lesion, err := models.GetLesion(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Lesion not found",
		})
		return
	}

	var bodyPart *models.BodyPart
	if includeBodyPart {
		bodyPart, err = models.GetBodyPart(r, a.DB, lesion.BodyPartID)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Body part not found",
			})
			return
		}
	}

	r.JSON(http.StatusOK, struct {
		*models.Lesion
		BodyPart *models.BodyPart `json:"body_part,omitempty"`
	}{
		Lesion:   lesion,
		BodyPart: bodyPart,
	})
}
