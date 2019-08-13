package gtfs

import (
	"fmt"
	"io"
	"strconv"
)

// DefaultRouteColor is the default color for a route with no specified color.
const DefaultRouteColor = "FFFFFF"

// DefaultRouteTextColor is the default text color for a route with no specified
// text color.
const DefaultRouteTextColor = "000000"

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

func (g *GTFS) processRoutes(r io.Reader) error {
	res, err := readCSVWithHeadings(r, routeFields, g.strictMode)
	if err != nil {
		return err
	}

	g.routesByID = map[string]*Route{}

	for _, row := range res {
		sortOrder, err := parseRouteSortOrder(row["route_sort_order"])
		if err != nil {
			return fmt.Errorf("invalid route sort order: %v", err)
		}

		routeType, err := parseRouteType(row["route_type"])
		if err != nil {
			return err
		}

		agency, err := g.agencyByIDOrDefault(row["agency_id"])
		if err != nil {
			return err
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
			r.Color = DefaultRouteColor
		}

		if r.TextColor == "" {
			r.TextColor = DefaultRouteTextColor
		}

		g.Routes = append(g.Routes, r)
		g.routesByID[r.ID] = r
	}

	return nil
}

func (g *GTFS) routeByID(id string) *Route {
	return g.routesByID[id]
}

func parseRouteSortOrder(val string) (uint64, error) {
	if val == "" {
		return 0, nil
	}

	return strconv.ParseUint(val, 10, 64)
}
