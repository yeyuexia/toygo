package toygo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Field struct {
	Name          string
	Length        int
	AutoIncrement bool
	Nullable      bool
	Updatable     bool
	PrimaryKey    bool
	Type          reflect.Type
	Default       interface{}
	Value         interface{}
	Tags          reflect.StructTag
}

type Model struct {
	Fields    []*Field
	TableName string
}

func getFieldName(tags reflect.StructTag, fieldName string) (name string) {
	if tags.Get(FIELD_KEY_NAME) == "" {
		return strings.ToLower(fieldName)
	} else {
		return tags.Get(FIELD_KEY_NAME)
	}
}

func getFieldLength(tags reflect.StructTag) (length int) {
	if tags.Get(FIELD_KEY_LENGTH) != "" {
		length, err := strconv.Atoi(tags.Get(FIELD_KEY_LENGTH))
		if err != nil {
			panic("length must be int")
		}
		return length
	}
	return 255
}

func getAutoIncrement(tags reflect.StructTag) (isAutoIncrement bool) {
	if tags.Get(FIELD_KEY_AUTO_INCR) != "" {
		isAutoIncrement, err := strconv.ParseBool(tags.Get(FIELD_KEY_AUTO_INCR))
		if err != nil {
			panic("auto_increment value type must be bool")
		}
		return isAutoIncrement
	}
	return false
}

func getNullable(tags reflect.StructTag) (nullable bool) {
	if tags.Get(FIELD_KEY_NULLABLE) != "" {
		nullable, err := strconv.ParseBool(tags.Get(FIELD_KEY_NULLABLE))
		if err != nil {
			panic("nullable value type must be bool")
		}
		return nullable
	}
	return true
}

func getType(tags reflect.StructTag, fieldType reflect.Type) (fType reflect.Type) {
	if tags.Get(FIELD_KEY_TYPE) != "" {
		// TODO: could get type by string
	}
	return fieldType
}

func getDefault(tags reflect.StructTag, fType reflect.Type) (value reflect.Value) {
	if tags.Get(FIELD_KEY_DEFAULT) != "" {
	}
	return reflect.Zero(fType)
}

func getPrimaryKey(tags reflect.StructTag) (result bool) {
	if tags.Get(FIELD_KEY_PK) != "" {
		result, err := strconv.ParseBool(tags.Get(FIELD_KEY_PK))
		if err != nil {
			panic("nullable value type must be bool")
		}
		return result
	}
	return false
}

func generateField(structField reflect.StructField) (field *Field) {
	tags := structField.Tag

	field.Tags = tags
	field.PrimaryKey = getPrimaryKey(tags)
	field.Name = getFieldName(tags, structField.Name)
	field.Length = getFieldLength(tags)
	field.AutoIncrement = getAutoIncrement(tags)
	field.Nullable = getNullable(tags)
	field.Type = getType(tags, structField.Type)
	field.Default = getDefault(tags, field.Type)

	return field
}

func getValue(field *Field, v reflect.Value) (value interface{}, err error) {
	if v == reflect.Zero(field.Type) {
		if field.Default == reflect.Zero(field.Type) {
			if !field.Nullable {
				return nil, &Error{fmt.Sprintf("%s could not be nil", field.Name)}
			}
			return field.Default, nil
		}
	}
	return v, nil
}

func getTableName(model *Model) (name string) {
	for _, field := range model.Fields {
		if field.PrimaryKey {
			if field.Tags.Get(FIELD_KEY_TABLE_NAME) != "" {
				return field.Tags.Get(FIELD_KEY_TABLE_NAME)
			}
		}
	}

	name = model.TableName
	for _, raneValue := range model.TableName {
		if unicode.IsUpper(raneValue) {
			// TODO: convert tableName to table_name
		}
	}
	return name
}

func getModel(entity interface{}) (model *Model, err error) {
	st := reflect.TypeOf(entity)

	fieldNum := st.NumField()

	var fields []*Field
	for i := 0; i < fieldNum; i++ {
		field := generateField(st.Field(i))

		value := reflect.ValueOf(entity)
		field.Value, err = getValue(field, reflect.Indirect(value).Field(i))
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	model = &Model{fields, st.Name()}
	model.TableName = getTableName(model)

	return model, err
}

func (s *Session) Save(entity interface{}) (num int64, err error) {
	model, err := getModel(entity)
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec(generateInsertCommand(model))
	return result.RowsAffected()
}
