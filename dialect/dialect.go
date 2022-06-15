package dialect

import (
	"reflect"
)

type Dialect interface {
	DataTypeof(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
