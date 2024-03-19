package sqlbuilder

import (
	"fmt"
	"github.com/mradmacher/audiofeeler/optiomist"
	"strings"
)

type Fields map[string]optiomist.Optionable

func (fields *Fields) BuildInsert(tableName string) (string, []any) {
	// TODO:
	// names := make([]string, len(*fields))
	// params := make([]string, len(*fields))
	// values := make([]any, len(*fields))
	var names []string
	var params []string
	values := []any{}
	n := 1
	for name, value := range *fields {
		if value.IsSome() {
			names = append(names, name)
			if value.IsNil() {
				params = append(params, "NULL")
			} else {
				params = append(params, fmt.Sprintf("$%d", n))
				n++
				values = append(values, value.AnyValue())
			}
		}
	}
	query := "INSERT INTO " + tableName
	if len(names) > 0 {
		query = query + " (" +
			strings.Join(names, ", ") + ") VALUES (" +
			strings.Join(params, ", ") + ")"
	} else {
		query = query + " DEFAULT VALUES"
	}
	query = query + " RETURNING id"
	return query, values
}
