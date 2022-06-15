package schema

import (
	"gee_orm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

type ITableName interface {
	TableName() string
}

func (schema *Schema) GetField(name string) *Field {
	return schema.FieldMap[name]
}

func (schema *Schema) RecordValues(dest interface{}) (fieldValues []interface{}) {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		tmp := modelType.Field(i)
		if !tmp.Anonymous && ast.IsExported(tmp.Name) {
			field := &Field{
				Name: tmp.Name,
				Type: d.DataTypeof(reflect.Indirect(reflect.New(tmp.Type))),
			}
			if v, ok := tmp.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, tmp.Name)
			schema.FieldMap[tmp.Name] = field
		}
	}
	return schema
}
