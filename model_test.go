package toygo

import "testing"
import "reflect"

func TestGetFieldName(t *testing.T) {
	var tags reflect.StructTag
	tags = "name:\"abc\""
	result := getFieldName(tags, "tt")
	if result != "abc" {
		t.Fatalf("getFieldName %s, expected %s, got %s.", tags, "abc", result)
	}
}

func TestSuccessGetFieldNameWithMultiTag(t *testing.T) {
	var tags reflect.StructTag
	tags = "length:\"10\" name:\"abc\" "
	result := getFieldName(tags, "tt")
	if result != "abc" {
		t.Fatalf("getFieldName %s, expected %s, got %s.", tags, "abc", result)
	}
}

func TestUseDefaultFieldName(t *testing.T) {
	var tags reflect.StructTag
	tags = "length:\"100\""
	result := getFieldName(tags, "default")

	if result != "default" {
		t.Fatalf("getFieldName %s, expected %s, got %s.", tags, "default", result)
	}
}

func TestGetLength(t *testing.T) {
	var tags reflect.StructTag
	tags = "length:\"100\""
	result := getFieldLength(tags)

	if result != 100 {
		t.Fatalf("getFieldName %s, expected %s, got %s.", tags, 100, result)
	}
}

func TestGetAutoIncrement(t *testing.T) {
	var tags reflect.StructTag
	tags = "auto_increment:\"true\""
	result := getAutoIncrement(tags)

	if !result {
		t.Fatalf("getFieldName %s, expected %s, got %s.", tags, 100, result)
	}
}

func TestGetToyModel(t *testing.T) {
	type testModel struct {
		Id    int64  `id:"true" auto_increment:"true" table_name:"test_my_model"`
		Name  string `name:"name" length:"10"`
		Value int    `name:"value"`
	}

	test := testModel{Name: "name", Value: 10}
	result, err := getToyModel(test)
	if err != nil {
		t.Fatalf("testGetToyModel", err.Error())
	}

	if len(result.Fields) != 3 {
		t.Fatalf("testGetToyModel, except 3, got %d", result.Fields)
	} else if result.TableName != "test_my_model" {
		t.Fatalf("testGetToyModel, error")
	}
}
