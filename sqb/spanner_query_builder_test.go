package sqb

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
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
			expected: strings.TrimSpace(heredoc.Doc(`
				SELECT
				  1
			`)),
		},
		{
			name: "simple SELECT wo indent",
			builder: func() Builder {
				qb := NewBuilder(NoIndent())
				qb.Select().C("1")
				return qb
			},
			expected: heredoc.Doc("SELECT 1"),
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
			expected: strings.TrimSpace(heredoc.Doc(`
				SELECT
				  ID,
				  CreatedAt AT,
				  A,
				  B,
				  C
				FROM
				  KeyData
				WHERE
				  Disabled=@disabled
				  AND PEMData NOT NULL
				ORDER BY
				  CreatedAt DESC,
				  KeyDataID ASC
				LIMIT @limit
			`)),
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
			expected: strings.TrimSpace(heredoc.Doc(`
				SELECT
				  ID,
				  CreatedAt AT
				FROM
				  KeyData
				WHERE
				  Disabled=@disabled
				  AND PEMData NOT NULL
				ORDER BY
				  CreatedAt DESC,
				  KeyDataID ASC
				LIMIT @limit
			`)),
		},
		{
			name: "simple SELECT DISTINCT",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().Distinct().C("Foo")
				return qb
			},
			expected: strings.TrimSpace(heredoc.Doc(`
				SELECT DISTINCT
				  Foo
			`)),
		},
		{
			name: "simple SELECT AS STRUCT",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Select().AsStruct().C("*")
				return qb
			},
			expected: strings.TrimSpace(heredoc.Doc(`
				SELECT AS STRUCT
				  *
			`)),
		},
		{
			name: "simple UPDATE",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Update("FOO").Set().U("String", "=", "@s")
				return qb
			},
			expected: strings.TrimSpace(heredoc.Doc(`
				UPDATE
				  FOO
				SET
				  String = @s
			`)),
		},
		{
			name: "simple UPDATE chain",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Update("FOO").Set().U("String", "=", "@s").U("Int = @int")
				qb.Where().E("ID = @id").E("B = TRUE")
				return qb
			},
			expected: strings.TrimSpace(heredoc.Doc(`
				UPDATE
				  FOO
				SET
				  String = @s,
				  Int = @int
				WHERE
				  ID = @id
				  AND B = TRUE
			`)),
		},
		{
			name: "simple DELETE",
			builder: func() Builder {
				qb := NewBuilder()
				qb.Delete().From().Name("FOO")
				return qb
			},
			expected: strings.TrimSpace(heredoc.Doc(`
				DELETE
				FROM
				  FOO
			`)),
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
			expected: strings.TrimSpace(heredoc.Doc(`
				DELETE
				FROM
				  KeyData
				WHERE
				  Disabled=@disabled
				  AND PEMData NOT NULL
			`)),
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
			wantError: strings.TrimSpace(heredoc.Docf(`
				2 error(s) occured! SELECT
				  ID AT !ERR1:%[1]stoo many arguments: [AT UNKNOWN]%[1]s!
				LIMIT 1 !ERR2:%[1]sunexpected LIMIT keyword%[1]s! 2
			`, "`")),
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
				t.Fatalf("unexpected: %v. expected: %v", s, tt.expected)
			}
		})
	}
}
