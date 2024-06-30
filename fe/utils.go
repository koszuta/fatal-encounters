package fe

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// NewSQLString trims and wraps a string in sql.NullString; valid when the trimmed string is not blank
func NewSQLString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{String: s, Valid: len(s) != 0}
}

// StructDiff ...
func StructDiff(v, w any) string {
	v1, v2 := reflect.ValueOf(v), reflect.ValueOf(w)
	t1, t2 := v1.Type(), v2.Type()
	typeName := t1.Name()
	if t1.Kind() != reflect.Struct || t2.Kind() != reflect.Struct {
		return fmt.Sprintf("%s and %s must be structs", typeName, t2.Name())
	}
	if t1 != t2 {
		return fmt.Sprintf("type mismatch: %s and %s", typeName, t2.Name())
	}

	var b strings.Builder
	for i := 0; i < v1.NumField(); i++ {
		f1, f2 := v1.Field(i), v2.Field(i)
		if f1.CanInterface() && f2.CanInterface() && f1.Interface() != f2.Interface() {
			if b.Len() != 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("%s:%v->%v", t1.Field(i).Name, f1.Interface(), f2.Interface()))
		}
	}

	if b.Len() == 0 {
		if strings.HasSuffix(typeName, "s") || strings.HasSuffix(typeName, "x") || strings.HasSuffix(typeName, "ch") || strings.HasSuffix(typeName, "sh") {
			return fmt.Sprintf("equal %ses", typeName)
		}
		return fmt.Sprintf("equal %ss", typeName)
	}
	return fmt.Sprintf("%s{ %s }", typeName, b.String())
}
