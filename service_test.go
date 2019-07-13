package gtfs

import (
	"reflect"
	"testing"
)

func TestGTFS_serviceByID(t *testing.T) {
	testService1 := &Service{
		ID: "test_service_1",
	}
	type fields struct {
		Services     []*Service
		servicesByID map[string]*Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Service
	}{
		{
			name: "Exists",
			fields: fields{
				Services: []*Service{
					testService1,
				},
				servicesByID: map[string]*Service{
					"test_service_1": testService1,
				},
			},
			args: args{
				id: "test_service_1",
			},
			want: testService1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Services: []*Service{
					testService1,
				},
				servicesByID: map[string]*Service{
					"test_service_1": testService1,
				},
			},
			args: args{
				id: "test_service_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Services:     nil,
				servicesByID: nil,
			},
			args: args{
				id: "test_service_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Services:     tt.fields.Services,
				servicesByID: tt.fields.servicesByID,
			}
			if got := g.serviceByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.serviceByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
