package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testTranslationsCSVValid = `trans_id,lang,translation
foo,en,bar
foo,es,baz`

func TestGTFS_Translate(t *testing.T) {
	testTranslation1 := &Translation{
		ID:          "station-001",
		Language:    "en",
		Translation: "City Center",
	}

	testTranslation2 := &Translation{
		ID:          "station-001",
		Language:    "fr",
		Translation: "Centre-ville",
	}

	type fields struct {
		Agencies         []*Agency
		Stops            []*Stop
		Routes           []*Route
		Services         []*Service
		Shapes           []*Shape
		Trips            []*Trip
		Fares            []*Fare
		Transfers        []*Transfer
		FeedInfo         FeedInfo
		Translations     []*Translation
		agenciesByID     map[string]*Agency
		stopsByID        map[string]*Stop
		routesByID       map[string]*Route
		servicesByID     map[string]*Service
		shapesByID       map[string]*Shape
		tripsByID        map[string]*Trip
		faresByID        map[string]*Fare
		translationsByID map[string]map[string]*Translation
	}
	type args struct {
		sourceStr string
		lang      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Translation Available (EN)",
			fields: fields{
				Translations: []*Translation{
					testTranslation1,
					testTranslation2,
				},
				translationsByID: map[string]map[string]*Translation{
					"station-001": map[string]*Translation{
						"en": testTranslation1,
						"fr": testTranslation2,
					},
				},
			},
			args: args{
				sourceStr: "station-001",
				lang:      "en",
			},
			want: "City Center",
		},
		{
			name: "Translation Available (FR)",
			fields: fields{
				Translations: []*Translation{
					testTranslation1,
					testTranslation2,
				},
				translationsByID: map[string]map[string]*Translation{
					"station-001": map[string]*Translation{
						"en": testTranslation1,
						"fr": testTranslation2,
					},
				},
			},
			args: args{
				sourceStr: "station-001",
				lang:      "fr",
			},
			want: "Centre-ville",
		},
		{
			name: "Translation Unavailable (DE)",
			fields: fields{
				Translations: []*Translation{
					testTranslation1,
					testTranslation2,
				},
				translationsByID: map[string]map[string]*Translation{
					"station-001": map[string]*Translation{
						"en": testTranslation1,
						"fr": testTranslation2,
					},
				},
			},
			args: args{
				sourceStr: "station-001",
				lang:      "de",
			},
			want: "station-001",
		},
		{
			name: "Unknown Phrase",
			fields: fields{
				Translations: []*Translation{
					testTranslation1,
					testTranslation2,
				},
				translationsByID: map[string]map[string]*Translation{
					"station-001": map[string]*Translation{
						"en": testTranslation1,
						"fr": testTranslation2,
					},
				},
			},
			args: args{
				sourceStr: "station-002",
				lang:      "en",
			},
			want: "station-002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Agencies:         tt.fields.Agencies,
				Stops:            tt.fields.Stops,
				Routes:           tt.fields.Routes,
				Services:         tt.fields.Services,
				Shapes:           tt.fields.Shapes,
				Trips:            tt.fields.Trips,
				Fares:            tt.fields.Fares,
				Transfers:        tt.fields.Transfers,
				FeedInfo:         tt.fields.FeedInfo,
				Translations:     tt.fields.Translations,
				agenciesByID:     tt.fields.agenciesByID,
				stopsByID:        tt.fields.stopsByID,
				routesByID:       tt.fields.routesByID,
				servicesByID:     tt.fields.servicesByID,
				shapesByID:       tt.fields.shapesByID,
				tripsByID:        tt.fields.tripsByID,
				faresByID:        tt.fields.faresByID,
				translationsByID: tt.fields.translationsByID,
			}
			if got := g.Translate(tt.args.sourceStr, tt.args.lang); got != tt.want {
				t.Errorf("GTFS.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGTFS_processTranslations(t *testing.T) {
	testTranslation1 := &Translation{
		ID:          "foo",
		Language:    "en",
		Translation: "bar",
	}
	testTranslation2 := &Translation{
		ID:          "foo",
		Language:    "es",
		Translation: "baz",
	}

	type fields struct {
		strictMode bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantErr              bool
		wantTranslations     []*Translation
		wantTranslationsById map[string]map[string]*Translation
	}{
		{
			name: "Valid",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testTranslationsCSVValid),
			},
			wantErr: false,
			wantTranslations: []*Translation{
				testTranslation1,
				testTranslation2,
			},
			wantTranslationsById: map[string]map[string]*Translation{
				"foo": map[string]*Translation{
					"en": testTranslation1,
					"es": testTranslation2,
				},
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
			wantErr:              true,
			wantTranslations:     nil,
			wantTranslationsById: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				strictMode: tt.fields.strictMode,
			}
			if err := g.processTranslations(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processTranslations() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Translations, tt.wantTranslations) {
				t.Errorf("GTFS.processTranslations() Translations = %v, wantTranslations %v", g.Translations, tt.wantTranslations)
			}
			if !reflect.DeepEqual(g.translationsByID, tt.wantTranslationsById) {
				t.Errorf("GTFS.processTranslations() translationsByID = %v, wantTranslationsById %v", g.translationsByID, tt.wantTranslationsById)
			}
		})
	}
}
