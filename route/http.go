package route

import (
	"crm-service/config"
	"crm-service/repository"
	au "crm-service/route/auth"
	"crm-service/route/user"
	"fmt"
	"github.com/husol/libs"
	mid "github.com/husol/middleware/middleware"
	mod "github.com/husol/middleware/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func NewHTTPHandler(repo *repository.Repository) *echo.Echo {
	e := echo.New()
	husAjax := libs.HusAjax{}
	loggerCfg := middleware.DefaultLoggerConfig

	loggerCfg.Skipper = func(c echo.Context) bool {
		return c.Request().URL.RequestURI() == "/check-health"
	}

	e.Use(middleware.LoggerWithConfig(loggerCfg))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch,
			http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if c.Request().URL.RequestURI() != "/check-health" {
			request := fmt.Sprintf("%s", reqBody)

			if len(request) > 0 {
				log.Printf("%s", request)
			}
			log.Printf("%s", resBody)
		}
	}))

	e.GET("/check-health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, husAjax.OutData("OK"))
	})

	cfg := config.GetConfig()

	// Authenticate
	auth := e.Group("/auth")
	au.Init(auth, repo)

	// API
	api := e.Group("/v1")
	api.Use(mid.SetClaim(cfg.SecretKey, []mod.AllowedRoute{}))
	user.Init(api, repo)

	return e
}
