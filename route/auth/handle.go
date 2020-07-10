package auth

import (
	"crm-service/config"
	"crm-service/model"
	"crm-service/util"
	"github.com/husol/libs"
	mod "github.com/husol/middleware/model"
	"github.com/husol/middleware/token"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

func (r *Route) Login(c echo.Context) error {
	husAjax := libs.HusAjax{}

	var (
		ctx = &util.CustomEchoContext{Context: c}
		req = model.LoginRequest{}
	)

	err := c.Bind(&req)
	if err != nil {
		husAjax.SetError(501, "Invalid Email or Password.")

		return c.JSON(http.StatusOK, husAjax.OutData(err.Error()))
	}

	myUser, err := r.route.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		husAjax.SetError(401, "Invalid Email or Password.")

		return c.JSON(http.StatusOK, husAjax.OutData(err.Error()))
	}

	check := libs.VerifyPassword(req.Password, myUser.Password)
	if !check {
		husAjax.SetError(401, "Invalid Email or Password.")

		return c.JSON(http.StatusOK, husAjax.OutData("invalid_password"))
	}

	var userInfo model.UserInfo
	if err := copier.Copy(&userInfo, myUser); err != nil {
		husAjax.SetError(401, "Error copying user info.")

		return c.JSON(http.StatusOK, husAjax.OutData(err.Error()))
	}

	var (
		tokenObj   = token.Token{}
		cfg        = config.GetConfig()
		expireTime = time.Now().Add(12 * time.Hour)
		userClaim  = mod.UserClaims{UserID: int(myUser.ID)}
	)

	accessToken, err := tokenObj.Encode(userClaim, cfg.SecretKey, expireTime)
	if err != nil {
		husAjax.SetError(401, "Error to get Access Token.")

		return c.JSON(http.StatusOK, husAjax.OutData(err.Error()))
	}

	resp := model.ResponseLogin{
		Token: accessToken,
		Info:  userInfo,
	}

	return c.JSON(http.StatusOK, husAjax.OutData(resp))
}
