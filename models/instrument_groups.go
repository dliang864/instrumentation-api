package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// InstrumentGroup holds information for entity instrument_group
type InstrumentGroup struct {
	ID          uuid.UUID  `json:"id"`
	Deleted     bool       `json:"-"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProjectID   *uuid.UUID `json:"project_id" db:"project_id"`
	InstrumentCount int    `json:"instrument_count" db:"instrument_count"`
	TimeseriesCount int    `json:"timeseries_count" db:"timeseries_count"`
	AuditInfo
}

// InstrumentGroupCollection is a collection of Instrument items
type InstrumentGroupCollection struct {
	Items []InstrumentGroup
}

// Shorten returns an instrument collection with individual objects limited to ID and Struct fields
func (c InstrumentGroupCollection) Shorten() IDAndSlugCollection {

	ss := IDAndSlugCollection{Items: make([]IDAndSlug, 0)}
	for _, n := range c.Items {
		s := IDAndSlug{ID: n.ID, Slug: n.Slug}

		ss.Items = append(ss.Items, s)
	}
	return ss
}

// UnmarshalJSON implements UnmarshalJSON interface
// Allows unpacking object or array of objects into array of objects
func (c *InstrumentGroupCollection) UnmarshalJSON(b []byte) error {

	switch JSONType(b) {
	case "ARRAY":
		if err := json.Unmarshal(b, &c.Items); err != nil {
			return err
		}
	case "OBJECT":
		var g InstrumentGroup
		if err := json.Unmarshal(b, &g); err != nil {
			return err
		}
		c.Items = []InstrumentGroup{g}
	default:
		c.Items = make([]InstrumentGroup, 0)
	}
	return nil
}

// ListInstrumentGroupSlugs lists used instrument group slugs in the database
func ListInstrumentGroupSlugs(db *sqlx.DB) ([]string, error) {

	ss := make([]string, 0)
	if err := db.Select(&ss, "SELECT slug FROM instrument_group"); err != nil {
		return make([]string, 0), err
	}
	return ss, nil
}

// ListInstrumentGroups returns a list of instrument groups
func ListInstrumentGroups(db *sqlx.DB) ([]InstrumentGroup, error) {

	gg := make([]InstrumentGroup, 0)
	if err := db.Select(
		&gg, listInstrumentGroupsSQL+" WHERE NOT deleted",
	); err != nil {
		return make([]InstrumentGroup, 0), err
	}
	return gg, nil
}

// GetInstrumentGroup returns a single instrument group
func GetInstrumentGroup(db *sqlx.DB, ID uuid.UUID) (*InstrumentGroup, error) {

	var g InstrumentGroup
	if err := db.QueryRowx(
		listInstrumentGroupsSQL+" WHERE id = $1",
		ID,
	).StructScan(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// CreateInstrumentGroup creates many instruments from an array of instruments
func CreateInstrumentGroup(db *sqlx.DB, groups []InstrumentGroup) ([]InstrumentGroup, error) {

	txn, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	stmt, err := txn.Preparex(
		`INSERT INTO instrument_group (slug, name, description, creator, create_date, project_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, slug, name, description, creator, create_date, updater, update_date, project_id`,
	)
	if err != nil {
		return make([]InstrumentGroup, 0), err
	}

	gg := make([]InstrumentGroup, len(groups))
	for idx, g := range groups {
		if err := stmt.Get(&gg[idx], g.Slug, g.Name, g.Description, g.Creator, g.CreateDate, g.ProjectID); err != nil {
			return make([]InstrumentGroup, 0), err
		}
	}

	if err := stmt.Close(); err != nil {
		return make([]InstrumentGroup, 0), err
	}

	if err := txn.Commit(); err != nil {
		return make([]InstrumentGroup, 0), err
	}

	return gg, nil
}

// UpdateInstrumentGroup updates an instrument group
func UpdateInstrumentGroup(db *sqlx.DB, g *InstrumentGroup) (*InstrumentGroup, error) {

	var gUpdated InstrumentGroup
	if err := db.QueryRowx(
		`UPDATE instrument_group
		 SET    name = $2,
			    deleted = $3,
			    description = $4,
			    updater = $5,
				update_date = $6,
				project_id = $7
		 WHERE id = $1
		 RETURNING *
		`, g.ID, g.Name, g.Deleted, g.Description, g.Updater, g.UpdateDate, g.ProjectID,
	).StructScan(&gUpdated); err != nil {
		return nil, err
	}
	return &gUpdated, nil
}

// DeleteFlagInstrumentGroup sets the deleted field to true
func DeleteFlagInstrumentGroup(db *sqlx.DB, ID uuid.UUID) error {
	if _, err := db.Exec(`UPDATE instrument_group SET deleted = true WHERE id = $1`, ID); err != nil {
		return err
	}
	return nil
}

// ListInstrumentGroupInstruments returns a list of instrument group instruments for a given instrument
func ListInstrumentGroupInstruments(db *sqlx.DB, ID uuid.UUID) ([]Instrument, error) {

	sql := fmt.Sprintf(
		`SELECT B.*
         FROM   instrument_group_instruments A
		 INNER JOIN (%s) B ON A.instrument_id = B.id
		 WHERE A.instrument_group_id = $1 and B.deleted = false`,
		listInstrumentsSQL,
	)

	rows, err := db.Queryx(sql, ID)
	if err != nil {
		return make([]Instrument, 0), err
	}
	return InstrumentsFactory(rows)
}

// CreateInstrumentGroupInstruments adds an instrument to an instrument group
func CreateInstrumentGroupInstruments(db *sqlx.DB, instrumentGroupID uuid.UUID, instrumentID uuid.UUID) error {

	if _, err := db.Exec(
		`INSERT INTO instrument_group_instruments (instrument_group_id, instrument_id) VALUES ($1, $2)`,
		instrumentGroupID,
		instrumentID,
	); err != nil {
		return err
	}

	return nil
}

// DeleteInstrumentGroupInstruments adds an instrument to an instrument group
func DeleteInstrumentGroupInstruments(db *sqlx.DB, instrumentGroupID uuid.UUID, instrumentID uuid.UUID) error {

	if _, err := db.Exec(
		`DELETE FROM instrument_group_instruments WHERE instrument_group_id = $1 and instrument_id = $2`,
		instrumentGroupID, instrumentID,
	); err != nil {
		return err
	}

	return nil
}

var listInstrumentGroupsSQL = `SELECT id,
				                      slug,
				                      name,
				                      description,
				                      creator,
				                      create_date,
				                      updater,
				                      update_date,
				                      project_id,
									  instrument_count,
									  timeseries_count 
	                            FROM  v_instrument_group`
