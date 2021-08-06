package scur_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/vvakame/spatk/scur"
)

var _ spanner.Encoder = FooID("")
var _ spanner.Decoder = (*FooID)(nil)

type FooID string

func (id FooID) EncodeSpanner() (interface{}, error) {
	return string(id), nil
}

func (id *FooID) DecodeSpanner(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return fmt.Errorf("unexpected id type: %T", input)
	}

	*id = FooID(s)
	return nil
}

func TestCursor_EncodeDecodeParameters(t *testing.T) {
	tests := []struct {
		name     string
		cursor   scur.Cursor
		expected scur.Cursor // nil の場合 cursor と同じものにする
	}{
		{
			name: "1 column",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "String", Order: scur.OrderAsc, Value: "foobar"},
			},
		},
		{
			name: "3 columns",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "String", Order: scur.OrderAsc, Value: "foobar"},
				scur.CursorParameter{Name: "Int", Order: scur.OrderAsc, Value: 1},
				scur.CursorParameter{Name: "Time", Order: scur.OrderAsc, Value: time.Now().Round(0)},
			},
		},
		{
			name: "5 columns",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "String", Order: scur.OrderAsc, Value: "foobar"},
				scur.CursorParameter{Name: "Int", Order: scur.OrderAsc, Value: 1},
				scur.CursorParameter{Name: "Time", Order: scur.OrderAsc, Value: time.Now().Round(0)},
				scur.CursorParameter{Name: "Bool", Order: scur.OrderAsc, Value: true},
				scur.CursorParameter{Name: "Float", Order: scur.OrderAsc, Value: 1.25},
			},
		},
		{
			name: "all types",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "String", Order: scur.OrderAsc, Value: "foobar"},
				scur.CursorParameter{Name: "Int", Order: scur.OrderAsc, Value: 1},
				scur.CursorParameter{Name: "Int64", Order: scur.OrderAsc, Value: int64(1)},
				scur.CursorParameter{Name: "Bool", Order: scur.OrderAsc, Value: true},
				scur.CursorParameter{Name: "Float64", Order: scur.OrderAsc, Value: 1.25},
				scur.CursorParameter{Name: "Time", Order: scur.OrderAsc, Value: time.Now().Round(0)},
			},
		},
		{
			name: "complex string",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "String", Order: scur.OrderAsc, Value: "https://godoc.org/cloud.google.com/go/spanner :::"},
			},
		},
		{
			name: "custom type value",
			cursor: scur.Cursor{
				scur.CursorParameter{Name: "CustomType", Order: scur.OrderAsc, Value: FooID("aaa")},
			},
			expected: scur.Cursor{
				scur.CursorParameter{Name: "CustomType", Order: scur.OrderAsc, Value: "aaa"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.cursor.EncodeParameters()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("encoded: %d %s", len(s), s)

			var actual scur.Cursor
			for _, cc := range tt.cursor {
				actual = append(actual, scur.CursorParameter{
					Name:  cc.Name,
					Order: cc.Order,
				})
			}

			err = scur.DecodeCursorParameters(actual, s)
			if err != nil {
				t.Fatal(err)
			}

			expected := tt.expected
			if expected == nil {
				expected = tt.cursor
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("unexpected expected = %v, decoded %v", expected, actual)
			}
		})
	}
}

func TestCursor_WhereExpression(t *testing.T) {
	tests := []struct {
		name   string
		cursor scur.Cursor
		sql    string
		params map[string]interface{}
	}{
		{
			"1 column",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderAsc, Value: 1},
			},
			heredoc.Doc(`
				(
				  A > @cursor1
				)
			`),
			map[string]interface{}{
				"cursor1": 1,
			},
		},
		{
			"2 columns",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderAsc, Value: 1},
				scur.CursorParameter{Name: "B", Order: scur.OrderDesc, Value: "b"},
			},
			heredoc.Doc(`
				(
				  ( A > @cursor1 )
				  OR
				  ( A = @cursor1 AND B < @cursor2 )
				)
			`),
			map[string]interface{}{
				"cursor1": 1,
				"cursor2": "b",
			},
		},
		{
			"3 columns",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderDesc, Value: 1},
				scur.CursorParameter{Name: "B", Order: scur.OrderAsc, Value: "b"},
				scur.CursorParameter{Name: "C", Order: scur.OrderAsc, Value: 1.25},
			},
			heredoc.Doc(`
				(
				  ( A < @cursor1 )
				  OR
				  ( A = @cursor1 AND B > @cursor2 )
				  OR
				  ( A = @cursor1 AND B = @cursor2 AND C > @cursor3 )
				)
			`),
			map[string]interface{}{
				"cursor1": 1,
				"cursor2": "b",
				"cursor3": 1.25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, params, err := tt.cursor.WhereExpression()
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.sql {
				t.Errorf("WhereExpression() sql = %v, want %v", sql, tt.sql)
			}
			if !reflect.DeepEqual(params, tt.params) {
				t.Errorf("WhereExpression() params = %v, want %v", params, tt.params)
			}
		})
	}
}

func TestCursor_OrderByExpression(t *testing.T) {
	tests := []struct {
		name   string
		cursor scur.Cursor
		sql    string
	}{
		{
			"1 column",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderAsc, Value: 1},
			},
			" A ASC",
		},
		{
			"2 columns",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderAsc, Value: 1},
				scur.CursorParameter{Name: "B", Order: scur.OrderDesc, Value: "b"},
			},
			" A ASC, B DESC",
		},
		{
			"3 columns",
			scur.Cursor{
				scur.CursorParameter{Name: "A", Order: scur.OrderDesc, Value: 1},
				scur.CursorParameter{Name: "B", Order: scur.OrderAsc, Value: "b"},
				scur.CursorParameter{Name: "C", Order: scur.OrderAsc, Value: 1.25},
			},
			" A DESC, B ASC, C ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, err := tt.cursor.OrderByExpression()
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.sql {
				t.Errorf("OrderByExpression() sql = %v, want %v", sql, tt.sql)
			}
		})
	}
}
