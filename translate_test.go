package gtfs

import "testing"

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
