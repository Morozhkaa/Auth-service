package http

import (
	"auth/internal/domain/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
)

func (a *Adapter) ErrorHandler(ctx *gin.Context, err error) {

	l := zapctx.Logger(ctx)
	l.Sugar().Errorf("request failed: %s", err.Error())

	switch {
	case errors.Is(err, models.ErrForbidden), errors.Is(err, models.ErrNotFound),
		errors.Is(err, models.ErrGenerateToken), errors.Is(err, models.ErrTokenExpired):
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	case errors.Is(err, models.ErrBadRequest):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
