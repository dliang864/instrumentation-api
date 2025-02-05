package handlers

import (
	"net/http"
	"time"

	"github.com/USACE/instrumentation-api/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// ListInstrumentNotes returns instrument notes
func ListInstrumentNotes(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		notes, err := models.ListInstrumentNotes(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, notes)
	}
}

// ListInstrumentInstrumentNotes returns instrument notes for a single instrument
func ListInstrumentInstrumentNotes(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		iID, err := uuid.Parse(c.Param("instrument_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}
		notes, err := models.ListInstrumentInstrumentNotes(db, &iID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, notes)
	}
}

// GetInstrumentNote returns a single instrument note
func GetInstrumentNote(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		nID, err := uuid.Parse(c.Param("note_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Malformed ID")
		}
		note, err := models.GetInstrumentNote(db, &nID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, note)
	}
}

// CreateInstrumentNote creates instrument notes
func CreateInstrumentNote(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		nc := models.InstrumentNoteCollection{}
		if err := c.Bind(&nc); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// profile and timestamp
		p := c.Get("profile").(*models.Profile)

		t := time.Now()
		for idx := range nc.Items {
			nc.Items[idx].Creator = p.ID
			nc.Items[idx].CreateDate = t
		}
		nn, err := models.CreateInstrumentNote(db, nc.Items)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, nn)
	}
}

// UpdateInstrumentNote updates an instrument note
func UpdateInstrumentNote(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		noteID, err := uuid.Parse(c.Param("note_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		n := models.InstrumentNote{ID: noteID}
		if err := c.Bind(&n); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// check :id in url params matches id in request body
		if noteID != n.ID {
			return c.String(
				http.StatusBadRequest,
				"url note_id does not match object id in body",
			)
		}
		// profile and timestamp
		p := c.Get("profile").(*models.Profile)
		t := time.Now()
		n.Updater, n.UpdateDate = &p.ID, &t

		// update
		nUpdated, err := models.UpdateInstrumentNote(db, &n)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// return updated instrument note
		return c.JSON(http.StatusOK, nUpdated)
	}
}

// DeleteInstrumentNote deletes an instrument note
func DeleteInstrumentNote(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		noteID, err := uuid.Parse(c.Param("note_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if err := models.DeleteInstrumentNote(db, &noteID); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}
}
