package gtfs

import (
	"archive/zip"
	"fmt"
)

// An Agency is a single agency from a GTFS file.
//
// Fields correspond directly to columns in agency.txt.
type Agency struct {
	ID       string
	Name     string
	URL      string
	Timezone string
	Lang     string
	Phone    string
	FareURL  string
	Email    string
}

var agencyFields = map[string]bool{
	"agency_id":       false,
	"agency_name":     true,
	"agency_url":      true,
	"agency_timezone": true,
	"agency_lang":     false,
	"agency_phone":    false,
	"agency_fare_url": false,
	"agency_email":    false,
}

func (g *GTFS) processAgencies(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() // nolint: errcheck

	res, err := readCSVWithHeadings(rc, agencyFields)
	if err != nil {
		return err
	}

	g.agenciesByID = map[string]*Agency{}

	hasAgencyWithoutID := false
	for _, row := range res {
		hasAgencyWithoutID = hasAgencyWithoutID || row["agency_id"] == ""

		a := &Agency{
			ID:       row["agency_id"],
			Name:     row["agency_name"],
			URL:      row["agency_url"],
			Timezone: row["agency_timezone"],
			Lang:     row["agency_lang"],
			Phone:    row["agency_phone"],
			FareURL:  row["agency_fare_url"],
			Email:    row["agency_email"],
		}

		g.Agencies = append(g.Agencies, a)
		g.agenciesByID[a.ID] = a
	}

	if hasAgencyWithoutID && len(g.Agencies) > 1 {
		return fmt.Errorf("agency IDs must be specified")
	}

	return nil
}

func (g *GTFS) agencyByID(id string) *Agency {
	return g.agenciesByID[id]
}
