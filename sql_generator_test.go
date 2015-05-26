package toygo

import "testing"

func TestGenerateInsertCommand(t *testing.T) {
	except := "INSERT INTO `table` (`name`, `name2`) VALUES (?, ?)"
	f := &toyField{Name: "name"}
	f2 := &toyField{Name: "name2"}
	model := &toyModel{[]*toyField{f, f2}, "table"}
	result := generateInsertCommand(model)

	AssertEqual(t, except, result)
}
