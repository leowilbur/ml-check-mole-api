package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MyUser returns details about the current user
func (a *API) MyUser(r *gin.Context) {
	r.JSON(http.StatusOK, r.MustGet("account"))
}
