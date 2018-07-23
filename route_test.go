package gtfs

import "testing"

func Test_parseRouteType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    RouteType
		wantErr bool
	}{
		{
			name:    "Light Rail",
			val:     "0",
			want:    RouteTypeLightRail,
			wantErr: false,
		},
		{
			name:    "Subway",
			val:     "1",
			want:    RouteTypeSubway,
			wantErr: false,
		},
		{
			name:    "Intercity Rail",
			val:     "2",
			want:    RouteTypeRail,
			wantErr: false,
		},
		{
			name:    "Bus",
			val:     "3",
			want:    RouteTypeBus,
			wantErr: false,
		},
		{
			name:    "Ferry",
			val:     "4",
			want:    RouteTypeFerry,
			wantErr: false,
		},
		{
			name:    "Cable Car",
			val:     "5",
			want:    RouteTypeCableCar,
			wantErr: false,
		},
		{
			name:    "Gondola",
			val:     "6",
			want:    RouteTypeGondola,
			wantErr: false,
		},
		{
			name:    "Funicular",
			val:     "7",
			want:    RouteTypeFunicular,
			wantErr: false,
		},
		{
			name:    "Empty",
			val:     "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid",
			val:     "8",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRouteType(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRouteType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRouteType() = %v, want %v", got, tt.want)
			}
		})
	}
}
