package gtfs

import (
	"archive/zip"
	"fmt"
	"strconv"
)

// A Route is a single route.
//
// Fields correspond directly to columns in routes.txt.
type Route struct {
	ID          string
	Agency      *Agency
	ShortName   string
	LongName    string
	Description string
	Type        RouteType
	URL         string
	Color       string
	TextColor   string
	SortOrder   uint64
}

var routeFields = map[string]bool{
	"route_id":         true,
	"agency_id":        false,
	"route_short_name": true,
	"route_long_name":  true,
	"route_desc":       false,
	"route_type":       true,
	"route_url":        false,
	"route_color":      false,
	"route_text_color": false,
	"route_sort_order": false,
}

func (g *GTFS) processRoutes(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() // nolint: errcheck

	res, err := readCSVWithHeadings(rc, routeFields, g.strictMode)
	if err != nil {
		return err
	}

	g.routesByID = map[string]*Route{}

	for _, row := range res {
		sortOrder := uint64(0)
		sortOrderStr := row["route_sort_order"]
		if sortOrderStr != "" {
			sortOrder, err = strconv.ParseUint(row["route_sort_order"], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid route sort order: %v", err)
			}
		}

		routeType, err := parseRouteType(row["route_type"])
		if err != nil {
			return err
		}

		var agency *Agency
		agencyID := row["agency_id"]
		if agencyID != "" {
			agency = g.agencyByID(row["agency_id"])
		} else if len(g.Agencies) != 1 {
			return fmt.Errorf("no agency_id specified, but there are %d agencies", len(g.Agencies))
		} else {
			agency = g.Agencies[0]
		}

		r := &Route{
			ID:          row["route_id"],
			Agency:      agency,
			ShortName:   row["route_short_name"],
			LongName:    row["route_long_name"],
			Description: row["route_desc"],
			Type:        routeType,
			URL:         row["route_url"],
			Color:       row["route_color"],
			TextColor:   row["route_text_color"],
			SortOrder:   sortOrder,
		}

		if r.Color == "" {
			r.Color = "FFFFFF"
		}

		if r.TextColor == "" {
			r.TextColor = "000000"
		}

		g.Routes = append(g.Routes, r)
		g.routesByID[r.ID] = r
	}

	return nil
}

func (g *GTFS) routeByID(id string) *Route {
	return g.routesByID[id]
}
