package gtfs

import (
	"reflect"
	"testing"
)

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
