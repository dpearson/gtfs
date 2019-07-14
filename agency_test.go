package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testAgenciesCSVValid = `agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email
1,"Test Agency",https://example.com,America/Toronto,en,0000000000,https://example.com/fares,test@example.com`
const testAgenciesCSVValidWithoutID = `agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email
"Test Agency",https://example.com,America/Toronto,en,0000000000,https://example.com/fares,test@example.com`

const testAgenciesCSVInvalidWithoutID = `agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email
"Test Agency",https://example.com,America/Toronto,en,0000000000,https://example.com/fares,test@example.com
"Test Agency",https://example.com,America/Toronto,en,0000000000,https://example.com/fares,test@example.com`

func TestGTFS_processAgencies(t *testing.T) {
	testAgency1 := &Agency{
		ID:       "1",
		Name:     "Test Agency",
		URL:      "https://example.com",
		Timezone: "America/Toronto",
		Lang:     "en",
		Phone:    "0000000000",
		FareURL:  "https://example.com/fares",
		Email:    "test@example.com",
	}

	testAgency2 := &Agency{
		ID:       "",
		Name:     "Test Agency",
		URL:      "https://example.com",
		Timezone: "America/Toronto",
		Lang:     "en",
		Phone:    "0000000000",
		FareURL:  "https://example.com/fares",
		Email:    "test@example.com",
	}

	type fields struct {
		strictMode bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantErr          bool
		wantAgencies     []*Agency
		wantAgenciesById map[string]*Agency
	}{
		{
			name: "Valid (with ID)",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testAgenciesCSVValid),
			},
			wantErr: false,
			wantAgencies: []*Agency{
				testAgency1,
			},
			wantAgenciesById: map[string]*Agency{
				"1": testAgency1,
			},
		},
		{
			name: "Valid (without ID)",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testAgenciesCSVValidWithoutID),
			},
			wantErr: false,
			wantAgencies: []*Agency{
				testAgency2,
			},
			wantAgenciesById: map[string]*Agency{
				"": testAgency2,
			},
		},
		{
			name: "Invalid (without ID)",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testAgenciesCSVInvalidWithoutID),
			},
			wantErr: true,
			wantAgencies: []*Agency{
				testAgency2,
				testAgency2,
			},
			wantAgenciesById: map[string]*Agency{
				"": testAgency2,
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
			wantErr:          true,
			wantAgencies:     nil,
			wantAgenciesById: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				strictMode: tt.fields.strictMode,
			}
			if err := g.processAgencies(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processAgencies() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Agencies, tt.wantAgencies) {
				t.Errorf("GTFS.processAgencies() Agencies = %v, wantAgencies %v", g.Agencies, tt.wantAgencies)
			}
			if !reflect.DeepEqual(g.agenciesByID, tt.wantAgenciesById) {
				t.Errorf("GTFS.processAgencies() agenciesByID = %v, wantAgenciesById %v", g.agenciesByID, tt.wantAgenciesById)
			}
		})
	}
}

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
