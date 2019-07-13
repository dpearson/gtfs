package gtfs

import (
	"reflect"
	"testing"
)

func TestGTFS_shapeByID(t *testing.T) {
	testShape1 := &Shape{
		ID: "test_shape_1",
	}
	type fields struct {
		Shapes     []*Shape
		shapesByID map[string]*Shape
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Shape
	}{
		{
			name: "Exists",
			fields: fields{
				Shapes: []*Shape{
					testShape1,
				},
				shapesByID: map[string]*Shape{
					"test_shape_1": testShape1,
				},
			},
			args: args{
				id: "test_shape_1",
			},
			want: testShape1,
		},
		{
			name: "Doesn't Exist",
			fields: fields{
				Shapes: []*Shape{
					testShape1,
				},
				shapesByID: map[string]*Shape{
					"test_shape_1": testShape1,
				},
			},
			args: args{
				id: "test_shape_2",
			},
			want: nil,
		},
		{
			name: "nil map",
			fields: fields{
				Shapes:     nil,
				shapesByID: nil,
			},
			args: args{
				id: "test_shape_1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GTFS{
				Shapes:     tt.fields.Shapes,
				shapesByID: tt.fields.shapesByID,
			}
			if got := g.shapeByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GTFS.shapeByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
