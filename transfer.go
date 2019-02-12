package gtfs

import (
	"archive/zip"
	"fmt"
	"strconv"
)

// A Transfer is specific transfer between two stops.
//
// Fields correspond to columns in transfers.txt.
type Transfer struct {
	From                *Stop
	To                  *Stop
	Type                TransferType
	MinimumTransferTime uint64
}

// TransferType specifies the specific type of a transfer.
type TransferType int

const (
	// TransferTypeRecommended indicates a recommended transfer point between
	// two routes.
	TransferTypeRecommended TransferType = iota

	// TransferTypeTimed indicates that departures will wait for arriving
	// passengers.
	TransferTypeTimed

	// TransferTypeMinimumTime indicates that a minimum amount of time is
	// needed to transfer.
	//
	// Transfers with this type will have MinimumTransferTime set.
	TransferTypeMinimumTime

	// TransferTypeNone indicates that a transfer between those two stops is not
	// possible.
	TransferTypeNone
)

var transferFields = map[string]bool{
	"from_stop_id":      true,
	"to_stop_id":        true,
	"transfer_type":     true,
	"min_transfer_time": false,
}

func (g *GTFS) processTransfers(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() // nolint: errcheck

	res, err := readCSVWithHeadings(rc, transferFields, g.strictMode)
	if err != nil {
		return err
	}

	for _, row := range res {
		minTimeStr := row["min_transfer_time"]
		minTime := uint64(0)
		if minTimeStr != "" {
			minTime, err = strconv.ParseUint(minTimeStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid min_transfer_time: %v", err)
			}
		}

		transferType, err := parseTransferType(row["transfer_type"])
		if err != nil {
			return err
		}

		t := &Transfer{
			From:                g.stopByID(row["from_stop_id"]),
			To:                  g.stopByID(row["to_stop_id"]),
			Type:                transferType,
			MinimumTransferTime: minTime,
		}

		g.Transfers = append(g.Transfers, t)
	}

	return nil
}

func parseTransferType(val string) (TransferType, error) {
	switch val {
	case "0", "":
		return TransferTypeRecommended, nil
	case "1":
		return TransferTypeTimed, nil
	case "2":
		return TransferTypeMinimumTime, nil
	case "3":
		return TransferTypeNone, nil
	default:
		return TransferTypeRecommended, fmt.Errorf("invalid transfer type: %s", val)
	}
}
