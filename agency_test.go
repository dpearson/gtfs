package gtfs

import (
	"reflect"
	"testing"
)

func TestGTFS_agencyByID(t *testing.T) {
	testAgency1 := &Agency{
		ID: "test_agency_1",
	}
	type fields struct {
		Agencies     []*Agency
		agenciesByID map[string]*Agency
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Agency
	}{
		{
			name: "Exists",
			fields: fields{
				Agencies: []*Agency{
					testAgency1,
				},
				agenciesByID: map[string]*Agency{
					"test_agency_1": testAgency1,
				},
			},
			args: args{
				id: "test_agency_1",
			},
			want: testAgency1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Agencies: []*Agency{
					testAgency1,
				},
				agenciesByID: map[string]*Agency{
					"test_agency_1": testAgency1,
				},
			},
			args: args{
				id: "test_agency_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Agencies:     nil,
				agenciesByID: nil,
			},
			args: args{
				id: "test_agency_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Agencies:     tt.fields.Agencies,
				agenciesByID: tt.fields.agenciesByID,
			}
			if got := g.agencyByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.agencyByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
