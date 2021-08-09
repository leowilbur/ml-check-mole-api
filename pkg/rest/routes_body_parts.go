package rest

import (
	"net/http"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/models"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

// BodyPartsList lists all body parts ordered by "order"
func (a *API) BodyPartsList(r *gin.Context) {
	rows, err := a.DB.QueryEx(
		r,
		`SELECT
			id, name, displayed, image, "order", parent
		FROM body_parts ORDER BY "order" ASC`,
		nil,
	)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch the body parts",
		})
		return
	}

	var filterDisplayed *bool
	if displayed := r.Query("displayed"); displayed == "true" {
		filterDisplayed = aws.Bool(true)
	} else if displayed == "false" {
		filterDisplayed = aws.Bool(false)
	}

	result := []*models.BodyPart{}
	for rows.Next() {
		item := &models.BodyPart{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Displayed,
			&item.Image,
			&item.Order,
			&item.Parent,
		); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to fetch the body parts",
			})
			return
		}

		// There won't be that many of them, might as well filter it here
		// than generate a query
		if filterDisplayed != nil &&
			(*filterDisplayed && !item.Displayed ||
				!*filterDisplayed && item.Displayed) {
			continue
		}

		result = append(result, item)
	}

	r.JSON(http.StatusOK, result)
}

// BodyPartsCreate inserts a new body part
func (a *API) BodyPartsCreate(r *gin.Context) {
	var input struct {
		Name      string             `json:"name"`
		Displayed bool               `json:"displayed"`
		Image     string             `json:"type"`
		Order     int                `json:"order"`
		Parent    types.ExtendedUUID `json:"parent"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input: " + err.Error(),
		})
		return
	}

	item := &models.BodyPart{
		Name:      input.Name,
		Displayed: input.Displayed,
		Image:     input.Image,
		Order:     input.Order,
		Parent:    input.Parent,
	}
	if err := models.CreateBodyPart(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to insert the body part",
		})
		return
	}

	r.JSON(http.StatusCreated, item)
}

// BodyPartsUpdate updates an existing body part by its ID
func (a *API) BodyPartsUpdate(r *gin.Context) {
	var input struct {
		Name      string             `json:"name"`
		Displayed bool               `json:"displayed"`
		Image     string             `json:"type"`
		Order     int                `json:"order"`
		Parent    types.ExtendedUUID `json:"parent"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input: " + err.Error(),
		})
		return
	}

	idString := r.Param("body-part")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	item := &models.BodyPart{
		ID:        id,
		Name:      input.Name,
		Displayed: input.Displayed,
		Image:     input.Image,
		Order:     input.Order,
		Parent:    input.Parent,
	}

	if err := models.UpdateBodyPart(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Body part not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}

// BodyPartsDelete deletes a body part.
func (a *API) BodyPartsDelete(r *gin.Context) {
	idString := r.Param("body-part")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	item, err := models.DeleteBodyPart(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Body part not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}
