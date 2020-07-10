package auth

import (
	"crm-service/repository"
	"github.com/labstack/echo/v4"
)

type Route struct {
	route *repository.Repository
}

func Init(group *echo.Group, controller *repository.Repository) {
	r := &Route{controller}

	group.POST("", r.Login)
}
