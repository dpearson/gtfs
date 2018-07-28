package gtfs

import (
	"archive/zip"
	"fmt"
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
}

func (g *GTFS) processStops(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, stopFields)
	if err != nil {
		return err
	}

	g.stopsByID = map[string]*Stop{}

	for _, row := range res {
		lat, err := strconv.ParseFloat(row["stop_lat"], 64)
		if err != nil {
			return fmt.Errorf("Invalid latitude: %v", err)
		}

		lon, err := strconv.ParseFloat(row["stop_lon"], 64)
		if err != nil {
			return fmt.Errorf("Invalid longitude: %v", err)
		}

		var locType LocationType
		switch row["location_type"] {
		case "0", "":
			locType = LocationTypeStop
		case "1":
			locType = LocationTypeStation
		case "2":
			locType = LocationTypeStationEntrance
		default:
			return fmt.Errorf("Invalid location type: %s", row["location_type"])
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
			return fmt.Errorf("Invalid location type with parent station: %d", s.LocationType)
		}

		parent, ok := g.stopsByID[s.parentStationID]
		if !ok {
			return fmt.Errorf("Invalid parent stop ID: %s", s.parentStationID)
		}

		s.ParentStation = parent
	}

	return nil
}

func (g *GTFS) stopByID(id string) *Stop {
	return g.stopsByID[id]
}
