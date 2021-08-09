package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leowilbur/ml-check-mole-api/pkg/ultis"
)

// UploadFileToAWS ...
func (a *API) UploadFileToAWS(r *gin.Context) {
	accountID := r.PostForm("account_id")
	fileName := r.PostForm("file_name")
	imageBase64 := r.PostForm("file_base64")

	awsFileContent := ultis.AWSFileContent{
		AccountID:  accountID,
		FileName:   fileName,
		FileEncode: imageBase64,
	}
	err := awsFileContent.UploadFileToAWS()

	if err != nil {
		r.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	r.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"url":    fmt.Sprintf("/media/images/%s/%s", accountID, fileName),
	})
}

// DownloadFileFromAWS ...
func (a *API) DownloadFileFromAWS(r *gin.Context) {
	accountID := r.Param("account_id")
	fileName := r.Param("file_name")

	awsFileContent := ultis.AWSFileContent{
		AccountID: accountID,
		FileName:  fileName,
	}

	url, _ := awsFileContent.DownloadFileFromAWS()

	r.Redirect(http.StatusSeeOther, url)
}
