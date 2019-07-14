package gtfs

import (
	"fmt"
	"io"
)

type rcOpener interface {
	Open() (io.ReadCloser, error)
}

func callWithOpenedReader(opener rcOpener, fn func(io.Reader) error) error {
	rc, err := opener.Open()
	if err != nil {
		return err
	}
	defer rc.Close() // nolint: errcheck

	return fn(rc)
}

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
