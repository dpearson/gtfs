package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testStopsCSVValid = `stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station,stop_timezone,wheelchair_boarding,platform_code
1,abc,Test Stop,A test stop,5.1,6,1,https://example/com/stops/abc,0,2,America/Chicago,0,
2,def,Test Station,A test station,5,6.1,1,https://example/com/stops/def,1,,America/Chicago,0,`

func TestGTFS_processStops(t *testing.T) {
	testStation1 := &Stop{
		ID:                 "2",
		Code:               "def",
		Name:               "Test Station",
		Description:        "A test station",
		Latitude:           5,
		Longitude:          6.1,
		ZoneID:             "1",
		URL:                "https://example/com/stops/def",
		LocationType:       LocationTypeStation,
		parentStationID:    "",
		ParentStation:      nil,
		Timezone:           "America/Chicago",
		WheelchairBoarding: "0",
		PlatformCode:       "",
	}
	testStop1 := &Stop{
		ID:                 "1",
		Code:               "abc",
		Name:               "Test Stop",
		Description:        "A test stop",
		Latitude:           5.1,
		Longitude:          6,
		ZoneID:             "1",
		URL:                "https://example/com/stops/abc",
		LocationType:       LocationTypeStop,
		parentStationID:    "2",
		ParentStation:      testStation1,
		Timezone:           "America/Chicago",
		WheelchairBoarding: "0",
		PlatformCode:       "",
	}
	type fields struct {
		strictMode bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantStops     []*Stop
		wantStopsByID map[string]*Stop
	}{
		{
			name: "Valid",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testStopsCSVValid),
			},
			wantErr: false,
			wantStops: []*Stop{
				testStop1,
				testStation1,
			},
			wantStopsByID: map[string]*Stop{
				"1": testStop1,
				"2": testStation1,
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
			wantErr:       true,
			wantStops:     nil,
			wantStopsByID: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				strictMode: tt.fields.strictMode,
			}
			if err := g.processStops(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processStops() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Stops, tt.wantStops) {
				t.Errorf("GTFS.processStops() Stops = %v, wantStops %v", g.Stops, tt.wantStops)
			}
			if !reflect.DeepEqual(g.stopsByID, tt.wantStopsByID) {
				t.Errorf("GTFS.processStops() stopsByID = %v, wantStopsByID %v", g.stopsByID, tt.wantStopsByID)
			}
		})
	}
}

func TestGTFS_stopByID(t *testing.T) {
	testStop1 := &Stop{
		ID: "test_stop_1",
	}
	type fields struct {
		Stops     []*Stop
		stopsByID map[string]*Stop
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Stop
	}{
		{
			name: "Exists",
			fields: fields{
				Stops: []*Stop{
					testStop1,
				},
				stopsByID: map[string]*Stop{
					"test_stop_1": testStop1,
				},
			},
			args: args{
				id: "test_stop_1",
			},
			want: testStop1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Stops: []*Stop{
					testStop1,
				},
				stopsByID: map[string]*Stop{
					"test_stop_1": testStop1,
				},
			},
			args: args{
				id: "test_stop_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Stops:     nil,
				stopsByID: nil,
			},
			args: args{
				id: "test_stop_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Stops:     tt.fields.Stops,
				stopsByID: tt.fields.stopsByID,
			}
			if got := g.stopByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.stopByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
