package toygo

import (
	"fmt"
	"strings"
)

func generateInsertCommand(ms ...*toyModel) string {
	var fieldNames []string
	var insertValues []string
	var batchInsert []string

	m := ms[0]
	for _, f := range m.Fields {
		insertValues = append(insertValues, argumentTemplate)
		fieldNames = append(fieldNames, "`"+f.Name+"`")
	}

	var inserts = "(" + strings.Join(insertValues, ", ") + ")"
	for range ms {
		batchInsert = append(batchInsert, inserts)
	}

	return fmt.Sprintf(insertTemplate, m.TableName, strings.Join(fieldNames, ", "), strings.Join(batchInsert, ", "))
}
