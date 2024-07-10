package auth

import (
	"net/http"
	"social/internal/entity"
	"social/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const layout = "02.06.2006"

type RegisterRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	BirthDate  string `json:"birthdate"`
	Sex        string `json:"sex"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

func HandleRegister(authSvc *service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(err)

			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}

		birthDate, err := time.Parse(layout, req.BirthDate)
		if err != nil {
			log.Error(err)

			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}

		opts := &entity.CreateUserOpts{
			FirstName:  req.FirstName,
			SecondName: req.SecondName,
			BirthDate:  birthDate,
			Sex:        entity.SexType(req.Sex),
			Biography:  req.Biography,
			City:       req.City,
			Password:   req.Password,
		}

		userID, err := authSvc.CreateUser(c.Request.Context(), opts)
		if err != nil {
			log.Error(err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, RegisterResponse{
			UserID: userID,
		})
	}
}
