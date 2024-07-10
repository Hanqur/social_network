package middleware

import (
	"net/http"
	"social/internal/service"
	rest "social/internal/transport/REST"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	authSvc *service.AuthSvc
}

const authHeader = "Authorization"

func New(authSvc *service.AuthSvc) *Middleware {
	return &Middleware{authSvc: authSvc}
}

func (m *Middleware) Identify(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		rest.NewErrResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	splitedHeader := strings.Split(header, " ")
	if len(splitedHeader) != 2 {
		rest.NewErrResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userID, err := m.authSvc.ParseToken(splitedHeader[1])
	if err != nil {
		rest.NewErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set("userCtx", userID)
	c.Next()
}
