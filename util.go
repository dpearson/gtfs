package gtfs

import "fmt"

func parseBool(val string) (bool, error) {
	switch val {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("invalid value: %s", val)
	}
}
