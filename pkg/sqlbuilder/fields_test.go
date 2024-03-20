package sqlbuilder

import (
	"github.com/mradmacher/audiofeeler/pkg/optiomist"
	"reflect"
	"testing"
)

type PossibleResult struct {
	sql    string
	values []any
}

func TestFields_BuildInsert(t *testing.T) {
	tests := []struct {
		fields          Fields
		possibleResults []PossibleResult
	}{
		{
			fields: Fields{
				"name": optiomist.Some("Jan Nowak"),
				"age":  optiomist.Some(21),
			},
			possibleResults: []PossibleResult{{
				sql: "INSERT INTO users (name, age) " +
					"VALUES ($1, $2) RETURNING id",
				values: []any{"Jan Nowak", 21},
			}, {
				sql: "INSERT INTO users (age, name) " +
					"VALUES ($1, $2) RETURNING id",
				values: []any{21, "Jan Nowak"},
			}},
		}, {
			fields: Fields{
				"first_name": optiomist.Some("Jan"),
				"last_name":  optiomist.Some("Nowak"),
				"age":        optiomist.Nil[int](),
			},
			possibleResults: []PossibleResult{{
				sql: "INSERT INTO users (first_name, last_name, age) " +
					"VALUES ($1, $2, NULL) RETURNING id",
				values: []any{"Jan", "Nowak"},
			}, {
				sql: "INSERT INTO users (last_name, first_name, age) " +
					"VALUES ($1, $2, NULL) RETURNING id",
				values: []any{"Nowak", "Jan"},
			}, {
				sql: "INSERT INTO users (first_name, age, last_name) " +
					"VALUES ($1, NULL, $2) RETURNING id",
				values: []any{"Jan", "Nowak"},
			}, {
				sql: "INSERT INTO users (last_name, age, first_name) " +
					"VALUES ($1, NULL, $2) RETURNING id",
				values: []any{"Nowak", "Jan"},
			}, {
				sql: "INSERT INTO users (age, first_name, last_name) " +
					"VALUES (NULL, $1, $2) RETURNING id",
				values: []any{"Jan", "Nowak"},
			}, {
				sql: "INSERT INTO users (age, last_name, first_name) " +
					"VALUES (NULL, $1, $2) RETURNING id",
				values: []any{"Nowak", "Jan"},
			}},
		}, {
			fields: Fields{
				"name": optiomist.Nil[string](),
				"age":  optiomist.Nil[int](),
			},
			possibleResults: []PossibleResult{{
				sql: "INSERT INTO users (name, age) " +
					"VALUES (NULL, NULL) RETURNING id",
				values: []any{},
			}, {
				sql: "INSERT INTO users (age, name) " +
					"VALUES (NULL, NULL) RETURNING id",
				values: []any{},
			}},
		}, {
			fields: Fields{},
			possibleResults: []PossibleResult{{
				sql:    "INSERT INTO users DEFAULT VALUES RETURNING id",
				values: []any{},
			}},
		}, {
			fields: Fields{},
			possibleResults: []PossibleResult{{
				sql:    "INSERT INTO users DEFAULT VALUES RETURNING id",
				values: []any{},
			}},
		}, {
			fields: Fields{
				"first_name": optiomist.None[string](),
				"last_name":  optiomist.None[string](),
				"age":        optiomist.None[int](),
			},
			possibleResults: []PossibleResult{{
				sql:    "INSERT INTO users DEFAULT VALUES RETURNING id",
				values: []any{},
			}},
		},
	}

	for _, test := range tests {
		sql, values := test.fields.BuildInsert("users")
		foundSql := false
		foundValues := false
		for _, result := range test.possibleResults {
			if result.sql == sql {
				foundSql = true
				if reflect.DeepEqual(result.values, values) {
					foundValues = true
				}
				break
			}
		}
		if !foundSql {
			t.Errorf("SQL = %q; expected %q", sql, test.possibleResults[0].sql)
		}
		if !foundValues {
			t.Errorf("Values = %v; expected %v", values, test.possibleResults[0].values)
		}
	}
}
