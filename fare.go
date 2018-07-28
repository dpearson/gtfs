package gtfs

import (
	"archive/zip"
	"fmt"
	"strconv"
)

// A Fare is a single fare type.
//
// Fields correspond directly to columns in fares.txt.
type Fare struct {
	ID               string
	Price            string
	CurrencyType     string
	PaymentMethod    PaymentMethod
	Transfers        string // TODO: Parse me
	TransferDuration uint64

	Routes           []*Route
	OriginZones      []string
	DestinationZones []string
	ContainsZones    []string
}

// A PaymentMethod indicates where fares are paid.
type PaymentMethod int

const (
	// PaymentMethodOnBoard indicates that fares are paid on board the vehicle.
	PaymentMethodOnBoard PaymentMethod = iota

	//PaymentMethodBeforeBoarding indicates that fares are paid prior to
	// boarding the vehicle.
	PaymentMethodBeforeBoarding
)

var fareFields = map[string]bool{
	"fare_id":           true,
	"price":             true,
	"currency_type":     true,
	"payment_method":    true,
	"transfers":         true,
	"transfer_duration": false,
}

var fareRuleFields = map[string]bool{
	"fare_id":        true,
	"route_id":       false,
	"origin_id":      false,
	"destination_id": false,
	"contains_id":    false,
}

func (g *GTFS) processFares(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, fareFields)
	if err != nil {
		return err
	}

	g.faresByID = map[string]*Fare{}

	for _, row := range res {
		paymentMethod, err := parsePaymentMethod(row["payment_method"])
		if err != nil {
			return err
		}

		transferDuration := uint64(0)
		transferDurationStr := row["transfer_duration"]
		if transferDurationStr != "" {
			transferDuration, err = strconv.ParseUint(transferDurationStr, 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid transfer duration: %v", err)
			}
		}

		fare := &Fare{
			ID:               row["fare_id"],
			Price:            row["price"],
			CurrencyType:     row["currency_type"],
			PaymentMethod:    paymentMethod,
			Transfers:        row["transfers"],
			TransferDuration: transferDuration,
		}

		g.Fares = append(g.Fares, fare)
		g.faresByID[fare.ID] = fare
	}

	return nil
}

func (g *GTFS) processFareRules(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, fareRuleFields)
	if err != nil {
		return err
	}

	for _, row := range res {
		fare := g.fareByID(row["fare_id"])
		if fare == nil {
			return fmt.Errorf("Invalid fare ID: %s", row["fare_id"])
		}

		routeID := row["route_id"]
		if routeID != "" {
			r := g.routeByID(routeID)
			if r == nil {
				return fmt.Errorf("Invalid route ID: %s", row["route_id"])
			}

			fare.Routes = append(fare.Routes, r)
		}

		originID := row["origin_id"]
		if originID != "" {
			fare.OriginZones = append(fare.OriginZones, originID)
		}

		destID := row["destination_id"]
		if destID != "" {
			fare.DestinationZones = append(fare.DestinationZones, destID)
		}

		containsID := row["contains_id"]
		if containsID != "" {
			fare.ContainsZones = append(fare.ContainsZones, containsID)
		}
	}

	return nil
}

func (g *GTFS) fareByID(id string) *Fare {
	return g.faresByID[id]
}

func parsePaymentMethod(val string) (PaymentMethod, error) {
	switch val {
	case "0":
		return PaymentMethodOnBoard, nil
	case "1":
		return PaymentMethodBeforeBoarding, nil
	default:
		return PaymentMethodOnBoard, fmt.Errorf("Invalid payment method: %s", val)
	}
}
