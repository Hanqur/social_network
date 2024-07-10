package user

import (
	"fmt"
	"net/http"
	"social/internal/entity"
	"social/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetUserResponse struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	BirthDate  string    `json:"birthdate"`
	Sex        string    `json:"sex"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
}

const userIDParam = "user_id"

func HandleGetUser(userSvc *service.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param(userIDParam)

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}

		user, err := userSvc.Get(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, prepareGetUserResponse(user))
	}
}

func prepareGetUserResponse(user *entity.User) *GetUserResponse {
	year, month, day := user.BirthDate.Date()

	dayStr := strconv.Itoa(day)
	monthStr := strconv.Itoa(int(month))

	if day < 10 {
		dayStr = "0" + dayStr
	}

	if int(month) < 10 {
		monthStr = "0" + monthStr
	}

	return &GetUserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		BirthDate:  fmt.Sprintf("%v-%v-%v", year, monthStr, dayStr),
		Sex:        user.Sex.String(),
		Biography:  user.Biography,
		City:       user.City,
	}
}
