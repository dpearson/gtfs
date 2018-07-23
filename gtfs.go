package gtfs

import (
	"archive/zip"
	"fmt"
)

var validFilenames = map[string]bool{
	"agency.txt":          true,
	"stops.txt":           true,
	"routes.txt":          true,
	"trips.txt":           true,
	"stop_times.txt":      true,
	"calendar.txt":        true,
	"calendar_dates.txt":  false,
	"fare_attributes.txt": false,
	"fare_rules.txt":      false,
	"shapes.txt":          false,
	"frequencies.txt":     false,
	"transfers.txt":       false,
	"feed_info.txt":       false,
}

type GTFS struct {
	Agencies  []*Agency
	Stops     []*Stop
	Routes    []*Route
	Services  []*Service
	Shapes    []*Shape
	Trips     []*Trip
	Fares     []*Fare
	Transfers []*Transfer
	FeedInfo  FeedInfo

	agenciesByID map[string]*Agency
	stopsByID    map[string]*Stop
	routesByID   map[string]*Route
	servicesByID map[string]*Service
	shapesByID   map[string]*Shape
	tripsByID    map[string]*Trip
	faresByID    map[string]*Fare
}

func Load(filePath string) (*GTFS, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	g := &GTFS{}

	files := map[string]*zip.File{}
	for _, f := range r.File {
		if _, ok := validFilenames[f.Name]; !ok {
			return g, fmt.Errorf("Invalid filename: %s", f.Name)
		}

		files[f.Name] = f
	}

	for name, required := range validFilenames {
		if _, ok := files[name]; !ok && required {
			return g, fmt.Errorf("No %s file found", name)
		}
	}

	err = g.processAgencies(files["agency.txt"])
	if err != nil {
		return g, err
	}

	err = g.processStops(files["stops.txt"])
	if err != nil {
		return g, err
	}

	err = g.processRoutes(files["routes.txt"])
	if err != nil {
		return g, err
	}

	err = g.processServices(files["calendar.txt"])
	if err != nil {
		return g, err
	}

	f, ok := files["calendar_dates.txt"]
	if ok {
		err = g.processServiceDates(f)
		if err != nil {
			return g, err
		}
	}

	f, ok = files["shapes.txt"]
	if ok {
		err = g.processShapes(f)
		if err != nil {
			return g, err
		}
	}

	err = g.processTrips(files["trips.txt"])
	if err != nil {
		return g, err
	}

	err = g.processStopTimes(files["stop_times.txt"])
	if err != nil {
		return g, err
	}

	f, ok = files["fare_attributes.txt"]
	if ok {
		err = g.processFares(f)
		if err != nil {
			return g, err
		}

		f, ok = files["fare_rules.txt"]
		if ok {
			err = g.processFareRules(f)
			if err != nil {
				return g, err
			}
		}
	}

	f, ok = files["frequencies.txt"]
	if ok {
		err = g.processFrequencies(f)
		if err != nil {
			return g, err
		}
	}

	f, ok = files["transfers.txt"]
	if ok {
		err = g.processTransfers(f)
		if err != nil {
			return g, err
		}
	}

	f, ok = files["feed_info.txt"]
	if ok {
		err = g.processFeedInfo(f)
		if err != nil {
			return g, err
		}
	}

	return g, nil
}
