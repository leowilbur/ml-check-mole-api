package rest

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

// QuestionsList lists all questins
func (a *API) QuestionsList(r *gin.Context) {
	rows, err := a.DB.QueryEx(
		r,
		`SELECT
			id, name, type, answers,
			displayed, "order",
			created_at, updated_at
		FROM questions ORDER BY "order" ASC`,
		nil,
	)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"message": err.Error(),
		})
		return
	}

	var filterDisplayed *bool
	if displayed := r.Query("displayed"); displayed == "true" {
		filterDisplayed = aws.Bool(true)
	} else if displayed == "false" {
		filterDisplayed = aws.Bool(false)
	}

	result := []*models.Question{}
	for rows.Next() {
		item := &models.Question{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Type,
			&item.Answers,
			&item.Displayed,
			&item.Order,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    201,
				"message": err.Error(),
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

// QuestionsCreate creates a new question
func (a *API) QuestionsCreate(r *gin.Context) {
	var input struct {
		Name      string              `json:"name"`
		Type      string              `json:"type"`
		Answers   types.ExtendedJSONB `json:"answers"`
		Displayed bool                `json:"displayed"`
		Order     int64               `json:"order"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    300,
			"message": err.Error(),
		})
		return
	}

	item := &models.Question{
		Name:      input.Name,
		Type:      input.Type,
		Answers:   input.Answers,
		Displayed: input.Displayed,
		Order:     input.Order,
	}
	if err := models.CreateQuestion(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to insert the question",
		})
		return
	}

	r.JSON(http.StatusCreated, item)
}

// QuestionsUpdate updates a single question
func (a *API) QuestionsUpdate(r *gin.Context) {
	var input struct {
		Name      string              `json:"name"`
		Type      string              `json:"type"`
		Answers   types.ExtendedJSONB `json:"answers"`
		Displayed bool                `json:"displayed"`
		Order     int64               `json:"order"`
	}
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input: " + err.Error(),
		})
		return
	}

	idString := r.Param("question")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	item := &models.Question{
		ID:        id,
		Name:      input.Name,
		Type:      input.Type,
		Answers:   input.Answers,
		Displayed: input.Displayed,
		Order:     input.Order,
	}
	if err := models.UpdateQuestion(r, a.DB, item); err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Question not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}

// QuestionsDelete deletes a single question
func (a *API) QuestionsDelete(r *gin.Context) {
	idString := r.Param("question")
	id, err := types.StringToUUID(idString)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID supplied: " + err.Error(),
		})
		return
	}

	item, err := models.DeleteQuestion(r, a.DB, id)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Question not found",
		})
		return
	}

	r.JSON(http.StatusOK, item)
}
