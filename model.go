package toygo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type toyField struct {
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

type toyModel struct {
	Fields    []*toyField
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

func generateField(structField reflect.StructField) *toyField {
	tags := structField.Tag

	f := new(toyField)
	f.Tags = tags
	f.PrimaryKey = getPrimaryKey(tags)
	f.Name = getFieldName(tags, structField.Name)
	f.Length = getFieldLength(tags)
	f.AutoIncrement = getAutoIncrement(tags)
	f.Nullable = getNullable(tags)
	f.Type = getType(tags, structField.Type)
	f.Default = getDefault(tags, f.Type)

	return f
}

func getValue(f *toyField, v reflect.Value) (value interface{}, err error) {
	if v == reflect.Zero(f.Type) {
		if f.Default == reflect.Zero(f.Type) {
			if !f.Nullable {
				return nil, &Error{fmt.Sprintf("%s could not be nil", f.Name)}
			}
			return f.Default, nil
		}
	}
	return v, nil
}

func getTableName(model *toyModel) (name string) {
	for _, f := range model.Fields {
		if f.PrimaryKey {
			if f.Tags.Get(FIELD_KEY_TABLE_NAME) != "" {
				return f.Tags.Get(FIELD_KEY_TABLE_NAME)
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

func getToyModel(entity interface{}) (model *toyModel, err error) {
	st := reflect.TypeOf(entity)

	fieldNum := st.NumField()

	var fields []*toyField
	for i := 0; i < fieldNum; i++ {
		f := generateField(st.Field(i))

		value := reflect.ValueOf(entity)
		f.Value, err = getValue(f, reflect.Indirect(value).Field(i))
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}

	model = &toyModel{fields, st.Name()}
	model.TableName = getTableName(model)

	return model, err
}

func (s *Session) Save(entity interface{}) (num int64, err error) {
	model, err := getToyModel(entity)
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec(generateInsertCommand(model))
	return result.RowsAffected()
}
