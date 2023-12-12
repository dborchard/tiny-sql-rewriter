package domain

// Meta With 'database' as key, DB map, metadata organized by db->table->column
type Meta map[string]*DB

// DB database related structures
type DB struct {
	Name  string
	Table map[string]*Table // ['table_name']*TableName
}

type Table struct {
	TableName    string
	TableAliases []string
	Column       map[string]*Column
}

// Column Contains column definition properties
type Column struct {
	Name        string   `json:"col_name"`    // 列名
	Alias       []string `json:"alias"`       // 别名
	Table       string   `json:"tb_name"`     // 表名
	DB          string   `json:"db_name"`     // 数据库名称
	DataType    string   `json:"data_type"`   // 数据类型
	Character   string   `json:"character"`   // 字符集
	Collation   string   `json:"collation"`   // collation
	Cardinality float64  `json:"cardinality"` // 散粒度
	Null        string   `json:"null"`        // 是否为空: YES/NO
	Key         string   `json:"key"`         // 键类型
	Default     string   `json:"default"`     // 默认值
	Extra       string   `json:"extra"`       // 其他
	Comment     string   `json:"comment"`     // 备注
	Privileges  string   `json:"privileges"`  // 权限
}

// TableColumns The elements in this structure are ordered map[db]map[table][]columns
type TableColumns map[string]map[string][]*Column

func NewDB(db string) *DB {
	return &DB{
		Name:  db,
		Table: make(map[string]*Table),
	}
}

func NewTable(tb string) *Table {
	return &Table{
		TableName:    tb,
		TableAliases: make([]string, 0),
		Column:       make(map[string]*Column),
	}
}
