package gtfs

import (
	"fmt"
	"io"
)

// A Service is a schedule of service over one or more routes.
//
// Fields correspond to columns in calendar.txt and calendar_dates.txt.
type Service struct {
	ID        string
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
	StartDate string
	EndDate   string

	AdditionalDates []string
	ExceptDates     []string
}

var serviceFields = map[string]bool{
	"service_id": true,
	"monday":     true,
	"tuesday":    true,
	"wednesday":  true,
	"thursday":   true,
	"friday":     true,
	"saturday":   true,
	"sunday":     true,
	"start_date": true,
	"end_date":   true,
}

var serviceDateFields = map[string]bool{
	"service_id":     true,
	"date":           true,
	"exception_type": true,
}

func (g *GTFS) processServices(r io.Reader) error {
	res, err := readCSVWithHeadings(r, serviceFields, g.strictMode)
	if err != nil {
		return err
	}

	if g.servicesByID == nil {
		g.servicesByID = map[string]*Service{}
	}

	for _, row := range res {
		monday, err := parseBool(row["monday"])
		if err != nil {
			return err
		}

		tuesday, err := parseBool(row["tuesday"])
		if err != nil {
			return err
		}

		wednesday, err := parseBool(row["wednesday"])
		if err != nil {
			return err
		}

		thursday, err := parseBool(row["thursday"])
		if err != nil {
			return err
		}

		friday, err := parseBool(row["friday"])
		if err != nil {
			return err
		}

		saturday, err := parseBool(row["saturday"])
		if err != nil {
			return err
		}

		sunday, err := parseBool(row["sunday"])
		if err != nil {
			return err
		}

		s := &Service{
			ID:        row["service_id"],
			Monday:    monday,
			Tuesday:   tuesday,
			Wednesday: wednesday,
			Thursday:  thursday,
			Friday:    friday,
			Saturday:  saturday,
			Sunday:    sunday,
			StartDate: row["start_date"],
			EndDate:   row["end_date"],
		}

		g.Services = append(g.Services, s)
		g.servicesByID[s.ID] = s
	}

	return nil
}

func (g *GTFS) processServiceDates(r io.Reader) error {
	res, err := readCSVWithHeadings(r, serviceDateFields, g.strictMode)
	if err != nil {
		return err
	}

	if g.servicesByID == nil {
		g.servicesByID = map[string]*Service{}
	}

	for _, row := range res {
		id := row["service_id"]
		date := row["date"]
		exceptionType := row["exception_type"]

		s := g.serviceByID(id)
		if s == nil {
			s = &Service{
				ID: id,
			}
			g.Services = append(g.Services, s)
			g.servicesByID[s.ID] = s
		}

		switch exceptionType {
		case "1":
			s.AdditionalDates = append(s.AdditionalDates, date)
		case "2":
			s.ExceptDates = append(s.ExceptDates, date)
		default:
			return fmt.Errorf("invalid exception_type: %s", exceptionType)
		}
	}

	return nil
}

func (g *GTFS) serviceByID(id string) *Service {
	return g.servicesByID[id]
}
