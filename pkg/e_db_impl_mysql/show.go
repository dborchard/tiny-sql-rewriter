package database

import "fmt"

// desc table
// https://dev.mysql.com/doc/refman/5.7/en/show-columns.html

// TableDesc show columns from rental;
type TableDesc struct {
	Name       string
	DescValues []TableDescValue
}

// TableDescValue 含有每一列的属性
type TableDescValue struct {
	Field      string // List
	Type       string // type of data
	Collation  []byte // character set
	Null       string // whether there is null no yes
	Key        string // key type
	Default    []byte // default value
	Extra      string // Else
	Privileges string // Permissions
	Comment    string // Remark
}

// NewTableDesc initialize a*TableDesc
func NewTableDesc(tableName string) *TableDesc {
	return &TableDesc{
		Name:       tableName,
		DescValues: make([]TableDescValue, 0),
	}
}

// ShowColumns 获取 DB 中所有的 columns
func (db *Connector) ShowColumns(tableName string) (*TableDesc, error) {
	tbDesc := NewTableDesc(tableName)

	// 执行 show create table
	res, err := db.Query(fmt.Sprintf("show full columns from `%s`.`%s`", Escape(db.Database, false), Escape(tableName, false)))
	if err != nil {
		return nil, err
	}

	// columns info
	tc := TableDescValue{}
	columnFields := make([]interface{}, 0)
	fields := map[string]interface{}{
		"Field":      &tc.Field,
		"Type":       &tc.Type,
		"Collation":  &tc.Collation,
		"Null":       &tc.Null,
		"Key":        &tc.Key,
		"Default":    &tc.Default,
		"Extra":      &tc.Extra,
		"Privileges": &tc.Privileges,
		"Comment":    &tc.Comment,
	}
	cols, err := res.Rows.Columns()
	if err != nil {
		panic(err)
	}
	var colByPass []byte
	for _, col := range cols {
		if _, ok := fields[col]; ok {
			columnFields = append(columnFields, fields[col])
		} else {
			columnFields = append(columnFields, &colByPass)
		}
	}
	// get value
	for res.Rows.Next() {
		err := res.Rows.Scan(columnFields...)
		if err != nil {
			panic(err)
		}
		tbDesc.DescValues = append(tbDesc.DescValues, tc)
	}
	res.Rows.Close()
	return tbDesc, err
}

// Escape like C API mysql_escape_string()
func Escape(source string, NoBackslashEscapes bool) string {
	// NoBackslashEscapes https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_no_backslash_escapes
	// TODO: NoBackslashEscapes always false
	if NoBackslashEscapes {
		return quoteEscape(source)
	}
	return stringEscape(source)
}
