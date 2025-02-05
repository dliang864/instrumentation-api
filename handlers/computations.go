package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/USACE/instrumentation-api/models"
)

// ComputedTimeseries returns computed timeseries for a given instrument
// This is an endpoint for debugging at this time
func ComputedTimeseries(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Instrument ID
		instrumentID, err := uuid.Parse("a7540f69-c41e-43b3-b655-6e44097edb7e")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		instrumentIDs := make([]uuid.UUID, 1)
		instrumentIDs[0] = instrumentID
		// Time Window
		timeWindow := models.TimeWindow{
			After:  time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			Before: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
		}
		// Interval - Hard Code at 1 Hour
		interval := time.Hour

		tt, err := models.ComputedTimeseries(db, instrumentIDs, &timeWindow, &interval)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, &tt)
	}
}
