package gtfs

import (
	"reflect"
	"testing"
)

func TestGTFS_routeByID(t *testing.T) {
	testRoute1 := &Route{
		ID: "test_route_1",
	}
	type fields struct {
		Routes     []*Route
		routesByID map[string]*Route
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Route
	}{
		{
			name: "Exists",
			fields: fields{
				Routes: []*Route{
					testRoute1,
				},
				routesByID: map[string]*Route{
					"test_route_1": testRoute1,
				},
			},
			args: args{
				id: "test_route_1",
			},
			want: testRoute1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Routes: []*Route{
					testRoute1,
				},
				routesByID: map[string]*Route{
					"test_route_1": testRoute1,
				},
			},
			args: args{
				id: "test_route_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Routes:     nil,
				routesByID: nil,
			},
			args: args{
				id: "test_route_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Routes:     tt.fields.Routes,
				routesByID: tt.fields.routesByID,
			}
			if got := g.routeByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.routeByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
