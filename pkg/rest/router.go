package rest

import (
	"net/http"

	fcm "github.com/appleboy/go-fcm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/leowilbur/ml-check-mole-api/docs/api"
	"github.com/leowilbur/ml-check-mole-api/pkg/models"
)

// API is the REST API
type API struct {
	*gin.Engine
	DB     models.DB
	JWK    *jose.JSONWebKeySet
	Config *Config
	S3     *s3.S3
}

// Config contains the configuration parameters used by the API
type Config struct {
	AuthURL      string
	AWSSession   *session.Session
	FCMClient    *fcm.Client
	PhotosBucket string
}

// New creates a new API using the given dependencies
func New(
	db models.DB,
	jwk *jose.JSONWebKeySet,
	cfg *Config,
) (*API, error) {
	gin.SetMode(gin.ReleaseMode)

	r := &API{
		Engine: gin.New(),
		DB:     db,
		JWK:    jwk,
		Config: cfg,
	}

	if cfg.AWSSession != nil {
		r.S3 = s3.New(cfg.AWSSession)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsMiddleware := cors.AllowAll()
	r.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
	})

	// r.GET("/", r.Index)
	r.GET("/swagger.json", func(r *gin.Context) {
		r.Header("Content-Type", "application/json")
		r.String(http.StatusOK, api.JSON)
	})

	r.GET("/body-parts", r.BodyPartsList)
	r.POST(
		"/body-parts",
		r.AuthMiddleware("body-parts.write"), r.BodyPartsCreate,
	)
	r.PUT(
		"/body-parts/:body-part",
		r.AuthMiddleware("body-parts.write"), r.BodyPartsUpdate,
	)
	r.DELETE(
		"/body-parts/:body-part",
		r.AuthMiddleware("body-parts.write"), r.BodyPartsDelete,
	)

	r.GET("/questions", r.QuestionsList)
	r.POST(
		"/questions",
		r.AuthMiddleware("questions.write"), r.QuestionsCreate)
	r.PUT(
		"/questions/:question",
		r.AuthMiddleware("questions.write"), r.QuestionsUpdate,
	)
	r.DELETE(
		"/questions/:question",
		r.AuthMiddleware("questions.write"), r.QuestionsDelete,
	)

	r.GET("/users/me", r.AuthMiddleware(), r.MyUser)

	r.GET(
		"/users/me/lesions",
		r.AuthMiddleware("owned.lesions.read"), r.LesionsList,
	)
	r.POST(
		"/users/me/lesions",
		r.AuthMiddleware("owned.lesions.write"), r.LesionsCreate,
	)
	r.PUT(
		"/users/me/lesions/:lesion",
		r.AuthMiddleware("owned.lesions.write"), r.LesionsUpdate,
	)
	r.DELETE(
		"/users/me/lesions/:lesion",
		r.AuthMiddleware("owned.lesions.write"), r.LesionsDelete,
	)

	r.GET(
		"/users/me/lesions/:lesion/reports",
		r.AuthMiddleware("owned.lesions.read"), r.ReportsList,
	)
	r.POST(
		"/users/me/lesions/:lesion/reports",
		r.AuthMiddleware("owned.lesions.write"), r.ReportsCreate,
	)
	r.PUT(
		"/users/me/lesions/:lesion/reports/:report",
		r.AuthMiddleware("owned.lesions.write"), r.ReportsUpdate,
	)
	r.DELETE(
		"/users/me/lesions/:lesion/reports/:report",
		r.AuthMiddleware("owned.lesions.write"), r.ReportsDelete,
	)

	r.GET(
		"/users/me/requests",
		r.AuthMiddleware("owned.lesions.read"), r.RequestsList,
	)
	r.GET(
		"/users/me/requests/:request",
		r.AuthMiddleware("owned.lesions.write"), r.RequestsGet,
	)
	r.POST(
		"/users/me/requests",
		r.AuthMiddleware("owned.lesions.write", "owned.requests.create"), r.RequestsCreate,
	)
	r.PUT(
		"/users/me/requests/:request",
		r.AuthMiddleware("owned.lesions.write"), r.RequestsUpdate,
	)
	r.DELETE(
		"/users/me/requests/:request",
		r.AuthMiddleware("owned.lesions.write"), r.RequestsDelete,
	)

	r.GET("/requests", r.AuthMiddleware("requests.read"), r.DoctorRequestsList)
	r.GET("/requests/:request", r.AuthMiddleware("requests.read"), r.DoctorRequestsGet)
	r.GET("/reports", r.AuthMiddleware("reports.read"), r.DoctorReportsList)
	r.GET("/reports/:report", r.AuthMiddleware("reports.read"), r.DoctorReportsGet)
	r.GET("/lesions", r.AuthMiddleware("lesions.read"), r.DoctorLesionsList)
	r.GET("/lesions/:lesion", r.AuthMiddleware("lesions.read"), r.DoctorLesionsGet)

	r.GET("/media/images/:account_id/:file_name", r.AuthMiddleware(), r.DownloadFileFromAWS)
	r.POST("/media/images/upload", r.AuthMiddleware(), r.UploadFileToAWS)

	r.PUT("/requests/:request", r.AuthMiddleware("requests.respond"), r.DoctorRequestsRespond)

	return r, nil
}
