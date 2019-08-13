package gtfs

import (
	"fmt"
	"io"
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

func (g *GTFS) processAgencies(r io.Reader) error {
	res, err := readCSVWithHeadings(r, agencyFields, g.strictMode)
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

// agencyByIDOrDefault gets the agency specified by id, unless id is empty.
//
// If id is empty, it returns either the only agency contained within the file
// or, if there is more than one agency defined, an error.
//
// It is the caller's responsibility to ensure that processAgencies has been
// called before invoking this function.
func (g *GTFS) agencyByIDOrDefault(id string) (*Agency, error) {
	if id != "" {
		return g.agencyByID(id), nil
	}

	if len(g.Agencies) != 1 {
		return nil, fmt.Errorf("no agency_id specified, but there are %d agencies", len(g.Agencies))
	}

	return g.Agencies[0], nil
}
