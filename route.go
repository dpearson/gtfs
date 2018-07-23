package gtfs

import (
	"archive/zip"
	"fmt"
	"strconv"
)

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

type RouteType int

const (
	RouteTypeLightRail RouteType = iota
	RouteTypeSubway
	RouteTypeRail
	RouteTypeBus
	RouteTypeFerry
	RouteTypeCableCar
	RouteTypeGondola
	RouteTypeFunicular
)

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
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, routeFields)
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
				return fmt.Errorf("Invalid route sort order: %v", err)
			}
		}

		routeType, err := parseRouteType(row["route_type"])
		if err != nil {
			return err
		}

		r := &Route{
			ID:          row["route_id"],
			Agency:      g.agencyByID(row["agency_id"]),
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

func parseRouteType(val string) (RouteType, error) {
	switch val {
	case "0":
		return RouteTypeLightRail, nil
	case "1":
		return RouteTypeSubway, nil
	case "2":
		return RouteTypeRail, nil
	case "3":
		return RouteTypeBus, nil
	case "4":
		return RouteTypeFerry, nil
	case "5":
		return RouteTypeCableCar, nil
	case "6":
		return RouteTypeGondola, nil
	case "7":
		return RouteTypeFunicular, nil
	default:
		return 0, fmt.Errorf("Invalid route type: %s", val)
	}
}