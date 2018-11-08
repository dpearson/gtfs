package gtfs

import (
	"archive/zip"
	"fmt"
	"sort"
	"strconv"
)

// A Shape is a collection of points describing a trip's path.
type Shape struct {
	ID     string
	Points []*ShapePoint
}

// A ShapePoint is a single point in a larger shape.
type ShapePoint struct {
	Latitude  float64
	Longitude float64
	Sequence  uint64
	Distance  float64
}

var shapeFields = map[string]bool{
	"shape_id":            true,
	"shape_pt_lat":        true,
	"shape_pt_lon":        true,
	"shape_pt_sequence":   true,
	"shape_dist_traveled": false,
}

func (g *GTFS) processShapes(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() // nolint: errcheck

	res, err := readCSVWithHeadings(rc, shapeFields)
	if err != nil {
		return err
	}

	shapePoints := map[string][]*ShapePoint{}
	for _, row := range res {
		id := row["shape_id"]
		lat, err := strconv.ParseFloat(row["shape_pt_lat"], 64)
		if err != nil {
			return fmt.Errorf("Invalid latitude: %v", err)
		}

		lon, err := strconv.ParseFloat(row["shape_pt_lon"], 64)
		if err != nil {
			return fmt.Errorf("Invalid longitude: %v", err)
		}

		distStr := row["shape_dist_traveled"]
		dist := 0.0
		if distStr != "" {
			dist, err = strconv.ParseFloat(distStr, 64)
			if err != nil {
				return fmt.Errorf("Invalid distance: %v", err)
			}
		}

		seq, err := strconv.ParseUint(row["shape_pt_sequence"], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid point sequence: %v", err)
		}

		pt := &ShapePoint{
			Latitude:  lat,
			Longitude: lon,
			Sequence:  seq,
			Distance:  dist,
		}

		shapePoints[id] = append(shapePoints[id], pt)
	}

	g.shapesByID = map[string]*Shape{}

	for id, pts := range shapePoints {
		sort.Slice(pts, func(i, j int) bool {
			return pts[i].Sequence < pts[j].Sequence
		})

		s := &Shape{
			ID:     id,
			Points: pts,
		}

		g.Shapes = append(g.Shapes, s)
		g.shapesByID[s.ID] = s
	}

	return nil
}

func (g *GTFS) shapeByID(id string) *Shape {
	return g.shapesByID[id]
}
