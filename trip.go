package gtfs

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

// A Trip is a trip along a route with schedule information.
//
// Fields correspond directly to columns in trips.txt, stop_times.txt, and
// frequencies.txt.
type Trip struct { // nolint: maligned
	ID                   string
	Route                *Route
	Service              *Service
	Shape                *Shape
	Headsign             string
	ShortName            string
	DirectionID          string
	BlockID              string
	WheelchairAccessible WheelchairAccessible
	BikesAllowed         BikesAllowed

	AbsoluteTimes  bool
	StartTime      string
	EndTime        string
	HeadwaySeconds uint64
	ExactTimes     bool
	Stops          []*StopTime

	Exceptional bool
}

// StopTime provides details on a specific stop in a trip.
type StopTime struct {
	Stop                  *Stop
	ArrivalTime           string
	DepartureTime         string
	Sequence              uint64
	Headsign              string
	PickupType            PickupType
	DropoffType           DropoffType
	ShapeDistanceTraveled float64
	Timepoint             TimepointType
}

// WheelchairAccessible indicates whether a trip is accessible to passengers in
// wheelchairs.
type WheelchairAccessible int

const (
	// WheelchairAccessibleUnknown means that no wheelchair accessibility
	// information is available.
	WheelchairAccessibleUnknown WheelchairAccessible = iota

	// WheelchairAccessibleYes means that at least one passenger in a wheelchair
	// may be accommodated.
	WheelchairAccessibleYes

	// WheelchairAccessibleNo means that no passengers in wheelchairs may be
	// accommodated.
	WheelchairAccessibleNo
)

// BikesAllowed indicates whether bikes are allowed on a trip.
type BikesAllowed int

const (
	// BikesAllowedUnknown means that no information on whether bikes are
	// allowed is available.
	BikesAllowedUnknown BikesAllowed = iota

	// BikesAllowedYes means that at least one bike may be brought on this trip.
	BikesAllowedYes

	// BikesAllowedNo means that no bikes are allowed on this trip.
	BikesAllowedNo
)

// PickupType indicates the type of pickup available at a stop.
type PickupType int

const (
	// PickupTypeRegular indicates that regularly scheduled pickups are
	// available.
	PickupTypeRegular PickupType = iota

	// PickupTypeNone indicates that no pickups are available.
	PickupTypeNone

	// PickupTypePhoneAgency indicates that riders must phone the transit agency
	// to schedule pickups.
	PickupTypePhoneAgency

	// PickupTypeCoordinateWithDriver indicates that riders must coordinate with
	// the vehicle driver to schedule pickups.
	PickupTypeCoordinateWithDriver
)

// DropoffType indicates the type of dropoff available at a stop.
type DropoffType int

const (
	// DropoffTypeRegular indicates that regularly scheduled dropoffs are
	// available.
	DropoffTypeRegular DropoffType = iota

	// DropoffTypeNone indicates that no dropoffs are available.
	DropoffTypeNone

	// DropoffTypePhoneAgency indicates that riders must phone the transit
	// agency to schedule dropoffs.
	DropoffTypePhoneAgency

	// DropoffTypeCoordinateWithDriver indicates that riders must coordinate
	// with the vehicle driver to schedule dropoffs.
	DropoffTypeCoordinateWithDriver
)

// TimepointType specifies whether stop times are exact or approximate.
type TimepointType int

const (
	// TimepointTypeExact means that stop times are exact.
	TimepointTypeExact TimepointType = iota

	// TimepointTypeApproximate means that stop times are approximate.
	TimepointTypeApproximate
)

var tripFields = map[string]bool{
	"route_id":              true,
	"service_id":            true,
	"trip_id":               true,
	"trip_headsign":         false,
	"trip_short_name":       false,
	"direction_id":          false,
	"block_id":              false,
	"shape_id":              false,
	"wheelchair_accessible": false,
	"bikes_allowed":         false,

	// Extensions:
	"exceptional": false,
}

var stopTimeFields = map[string]bool{
	"trip_id":             true,
	"arrival_time":        true,
	"departure_time":      true,
	"stop_id":             true,
	"stop_sequence":       true,
	"stop_headsign":       false,
	"pickup_type":         false,
	"drop_off_type":       false,
	"shape_dist_traveled": false,
	"timepoint":           false,
}

var frequencyFields = map[string]bool{
	"trip_id":      true,
	"start_time":   true,
	"end_time":     true,
	"headway_secs": true,
	"exact_times":  false,
}

func (g *GTFS) processTrips(r io.Reader) error {
	res, err := readCSVWithHeadings(r, tripFields, g.strictMode)
	if err != nil {
		return err
	}

	g.tripsByID = map[string]*Trip{}

	for _, row := range res {
		wheelchairAccessible, err := parseWheelchairAccessible(row["wheelchair_accessible"])
		if err != nil {
			return err
		}

		bikesAllowed, err := parseBikesAllowed(row["bikes_allowed"])
		if err != nil {
			return err
		}

		exceptional, err := parseExceptional(row["exceptional"])
		if err != nil {
			return err
		}

		t := &Trip{
			ID:                   row["trip_id"],
			Route:                g.routeByID(row["route_id"]),
			Service:              g.serviceByID(row["service_id"]),
			Shape:                g.shapeByID(row["shape_id"]),
			Headsign:             row["trip_headsign"],
			ShortName:            row["trip_short_name"],
			DirectionID:          row["direction_id"],
			BlockID:              row["block_id"],
			WheelchairAccessible: wheelchairAccessible,
			BikesAllowed:         bikesAllowed,
			AbsoluteTimes:        true,

			Exceptional: exceptional,
		}

		g.Trips = append(g.Trips, t)
		g.tripsByID[t.ID] = t
	}

	return nil
}

