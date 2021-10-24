package sqb

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		name      string
		builder   func() Builder
		expected  string
		wantError string
	}{
		{
			name: "simple SELECT",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().C("1")
				return qb
			},
			expected: `SELECT 1`,
		},
		{
			name: "simple SELECT chain",
			builder: func() Builder {
				qb := NewBuilder()
				qb.
					Select().
					C("ID").C("CreatedAt", "AT").
					CS("A", "B", "C").
					From().
					Name("KeyData").
					Where().
					E("Disabled=@disabled").
					E("PEMData", "NOT NULL").
					OrderBy().O("CreatedAt DESC").O("KeyDataID", "ASC").
					Limit("@limit")
				return qb
			},
			expected: `SELECT ID, CreatedAt AT, A, B, C FROM KeyData WHERE Disabled=@disabled AND PEMData NOT NULL ORDER BY CreatedAt DESC, KeyDataID ASC LIMIT @limit`,
		},
		{
			name: "simple SELECT separate",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().C("ID").C("CreatedAt", "AT")
				qb.From().Name("KeyData")
				qb.Where().E("Disabled=@disabled")
				qb.Where().E("PEMData", "NOT NULL")
				qb.OrderBy().O("CreatedAt DESC").O("KeyDataID", "ASC")
				qb.Limit("@limit")
				return qb
			},
			expected: `SELECT ID, CreatedAt AT FROM KeyData WHERE Disabled=@disabled AND PEMData NOT NULL ORDER BY CreatedAt DESC, KeyDataID ASC LIMIT @limit`,
		},
		{
			name: "simple SELECT AS STRUCT",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().AsStruct().C("*")
				return qb
			},
			expected: `SELECT AS STRUCT *`,
		},
		{
			name: "simple DELETE",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Delete().From().Name("FOO")
				return qb
			},
			expected: `DELETE FROM FOO`,
		},
		{
			name: "simple DELETE chain",
			builder: func() Builder {
				qb := NewBuilder()
				qb.
					Delete().
					From().
					Name("KeyData").
					Where().
					E("Disabled=@disabled").
					E("PEMData", "NOT NULL")
				return qb
			},
			expected: `DELETE FROM KeyData WHERE Disabled=@disabled AND PEMData NOT NULL`,
		},
		{
			name: "error",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().C("ID", "AT", "UNKNOWN")
				qb.Limit("1")
				qb.Limit("2")
				return qb
			},
			wantError: "2 error(s) occured! SELECT ID AT !ERR1:`too many arguments: [AT UNKNOWN]`! LIMIT 1 !ERR2:`unexpected LIMIT keyword`! 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.builder().Build()
			if tt.wantError != "" {
				if err == nil {
					t.Fatal("unexpected")
				} else if v1, v2 := err.Error(), tt.wantError; v1 != v2 {
					t.Fatalf("unexpected: %v %v", v1, v2)
				}
				return
			} else if err != nil {
				t.Fatal(err)
			}
			if s != tt.expected {
				t.Fatalf("unexpected: %v %v", s, tt.expected)
			}
		})
	}
}
