package common

import (
	"log"
	"net/http"
	"time"

	"github.com/aunjaidev/aunjai-common/logger"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerBase struct {
	App     *echo.Echo
	AppName string
	Port    string
}

func init() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func CreateServer(appName string, port string) *ServerBase {
	app := echo.New()

	// Enable metrics middleware
	p := prometheus.NewPrometheus("go-"+appName, nil)
	p.Use(app)

	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	app.Logger = logger.GetEchoLogger()
	app.Use(logger.Hook())

	sb := &ServerBase{
		App:     app,
		AppName: appName,
		Port:    port,
	}

	return sb
}

func (sb *ServerBase) StartServer() {
	log.Panic(sb.App.Start(":" + sb.Port))
}

type HealthCheck struct {
	Topic          string
	URL            string
	MiddlewareFunc func() int
}

func (sb *ServerBase) HealthCheck(healths []HealthCheck) {
	response := make(map[string]bool)
	flag := 0
	sb.App.GET("/health", func(c echo.Context) error {
		for _, h := range healths {
			if h.MiddlewareFunc() == 200 {
				response[h.Topic] = true
				flag++
			} else {
				response[h.Topic] = false
			}
		}

		if flag == len(response) {
			return c.JSON(http.StatusOK, response)
		} else {
			return c.JSON(http.StatusBadRequest, response)
		}
	})
}
