package toygo

import (
	"fmt"
	"strings"
)

func generateInsertCommand(ms ...*Model) string {
	var fieldNames []string
	var insertValues []string
	var batchInsert []string

	m := ms[0]
	for _, field := range m.Fields {
		insertValues = append(insertValues, argumentTemplate)
		fieldNames = append(fieldNames, field.Name)
	}

	var inserts = "(" + strings.Join(insertValues, ", ") + ")"
	for range ms {
		batchInsert = append(batchInsert, inserts)
	}

	return fmt.Sprint(insertTemplate, m.TableName, strings.Join(batchInsert, ", "))
}
