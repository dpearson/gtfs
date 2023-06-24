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
	"calendar.txt":        false,
	"calendar_dates.txt":  false,
	"fare_attributes.txt": false,
	"fare_rules.txt":      false,
	"shapes.txt":          false,
	"frequencies.txt":     false,
	"transfers.txt":       false,
	"feed_info.txt":       false,
	"translations.txt":    false,
}

// GTFS represents a single GTFS feed.
type GTFS struct {
	Agencies     []*Agency
	Stops        []*Stop
	Routes       []*Route
	Services     []*Service
	Shapes       []*Shape
	Trips        []*Trip
	Fares        []*Fare
	Transfers    []*Transfer
	FeedInfo     FeedInfo
	Translations []*Translation

	agenciesByID     map[string]*Agency
	stopsByID        map[string]*Stop
	routesByID       map[string]*Route
	servicesByID     map[string]*Service
	shapesByID       map[string]*Shape
	tripsByID        map[string]*Trip
	faresByID        map[string]*Fare
	translationsByID map[string]map[string]*Translation
	strictMode       bool
}

// ParsingOptions specifies options used when parsing GTFS files.
type ParsingOptions struct {
	StrictMode bool
}

var defaultOptions = ParsingOptions{
	StrictMode: false,
}

// Load reads a GTFS feed, which is expected to be contained within a ZIP file,
// from filePath.
//
// This function loads a GTFS file as permissively as possible (i.e. all errors
// that can be ignored are ignored). For full control over options used when
// parsing, use LoadWithOptions instead.
func Load(filePath string) (*GTFS, error) {
	return LoadWithOptions(filePath, defaultOptions)
}

// LoadFromReader reads a GTFS feed from a *zip.Reader.
//
// This function loads a GTFS file as permissively as possible (i.e. all errors
// that can be ignored are ignored). For full control over options used when
// parsing, use LoadWithOptions instead.
func LoadFromReader(r *zip.Reader) (*GTFS, error) {
	return LoadFromReaderWithOptions(r, defaultOptions)
}

// LoadWithOptions reads a GTFS feed, which is expected to be contained within
// a ZIP file, from filePath using the specified options when parsing.
func LoadWithOptions(filePath string, opts ParsingOptions) (*GTFS, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close() // nolint: errcheck

	return LoadFromReaderWithOptions(&r.Reader, opts)
}

// LoadFromReaderWithOptions reads a GTFS feed from a *zip.Reader using the specified options when
// parsing.
func LoadFromReaderWithOptions(r *zip.Reader, opts ParsingOptions) (*GTFS, error) {
	g := &GTFS{
		strictMode: opts.StrictMode,
	}

	files := map[string]rcOpener{}
	for _, f := range r.File {
		if _, ok := validFilenames[f.Name]; !ok {
			continue
		}

		files[f.Name] = f
	}

	for name, required := range validFilenames {
		if _, ok := files[name]; !ok && required {
			return g, fmt.Errorf("no %s file found", name)
		}
	}

	err := g.doLoad(files)

	return g, err
}

func (g *GTFS) doLoad(files map[string]rcOpener) error {
	err := callWithOpenedReader(g.processAgencies, files["agency.txt"])
	if err != nil {
		return fmt.Errorf("error parsing agency.txt: %v", err)
	}

	err = callWithOpenedReader(g.processStops, files["stops.txt"])
	if err != nil {
		return fmt.Errorf("error parsing stops.txt: %v", err)
	}

	err = callWithOpenedReader(g.processRoutes, files["routes.txt"])
	if err != nil {
		return fmt.Errorf("error parsing routes.txt: %v", err)
	}

	f, hasCalendar := files["calendar.txt"]
	if hasCalendar {
		err = callWithOpenedReader(g.processServices, f)
		if err != nil {
			return fmt.Errorf("error parsing calendar.txt: %v", err)
		}
	}

	f, ok := files["calendar_dates.txt"]
	if ok {
		err = callWithOpenedReader(g.processServiceDates, f)
		if err != nil {
			return fmt.Errorf("error parsing calendar_dates.txt: %v", err)
		}
	} else if !hasCalendar {
		return fmt.Errorf("either calendar.txt or calendar_dates.txt is required")
	}

	f, ok = files["shapes.txt"]
	if ok {
		err = callWithOpenedReader(g.processShapes, f)
		if err != nil {
			return fmt.Errorf("error parsing shapes.txt: %v", err)
		}
	}

	err = callWithOpenedReader(g.processTrips, files["trips.txt"])
	if err != nil {
		return fmt.Errorf("error parsing trips.txt: %v", err)
	}

	err = callWithOpenedReader(g.processStopTimes, files["stop_times.txt"])
	if err != nil {
		return fmt.Errorf("error parsing stop_times.txt: %v", err)
	}

	f, ok = files["fare_attributes.txt"]
	if ok {
		err = callWithOpenedReader(g.processFares, f)
		if err != nil {
			return fmt.Errorf("error parsing fare_attributes.txt: %v", err)
		}

		f, ok = files["fare_rules.txt"]
		if ok {
			err = callWithOpenedReader(g.processFareRules, f)
			if err != nil {
				return fmt.Errorf("error parsing fare_rules.txt: %v", err)
			}
		}
	}

	f, ok = files["frequencies.txt"]
	if ok {
		err = callWithOpenedReader(g.processFrequencies, f)
		if err != nil {
			return fmt.Errorf("error parsing frequencies.txt: %v", err)
		}
	}

	f, ok = files["transfers.txt"]
	if ok {
		err = callWithOpenedReader(g.processTransfers, f)
		if err != nil {
			return fmt.Errorf("error parsing transfers.txt: %v", err)
		}
	}

	f, ok = files["feed_info.txt"]
	if ok {
		err = callWithOpenedReader(g.processFeedInfo, f)
		if err != nil {
			return fmt.Errorf("error parsing feed_info.txt: %v", err)
		}
	}

	f, ok = files["translations.txt"]
	if ok {
		err = callWithOpenedReader(g.processTranslations, f)
		if err != nil {
			return fmt.Errorf("error parsing translations.txt: %v", err)
		}
	}

	return nil
}
