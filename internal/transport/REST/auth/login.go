package auth

import (
	"net/http"
	"social/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type LoginRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func HandleLogin(authSvc *service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(err)

			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}

		token, err := authSvc.Login(c.Request.Context(), req.ID, req.Password)
		if err != nil {
			log.Error(err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}
