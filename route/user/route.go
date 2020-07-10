package user

import (
	"crm-service/repository"
	"github.com/labstack/echo/v4"
)

type Route struct {
	route *repository.Repository
}

func Init(group *echo.Group, controller *repository.Repository) {
	r := &Route{controller}

	user := group.Group("/users")
	user.GET("", r.GetUsers)
}
