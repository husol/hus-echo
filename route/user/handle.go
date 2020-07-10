package user

import (
	"crm-service/util"
	"github.com/husol/libs"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (r *Route) GetUsers(c echo.Context) error {
	ctx := &util.CustomEchoContext{Context: c}
	husAjax := libs.HusAjax{}

	users, err := r.route.UserRepo.GetAll(ctx)
	if err != nil {
		husAjax.SetError(501, err.Error())

		return c.JSON(http.StatusOK, husAjax.OutData(false))
	}

	return c.JSON(http.StatusOK, husAjax.OutData(users))
}
