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

func Test_parseRouteSortOrder(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "Empty String",
			args: args{
				val: "",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "0",
			args: args{
				val: "0",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "1",
			args: args{
				val: "1",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "-1",
			args: args{
				val: "-1",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "foo",
			args: args{
				val: "foo",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRouteSortOrder(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRouteSortOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRouteSortOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
