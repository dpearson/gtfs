package gtfs

import (
	"archive/zip"
	"fmt"
)

// FeedInfo specifies global information about a GTFS feed.
//
// Fields correspond directly to columns in feed_info.txt.
type FeedInfo struct {
	PublisherName string
	PublisherURL  string
	Lang          string
	StartDate     string
	EndDate       string
	Version       string
}

var feedInfoFields = map[string]bool{
	"feed_publisher_name": true,
	"feed_publisher_url":  true,
	"feed_lang":           true,
	"feed_start_date":     false,
	"feed_end_date":       false,
	"feed_version":        false,
}

func (g *GTFS) processFeedInfo(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	res, err := readCSVWithHeadings(rc, feedInfoFields)
	if err != nil {
		return err
	}

	if len(res) != 1 {
		return fmt.Errorf("Expected only one row of feed info, but there are %d rows", len(res))
	}

	row := res[0]
	g.FeedInfo = FeedInfo{
		PublisherName: row["feed_publisher_name"],
		PublisherURL:  row["feed_publisher_url"],
		Lang:          row["feed_lang"],
		StartDate:     row["feed_start_date"],
		EndDate:       row["feed_end_date"],
		Version:       row["feed_version"],
	}

	return nil
}
