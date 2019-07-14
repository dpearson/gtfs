package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testShapesCSVValid = `shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
1,5.2,6,1
1,5,6.1,3,100
1,5.1,6.05,2,50.2`

const testShapesCSVInvalidLatitude = `shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
1,foo,6,1,0`
const testShapesCSVInvalidLongitude = `shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
1,5.2,foo,1,0`
const testShapesCSVInvalidSequence = `shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
1,5.2,6,-1,0`
const testShapesCSVInvalidDistance = `shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
1,5.2,6,1,foo`

func TestGTFS_processShapes(t *testing.T) {
	testShape1 := &Shape{
		ID: "1",
		Points: []*ShapePoint{
			{
				Latitude:  5.2,
				Longitude: 6.0,
				Sequence:  1,
				Distance:  0,
			},
			{
				Latitude:  5.1,
				Longitude: 6.05,
				Sequence:  2,
				Distance:  50.2,
			},
			{
				Latitude:  5,
				Longitude: 6.1,
				Sequence:  3,
				Distance:  100,
			},
		},
	}

	type fields struct {
		strictMode bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantShapes     []*Shape
		wantShapesById map[string]*Shape
	}{
		{
			name: "Valid",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testShapesCSVValid),
			},
			wantErr: false,
			wantShapes: []*Shape{
				testShape1,
			},
			wantShapesById: map[string]*Shape{
				"1": testShape1,
			},
		},
		{
			name: "Empty",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(""),
			},
			wantErr:        true,
			wantShapes:     nil,
			wantShapesById: nil,
		},
		{
			name: "Invalid Latitude",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testShapesCSVInvalidLatitude),
			},
			wantErr:        true,
			wantShapes:     nil,
			wantShapesById: nil,
		},
		{
			name: "Invalid Longitude",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testShapesCSVInvalidLongitude),
			},
			wantErr:        true,
			wantShapes:     nil,
			wantShapesById: nil,
		},
		{
			name: "Invalid Sequence",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testShapesCSVInvalidSequence),
			},
			wantErr:        true,
			wantShapes:     nil,
			wantShapesById: nil,
		},
		{
			name: "Invalid Distance",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testShapesCSVInvalidDistance),
			},
			wantErr:        true,
			wantShapes:     nil,
			wantShapesById: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				strictMode: tt.fields.strictMode,
			}
			if err := g.processShapes(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processShapes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Shapes, tt.wantShapes) {
				t.Errorf("GTFS.processShapes() Shapes = %v, wantShapes %v", g.Shapes, tt.wantShapes)
			}
			if !reflect.DeepEqual(g.shapesByID, tt.wantShapesById) {
				t.Errorf("GTFS.processShapes() shapesByID = %v, wantShapesById %v", g.shapesByID, tt.wantShapesById)
			}
		})
	}
}

func TestGTFS_shapeByID(t *testing.T) {
	testShape1 := &Shape{
		ID: "test_shape_1",
	}
	type fields struct {
		Shapes     []*Shape
		shapesByID map[string]*Shape
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Shape
	}{
		{
			name: "Exists",
			fields: fields{
				Shapes: []*Shape{
					testShape1,
				},
				shapesByID: map[string]*Shape{
					"test_shape_1": testShape1,
				},
			},
			args: args{
				id: "test_shape_1",
			},
			want: testShape1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Shapes: []*Shape{
					testShape1,
				},
				shapesByID: map[string]*Shape{
					"test_shape_1": testShape1,
				},
			},
			args: args{
				id: "test_shape_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Shapes:     nil,
				shapesByID: nil,
			},
			args: args{
				id: "test_shape_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Shapes:     tt.fields.Shapes,
				shapesByID: tt.fields.shapesByID,
			}
			if got := g.shapeByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.shapeByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
