package handlers

import (
	"net/http"

	"github.com/USACE/instrumentation-api/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListAwareParameters(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		pp, err := models.ListAwareParameters(db)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, pp)
	}
}

func ListAwarePlatformParameterConfig(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc, err := models.ListAwarePlatformParameterConfig(db)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, cc)
	}
}
