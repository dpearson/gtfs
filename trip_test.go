package gtfs

import (
	"reflect"
	"testing"
)

func Test_parseBikesAllowed(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    BikesAllowed
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    BikesAllowedUnknown,
			wantErr: false,
		},
		{
			name:    "Unknown",
			val:     "0",
			want:    BikesAllowedUnknown,
			wantErr: false,
		},
		{
			name:    "Yes",
			val:     "1",
			want:    BikesAllowedYes,
			wantErr: false,
		},
		{
			name:    "No",
			val:     "2",
			want:    BikesAllowedNo,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "3",
			want:    BikesAllowedUnknown,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBikesAllowed(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBikesAllowed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseBikesAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePickupType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    PickupType
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    PickupTypeRegular,
			wantErr: false,
		},
		{
			name:    "Regular",
			val:     "0",
			want:    PickupTypeRegular,
			wantErr: false,
		},
		{
			name:    "None",
			val:     "1",
			want:    PickupTypeNone,
			wantErr: false,
		},
		{
			name:    "Phone Agency",
			val:     "2",
			want:    PickupTypePhoneAgency,
			wantErr: false,
		},
		{
			name:    "Coordinate With Driver",
			val:     "3",
			want:    PickupTypeCoordinateWithDriver,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "4",
			want:    PickupTypeRegular,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePickupType(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePickupType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePickupType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDropoffType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    DropoffType
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    DropoffTypeRegular,
			wantErr: false,
		},
		{
			name:    "Regular",
			val:     "0",
			want:    DropoffTypeRegular,
			wantErr: false,
		},
		{
			name:    "None",
			val:     "1",
			want:    DropoffTypeNone,
			wantErr: false,
		},
		{
			name:    "Phone Agency",
			val:     "2",
			want:    DropoffTypePhoneAgency,
			wantErr: false,
		},
		{
			name:    "Coordinate With Driver",
			val:     "3",
			want:    DropoffTypeCoordinateWithDriver,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "4",
			want:    DropoffTypeRegular,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDropoffType(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDropoffType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDropoffType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTimepointType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    TimepointType
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    TimepointTypeExact,
			wantErr: false,
		},
		{
			name:    "Approximate",
			val:     "0",
			want:    TimepointTypeApproximate,
			wantErr: false,
		},
		{
			name:    "Exact",
			val:     "1",
			want:    TimepointTypeExact,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "2",
			want:    TimepointTypeExact,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimepointType(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimepointType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseTimepointType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseWheelchairAccessible(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    WheelchairAccessible
		wantErr bool
	}{
		{
			name: "Unknown (empty)",
			args: args{
				val: "",
			},
			want:    WheelchairAccessibleUnknown,
			wantErr: false,
		},
		{
			name: "Unknown (zero)",
			args: args{
				val: "0",
			},
			want:    WheelchairAccessibleUnknown,
			wantErr: false,
		},
		{
			name: "Yes",
			args: args{
				val: "1",
			},
			want:    WheelchairAccessibleYes,
			wantErr: false,
		},
		{
			name: "No",
			args: args{
				val: "2",
			},
			want:    WheelchairAccessibleNo,
			wantErr: false,
		},
		{
			name: "Invalid (3)",
			args: args{
				val: "3",
			},
			want:    WheelchairAccessibleUnknown,
			wantErr: true,
		},
		{
			name: "Invalid (foo)",
			args: args{
				val: "foo",
			},
			want:    WheelchairAccessibleUnknown,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWheelchairAccessible(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWheelchairAccessible() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseWheelchairAccessible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExceptional(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "False (empty)",
			args: args{
				val: "",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "False (zero)",
			args: args{
				val: "0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "True",
			args: args{
				val: "1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid (2)",
			args: args{
				val: "2",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid (foo)",
			args: args{
				val: "foo",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExceptional(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExceptional() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseExceptional() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGTFS_tripByID(t *testing.T) {
	testTrip1 := &Trip{
		ID: "test_trip_1",
	}
	type fields struct {
		Trips     []*Trip
		tripsByID map[string]*Trip
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Trip
	}{
		{
			name: "Exists",
			fields: fields{
				Trips: []*Trip{
					testTrip1,
				},
				tripsByID: map[string]*Trip{
					"test_trip_1": testTrip1,
				},
			},
			args: args{
				id: "test_trip_1",
			},
			want: testTrip1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Trips: []*Trip{
					testTrip1,
				},
				tripsByID: map[string]*Trip{
					"test_trip_1": testTrip1,
				},
			},
			args: args{
				id: "test_trip_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Trips:     nil,
				tripsByID: nil,
			},
			args: args{
				id: "test_trip_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Trips:     tt.fields.Trips,
				tripsByID: tt.fields.tripsByID,
			}
			if got := g.tripByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.tripByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
