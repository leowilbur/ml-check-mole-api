package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/auth"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/models"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token in the headers and checks its
// scope and permissions.
func (a *API) AuthMiddleware(permissions ...string) func(r *gin.Context) {
	return func(r *gin.Context) {
		header := r.GetHeader("Authorization")
		//Get token from parameter
		if header == "" && len(r.Query("token")) > 0 && r.Request.Method == "GET" {
			header = "Bearer " + r.Query("token")
		}
		if header == "" {
			r.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  100,
				"error": "The Authorization header is missing.",
			})
			return
		}

		headerParts := strings.SplitN(header, " ", 2)
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			r.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  100,
				"error": "The Authorization header is malformed.",
			})
			return
		}

		token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
			// Cognito users RS256 for JWTs
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// kid is the RSA key id
			keyID, ok := token.Header["kid"]
			if !ok {
				return nil, errors.New("Key ID missing from the JWT token")
			}

			// Cast it as string
			castedKeyID, ok := keyID.(string)
			if !ok {
				return nil, errors.New("Key ID has an invalid type")
			}

			// Match the keys in our JWK store
			matchedKeys := a.JWK.Key(castedKeyID)
			if len(matchedKeys) == 0 {
				return nil, errors.New("Unknown key ID")
			}

			return matchedKeys[0].Key, nil
		})
		if err != nil {
			r.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  102,
				"error": err.Error(),
			})
			return
		}

		if !token.Valid {
			r.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  102,
				"error": "The token is not valid",
			})
			return
		}

		// Check the permissions
		if len(permissions) > 0 {
			var groups []string

			claims := token.Claims.(jwt.MapClaims)
			interfaceSlice, ok := claims["cognito:groups"]
			if ok {
				castedSlice, ok := interfaceSlice.([]interface{})
				if ok {
					groups = []string{}
					for _, el := range castedSlice {
						groups = append(groups, el.(string))
					}
				}
			}

			if !auth.LookupPermissions(groups, permissions...) {
				r.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  105,
					"error": "Insufficient permissions",
				})
				return
			}
		}

		r.Set("token", token)

		userIDString := token.Claims.(jwt.MapClaims)["sub"].(string)
		userID, err := types.StringToUUID(userIDString)
		if err != nil {
			r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		account, err := models.GetAccount(r, a.DB, userID)
		if err != nil || account.UpdatedAt.Time.Add(30*time.Minute).Before(time.Now()) {
			userInfo, err := auth.GetUserInfo(r, a.Config.AuthURL, userIDString)
			//userInfo, err := auth.GetUserInfo(r, a.Config.AuthURL, headerParts[1])
			if err != nil {
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			parsedID, err := types.StringToUUID(userInfo.Subject)
			if err != nil {
				panic(err)
			}

			if err := models.UpsertAccount(r, a.DB, &models.Account{
				ID:        parsedID,
				Name:      userInfo.Name,
				Email:     userInfo.Email,
				Phone:     userInfo.PhoneNumber,
				Gender:    userInfo.Gender,
				BirthDate: userInfo.BirthDate,
			}); err != nil {
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    201,
					"message": err.Error(),
				})
				return
			}

			account, err = models.GetAccount(r, a.DB, userID)
			if err != nil {
				r.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    202,
					"message": err.Error(),
				})
				return
			}
		}

		r.Set("account", account)

		r.Next()
	}
}
