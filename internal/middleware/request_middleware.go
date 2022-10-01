package middleware

import (
	"go-mongo-auth/internal/jwt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Headers struct {
	Authorization string `json:"Authorization" binding:"required,startswith=Bearer "`
	Source        string `json:"source" binding:"required,eq=test"`
}

type InvalidRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var exclude = []string{"/login", "/register"}

func RequestValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, e := range exclude {
			if strings.HasSuffix(ctx.Request.RequestURI, e) {
				return
			}
		}

		log.Println("Validating:", ctx.Request.RequestURI)

		var headers Headers
		if err := ctx.ShouldBindHeader(&headers); err != nil {
			log.Println("Error while binding headers.", err)
			ctx.JSON(http.StatusBadRequest, InvalidRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		// Extract token and validate
		token := strings.Split(headers.Authorization, " ")[1]

		entity, err := jwt.ValidateToken(token)
		if err != nil || entity == nil {
			log.Println("Error validating token.", err)
			ctx.JSON(http.StatusUnauthorized, InvalidRequest{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
		}
	}
}