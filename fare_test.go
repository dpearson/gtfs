package gtfs

import (
	"reflect"
	"testing"
)

func TestGTFS_fareByID(t *testing.T) {
	testFare1 := &Fare{
		ID: "test_fare_1",
	}
	type fields struct {
		Fares     []*Fare
		faresByID map[string]*Fare
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Fare
	}{
		{
			name: "Exists",
			fields: fields{
				Fares: []*Fare{
					testFare1,
				},
				faresByID: map[string]*Fare{
					"test_fare_1": testFare1,
				},
			},
			args: args{
				id: "test_fare_1",
			},
			want: testFare1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Fares: []*Fare{
					testFare1,
				},
				faresByID: map[string]*Fare{
					"test_fare_1": testFare1,
				},
			},
			args: args{
				id: "test_fare_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Fares:     nil,
				faresByID: nil,
			},
			args: args{
				id: "test_fare_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Fares:     tt.fields.Fares,
				faresByID: tt.fields.faresByID,
			}
			if got := g.fareByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.fareByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePaymentMethod(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    PaymentMethod
		wantErr bool
	}{
		{
			name:    "Empty",
			val:     "",
			want:    PaymentMethodOnBoard,
			wantErr: true,
		},
		{
			name:    "On Board",
			val:     "0",
			want:    PaymentMethodOnBoard,
			wantErr: false,
		},
		{
			name:    "Before Boarding",
			val:     "1",
			want:    PaymentMethodBeforeBoarding,
			wantErr: false,
		},
		{
			name:    "Invalid",
			val:     "2",
			want:    PaymentMethodOnBoard,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePaymentMethod(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePaymentMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePaymentMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