func (g *GTFS) processStopTimes(r io.Reader) error {
	res, err := readCSVWithHeadings(r, stopTimeFields, g.strictMode)
	if err != nil {
		return err
	}

	stopsByTrip := map[string][]*StopTime{}
	for _, row := range res {
		seq, err := strconv.ParseUint(row["stop_sequence"], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid stop sequence: %v", err)
		}

		distStr := row["shape_dist_traveled"]
		dist := 0.0
		if distStr != "" {
			dist, err = strconv.ParseFloat(distStr, 64)
			if err != nil {
				return fmt.Errorf("invalid distance: %v", err)
			}
		}

		pickupType, err := parsePickupType(row["pickup_type"])
		if err != nil {
			return err
		}

		dropoffType, err := parseDropoffType(row["drop_off_type"])
		if err != nil {
			return err
		}

		timepointType, err := parseTimepointType(row["timepoint"])
		if err != nil {
			return err
		}

		s := &StopTime{
			Stop:                  g.stopByID(row["stop_id"]),
			ArrivalTime:           row["arrival_time"],
			DepartureTime:         row["departure_time"],
			Sequence:              seq,
			Headsign:              row["stop_headsign"],
			PickupType:            pickupType,
			DropoffType:           dropoffType,
			ShapeDistanceTraveled: dist,
			Timepoint:             timepointType,
		}

		stopsByTrip[row["trip_id"]] = append(stopsByTrip[row["trip_id"]], s)
	}

	for _, t := range g.Trips {
		stops, ok := stopsByTrip[t.ID]
		if !ok {
			continue
		}

		sort.Slice(stops, func(i, j int) bool {
			return stops[i].Sequence < stops[j].Sequence
		})

		t.Stops = stops
	}

	return nil
}

func (g *GTFS) processFrequencies(r io.Reader) error {
	res, err := readCSVWithHeadings(r, frequencyFields, g.strictMode)
	if err != nil {
		return err
	}

	for _, row := range res {
		t := g.tripByID(row["trip_id"])
		if t == nil {
			return fmt.Errorf("invalid trip id: %s", row["trip_id"])
		}

		headwaySecs, err := strconv.ParseUint(row["headway_secs"], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid headway seconds: %v", err)
		}

		var exactTimes bool
		switch row["exact_times"] {
		case "1":
			exactTimes = true
		case "0", "":
			exactTimes = false
		default:
			return fmt.Errorf("invalid exact times: %s", row["exact_times"])
		}

		t.AbsoluteTimes = false
		t.StartTime = row["start_time"]
		t.EndTime = row["end_time"]
		t.HeadwaySeconds = headwaySecs
		t.ExactTimes = exactTimes
	}

	return nil
}

func (g *GTFS) tripByID(id string) *Trip {
	return g.tripsByID[id]
}

func parseWheelchairAccessible(val string) (WheelchairAccessible, error) {
	switch val {
	case "0", "":
		return WheelchairAccessibleUnknown, nil
	case "1":
		return WheelchairAccessibleYes, nil
	case "2":
		return WheelchairAccessibleNo, nil
	default:
		return WheelchairAccessibleUnknown, fmt.Errorf("invalid wheelchair accessible value: %s", val)
	}
}

func parseBikesAllowed(val string) (BikesAllowed, error) {
	switch val {
	case "0", "":
		return BikesAllowedUnknown, nil
	case "1":
		return BikesAllowedYes, nil
	case "2":
		return BikesAllowedNo, nil
	default:
		return BikesAllowedUnknown, fmt.Errorf("invalid bikes allowed value: %s", val)
	}
}

func parsePickupType(val string) (PickupType, error) {
	switch val {
	case "0", "":
		return PickupTypeRegular, nil
	case "1":
		return PickupTypeNone, nil
	case "2":
		return PickupTypePhoneAgency, nil
	case "3":
		return PickupTypeCoordinateWithDriver, nil
	default:
		return PickupTypeRegular, fmt.Errorf("invalid pickup type: %s", val)
	}
}

func parseDropoffType(val string) (DropoffType, error) {
	switch val {
	case "0", "":
		return DropoffTypeRegular, nil
	case "1":
		return DropoffTypeNone, nil
	case "2":
		return DropoffTypePhoneAgency, nil
	case "3":
		return DropoffTypeCoordinateWithDriver, nil
	default:
		return DropoffTypeRegular, fmt.Errorf("invalid drop off type: %s", val)
	}
}

func parseTimepointType(val string) (TimepointType, error) {
	switch val {
	case "1", "":
		return TimepointTypeExact, nil
	case "0":
		return TimepointTypeApproximate, nil
	default:
		return TimepointTypeExact, fmt.Errorf("invalid timepoint type: %s", val)
	}
}

func parseExceptional(val string) (bool, error) {
	switch val {
	case "0", "":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("invalid value for exceptional: %s", val)
	}
}
