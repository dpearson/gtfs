package gtfs

import (
	"archive/zip"
	"fmt"
)

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

func (g *GTFS) processServices(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, serviceFields)
	if err != nil {
		return err
	}

	g.servicesByID = map[string]*Service{}

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

func (g *GTFS) processServiceDates(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, serviceDateFields)
	if err != nil {
		return err
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
		}

		switch exceptionType {
		case "1":
			s.AdditionalDates = append(s.AdditionalDates, date)
		case "2":
			s.ExceptDates = append(s.ExceptDates, date)
		default:
			return fmt.Errorf("Invalid exception_type: %s", exceptionType)
		}
	}

	return nil
}

func (g *GTFS) serviceByID(id string) *Service {
	return g.servicesByID[id]
}
