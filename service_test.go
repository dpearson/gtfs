package gtfs

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const testServicesCSVValid = `service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date
1,1,1,1,1,1,0,0,20190101,20191231
2,0,0,0,0,0,1,1,20190101,20191231`
const testServiceDatesCSVValid = `service_id,date,exception_type
1,20190720,1
3,20190721,1
2,20190727,2`

var testService1 = &Service{
	ID:        "1",
	Monday:    true,
	Tuesday:   true,
	Wednesday: true,
	Thursday:  true,
	Friday:    true,
	Saturday:  false,
	Sunday:    false,
	StartDate: "20190101",
	EndDate:   "20191231",
}
var testService2 = &Service{
	ID:        "2",
	Monday:    false,
	Tuesday:   false,
	Wednesday: false,
	Thursday:  false,
	Friday:    false,
	Saturday:  true,
	Sunday:    true,
	StartDate: "20190101",
	EndDate:   "20191231",
}
var testService3 = &Service{
	ID:        "3",
	Monday:    false,
	Tuesday:   false,
	Wednesday: false,
	Thursday:  false,
	Friday:    false,
	Saturday:  false,
	Sunday:    false,
	StartDate: "",
	EndDate:   "",
}

func TestGTFS_processServices(t *testing.T) {
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
		wantServices     []*Service
		wantServicesByID map[string]*Service
	}{
		{
			name: "Valid",
			fields: fields{
				strictMode: false,
			},
			args: args{
				r: strings.NewReader(testServicesCSVValid),
			},
			wantErr: false,
			wantServices: []*Service{
				testService1,
				testService2,
			},
			wantServicesByID: map[string]*Service{
				"1": testService1,
				"2": testService2,
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
			wantServices:     nil,
			wantServicesByID: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				strictMode: tt.fields.strictMode,
			}
			if err := g.processServices(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processServices() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Services, tt.wantServices) {
				t.Errorf("GTFS.processServices() Services = %v, wantServices %v", g.Services, tt.wantServices)
			}
			if !reflect.DeepEqual(g.servicesByID, tt.wantServicesByID) {
				t.Errorf("GTFS.processServices() servicesByID = %v, wantServicesByID %v", g.servicesByID, tt.wantServicesByID)
			}
		})
	}
}

func TestGTFS_processServiceDates(t *testing.T) {
	testService1After := &Service{
		ID:        testService1.ID,
		Monday:    testService1.Monday,
		Tuesday:   testService1.Tuesday,
		Wednesday: testService1.Wednesday,
		Thursday:  testService1.Thursday,
		Friday:    testService1.Friday,
		Saturday:  testService1.Saturday,
		Sunday:    testService1.Sunday,
		StartDate: testService1.StartDate,
		EndDate:   testService1.EndDate,
		AdditionalDates: []string{
			"20190720",
		},
	}
	testService2After := &Service{
		ID:        testService2.ID,
		Monday:    testService2.Monday,
		Tuesday:   testService2.Tuesday,
		Wednesday: testService2.Wednesday,
		Thursday:  testService2.Thursday,
		Friday:    testService2.Friday,
		Saturday:  testService2.Saturday,
		Sunday:    testService2.Sunday,
		StartDate: testService2.StartDate,
		EndDate:   testService2.EndDate,
		ExceptDates: []string{
			"20190727",
		},
	}
	testService3After := &Service{
		ID:        testService3.ID,
		Monday:    testService3.Monday,
		Tuesday:   testService3.Tuesday,
		Wednesday: testService3.Wednesday,
		Thursday:  testService3.Thursday,
		Friday:    testService3.Friday,
		Saturday:  testService3.Saturday,
		Sunday:    testService3.Sunday,
		StartDate: testService3.StartDate,
		EndDate:   testService3.EndDate,
		AdditionalDates: []string{
			"20190721",
		},
	}
	type fields struct {
		Services     []*Service
		servicesByID map[string]*Service
		strictMode   bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantErr          bool
		wantServices     []*Service
		wantServicesByID map[string]*Service
	}{
		{
			name: "Valid",
			fields: fields{
				Services: []*Service{
					testService1,
					testService2,
				},
				servicesByID: map[string]*Service{
					"1": testService1,
					"2": testService2,
				},
				strictMode: true,
			},
			args: args{
				r: strings.NewReader(testServiceDatesCSVValid),
			},
			wantErr: false,
			wantServices: []*Service{
				testService1After,
				testService2After,
				testService3After,
			},
			wantServicesByID: map[string]*Service{
				"1": testService1After,
				"2": testService2After,
				"3": testService3After,
			},
		},
		{
			name: "Empty",
			fields: fields{
				Services: []*Service{
					testService1,
					testService2,
				},
				servicesByID: map[string]*Service{
					"1": testService1,
					"2": testService2,
				},
				strictMode: true,
			},
			args: args{
				r: strings.NewReader(""),
			},
			wantErr: true,
			wantServices: []*Service{
				testService1,
				testService2,
			},
			wantServicesByID: map[string]*Service{
				"1": testService1,
				"2": testService2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Services:     tt.fields.Services,
				servicesByID: tt.fields.servicesByID,
				strictMode:   tt.fields.strictMode,
			}
			if err := g.processServiceDates(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("GTFS.processServiceDates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g.Services, tt.wantServices) {
				t.Errorf("GTFS.processServices() Services = %v, wantServices %v", g.Services, tt.wantServices)
			}
			if !reflect.DeepEqual(g.servicesByID, tt.wantServicesByID) {
				t.Errorf("GTFS.processServices() servicesByID = %v, wantServicesByID %v", g.servicesByID, tt.wantServicesByID)
			}
		})
	}
}

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
