package handlers

import (
	"net/http"

	"github.com/USACE/instrumentation-api/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// ListInstrumentStatus lists all Status for an instrument
func ListInstrumentStatus(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("instrument_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}

		ss, err := models.ListInstrumentStatus(db, &id)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, ss)
	}
}

// GetInstrumentStatus returns a single Status
func GetInstrumentStatus(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("status_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}

		s, err := models.GetInstrumentStatus(db, &id)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}

// CreateOrUpdateInstrumentStatus creates a Status for an instrument
func CreateOrUpdateInstrumentStatus(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		instrumentID, err := uuid.Parse(c.Param("instrument_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}

		var sc models.InstrumentStatusCollection
		if err := c.Bind(&sc); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// Assign Fresh UUID to each Status
		for idx := range sc.Items {
			id, err := uuid.NewRandom()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			sc.Items[idx].ID = id
		}

		if err := models.CreateOrUpdateInstrumentStatus(db, &instrumentID, sc.Items); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusCreated, make(map[string]interface{}))
	}
}

// DeleteInstrumentStatus deletes a Status for an instrument
func DeleteInstrumentStatus(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("status_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}
		if err := models.DeleteInstrumentStatus(db, &id); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}
}
