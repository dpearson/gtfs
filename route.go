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

// RouteType specifies the type of vehicles operating on a route.
type RouteType int

const (
	// RouteTypeNotSpecified indicates that no route type was specified.
	RouteTypeNotSpecified RouteType = iota

	// RouteTypeLightRail indicates that the route is a light rail route.
	RouteTypeLightRail

	// RouteTypeSubway indicates that the route is a subway or metro route.
	RouteTypeSubway

	// RouteTypeRail indicates that the route is an intercity/long-distance rail
	// route.
	RouteTypeRail

	// RouteTypeBus indicates that the route is a bus route.
	RouteTypeBus

	// RouteTypeFerry indicates that the route is a ferry route.
	RouteTypeFerry

	// RouteTypeCableCar indicates that the route is a cable car route.
	RouteTypeCableCar

	// RouteTypeGondola indicates that the route is an aerial gondola route.
	RouteTypeGondola

	// RouteTypeFunicular indicates that the route is a funicular route.
	RouteTypeFunicular

	// EXTENDED ROUTE TYPES
	//
	// The following types are extended route types proposed by Google, although
	// they are not part of the GTFS standard.

	RouteTypeExtendedRailwayService
	RouteTypeExtendedHighSpeedRail
	RouteTypeExtendedLongDistanceRail
	RouteTypeExtendedInterRegionalRail
	RouteTypeExtendedCarTransportRail
	RouteTypeExtendedSleeperRail
	RouteTypeExtendedRegionalRail
	RouteTypeExtendedTouristRail
	RouteTypeExtendedRailShuttle
	RouteTypeExtendedSuburbanRail
	RouteTypeExtendedReplacementRail
	RouteTypeExtendedSpecialRail
	RouteTypeExtendedLorryTransportRail
	RouteTypeExtendedAllRail
	RouteTypeExtendedCrossCountryRail
	RouteTypeExtendedVehicleTransportRail
	RouteTypeExtendedRackAndPinionRail
	RouteTypeExtendedAdditionalRail
)

var routeTypes = map[string]RouteType{
	"0":   RouteTypeLightRail,
	"1":   RouteTypeSubway,
	"2":   RouteTypeRail,
	"3":   RouteTypeBus,
	"4":   RouteTypeFerry,
	"5":   RouteTypeCableCar,
	"6":   RouteTypeGondola,
	"7":   RouteTypeFunicular,
	"100": RouteTypeExtendedRailwayService,
	"101": RouteTypeExtendedHighSpeedRail,
	"102": RouteTypeExtendedLongDistanceRail,
	"103": RouteTypeExtendedInterRegionalRail,
	"104": RouteTypeExtendedCarTransportRail,
	"105": RouteTypeExtendedSleeperRail,
	"106": RouteTypeExtendedRegionalRail,
	"107": RouteTypeExtendedTouristRail,
	"108": RouteTypeExtendedRailShuttle,
	"109": RouteTypeExtendedSuburbanRail,
	"110": RouteTypeExtendedReplacementRail,
	"111": RouteTypeExtendedSpecialRail,
	"112": RouteTypeExtendedLorryTransportRail,
	"113": RouteTypeExtendedAllRail,
	"114": RouteTypeExtendedCrossCountryRail,
	"115": RouteTypeExtendedVehicleTransportRail,
	"116": RouteTypeExtendedRackAndPinionRail,
	"117": RouteTypeExtendedAdditionalRail,
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

func parseRouteType(val string) (RouteType, error) {
	routeType, ok := routeTypes[val]
	if !ok {
		return 0, fmt.Errorf("invalid route type: %s", val)
	}

	return routeType, nil
}
