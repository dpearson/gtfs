package gtfs

import (
	"fmt"
	"io"
	"strconv"
)

// A Stop is a single stop served by an agency referenced in a GTFS feed.
//
// Fields correspond directly to columns in stops.txt.
type Stop struct {
	ID                 string
	Code               string
	Name               string
	Description        string
	Latitude           float64
	Longitude          float64
	ZoneID             string
	URL                string
	LocationType       LocationType
	ParentStation      *Stop
	Timezone           string
	WheelchairBoarding string // TODO: parse me

	// Extensions:
	PlatformCode string
	VehicleType  RouteType

	parentStationID string
}

// LocationType specifies the specific type of a stop.
type LocationType int

const (
	// LocationTypeStop is a stop where passengers board or exit a vehicle.
	LocationTypeStop LocationType = iota

	// LocationTypeStation is a station containing at least one stop.
	LocationTypeStation

	// LocationTypeStationEntrance is the entrance to a station.
	LocationTypeStationEntrance
)

var stopFields = map[string]bool{
	"stop_id":             true,
	"stop_code":           false,
	"stop_name":           true,
	"stop_desc":           false,
	"stop_lat":            true,
	"stop_lon":            true,
	"zone_id":             false,
	"stop_url":            false,
	"location_type":       false,
	"parent_station":      false,
	"stop_timezone":       false,
	"wheelchair_boarding": false,

	// Extensions:
	"platform_code": false,
	"vehicle_type":  false,
}

func (g *GTFS) processStops(r io.Reader) error {
	res, err := readCSVWithHeadings(r, stopFields, g.strictMode)
	if err != nil {
		return err
	}

	g.stopsByID = map[string]*Stop{}

	for _, row := range res {
		lat, err := strconv.ParseFloat(row["stop_lat"], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude: %v", err)
		}

		lon, err := strconv.ParseFloat(row["stop_lon"], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude: %v", err)
		}

		locType, err := parseLocationType(row["location_type"])
		if err != nil {
			return err
		}

		var vehicleType RouteType
		if row["vehicle_type"] != "" {
			vehicleType, err = parseRouteType(row["vehicle_type"])
			if err != nil {
				return fmt.Errorf("invalid vehicle_type: %v", err)
			}
		}

		s := &Stop{
			ID:                 row["stop_id"],
			Code:               row["stop_code"],
			Name:               row["stop_name"],
			Description:        row["stop_desc"],
			Latitude:           lat,
			Longitude:          lon,
			ZoneID:             row["zone_id"],
			URL:                row["stop_url"],
			LocationType:       locType,
			Timezone:           row["stop_timezone"],
			WheelchairBoarding: row["wheelchair_boarding"],

			PlatformCode: row["platform_code"],
			VehicleType:  vehicleType,

			parentStationID: row["parent_station"],
		}

		g.Stops = append(g.Stops, s)
		g.stopsByID[s.ID] = s
	}

	for _, s := range g.Stops {
		if s.parentStationID == "" {
			continue
		}

		if s.LocationType != LocationTypeStop {
			return fmt.Errorf("invalid location type with parent station: %d", s.LocationType)
		}

		parent, ok := g.stopsByID[s.parentStationID]
		if !ok {
			if g.strictMode {
				return fmt.Errorf("invalid parent stop ID: %s for stop %s", s.parentStationID, s.ID)
			}

			continue
		}

		s.ParentStation = parent
	}

	return nil
}

func (g *GTFS) stopByID(id string) *Stop {
	return g.stopsByID[id]
}

func parseLocationType(val string) (LocationType, error) {
	switch val {
	case "0", "":
		return LocationTypeStop, nil
	case "1":
		return LocationTypeStation, nil
	case "2":
		return LocationTypeStationEntrance, nil
	default:
		return LocationTypeStop, fmt.Errorf("invalid location type: %s", val)
	}
}
