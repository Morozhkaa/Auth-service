package http

import (
	"auth/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

const (
	accessHeader  = "access_token"
	refreshHeader = "refresh_token"
)

// @ID login
// @tags login
// @Summary Log in
// @Description Выполняет авторизацию пользователя
// @Security BasicAuth
// @Success 200 {string} string "Success"
// @Failure 400 {object} error "Required parameters not filled"
// @Failure 403 {object} error "Insufficient rights to perform the operation"
// @Router /login [post]
func (a *Adapter) login(ctx *gin.Context) {
	login, password, ok := ctx.Request.BasicAuth()
	zap.L().Sugar().Infof("http/login: login: %s, password: %s", login, password)

	if !ok || login == "" || password == "" {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	ctx.Request = ctx.Request.WithContext(zapctx.WithFields(ctx.Request.Context(),
		zap.String("login", login),
	))
	accesstoken, refreshToken, err := a.auth.Login(ctx, login, password)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}

	accessCookie := &http.Cookie{
		Name:  accessHeader,
		Value: accesstoken,
	}
	refreshCookie := &http.Cookie{
		Name:  refreshHeader,
		Value: refreshToken,
	}
	http.SetCookie(ctx.Writer, accessCookie)
	http.SetCookie(ctx.Writer, refreshCookie)
}

// @ID verify
// @tags verify
// @Summary Verify
// @Description Выполняет проверку токенов пользователя, которые получает из cookies после выполнения /login запроса
// @Success 200 {string} string "Success"
// @Failure 403 {object} error "Insufficient rights to perform the operation"
// @Router /verify [post]
func (a *Adapter) verify(ctx *gin.Context) {
	access, err := ctx.Cookie(accessHeader)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	refresh, err := ctx.Cookie(refreshHeader)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}

	resp, err := a.auth.Verify(ctx, access, refresh)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	cookie := &http.Cookie{
		Name:  accessHeader,
		Value: resp.AccessToken,
	}
	cookie2 := &http.Cookie{
		Name:  refreshHeader,
		Value: resp.RefreshToken,
	}
	http.SetCookie(ctx.Writer, cookie)
	http.SetCookie(ctx.Writer, cookie2)

	ctx.JSON(http.StatusOK, gin.H{
		"login": resp.Login,
		"email": resp.Email,
	})
}
