package env

import (
	"strings"
	domain "tiny_rewriter/pkg/c_domain"
	database "tiny_rewriter/pkg/e_catalog_mysql"
)

// VirtualEnv SQL optimization review test environment
// The information used by DB is obtained from the configuration file
type VirtualEnv struct {
	*database.Connector

	//Save the DB test environment mapping relationship to prevent vEnv environment conflicts.
	DBRef   map[string]string // db -> optimizer_xxx
	Hash2DB map[string]string // optimizer_xxx -> db
	//Save the Table creation relationship to prevent repeated table creation
	TableMap map[string]map[string]string
	// error
	Error error
}

// BuildEnv test environment initialization & connection online environment check
// @output *VirtualEnv test environment
// @output *y_database.Connector online environment connection handle
func BuildEnv() (*VirtualEnv, *database.Connector) {
	connTest, _ := database.NewConnector("")
	vEnv := NewVirtualEnv(connTest)
	connOnline, _ := database.NewConnector("")
	return vEnv, connOnline
}

// NewVirtualEnv initialize a new test environment
func NewVirtualEnv(vEnv *database.Connector) *VirtualEnv {
	return &VirtualEnv{
		Connector: vEnv,
		DBRef:     make(map[string]string),
		Hash2DB:   make(map[string]string),
		TableMap:  make(map[string]map[string]string),
	}
}

// GenTableColumns Initialization of the structure provided for Rewrite
func (vEnv *VirtualEnv) GenTableColumns(meta domain.Meta) domain.TableColumns {
	tableColumns := make(domain.TableColumns)
	for dbName, db := range meta {
		for _, tb := range db.Table {
			// prevent unexpected values from being passed in
			if tb == nil {
				break
			}
			td, err := vEnv.Connector.ShowColumns(tb.TableName)
			if err != nil {
				panic(err)
			}

			// tableColumns Initialize
			if dbName == "" {
				dbName = vEnv.RealDB(vEnv.Connector.Database)
			}

			if _, ok := tableColumns[dbName]; !ok {
				tableColumns[dbName] = make(map[string][]*domain.Column)
			}

			if _, ok := tableColumns[dbName][tb.TableName]; !ok {
				tableColumns[dbName][tb.TableName] = make([]*domain.Column, 0)
			}

			if len(tb.Column) == 0 {
				// tb.column. If it is empty, it means that this table is queried using * in SQL.
				if err != nil {
					panic(err)
				}

				for _, colInfo := range td.DescValues {
					tableColumns[dbName][tb.TableName] = append(tableColumns[dbName][tb.TableName], &domain.Column{
						Name:       colInfo.Field,
						DB:         dbName,
						Table:      tb.TableName,
						DataType:   colInfo.Type,
						Character:  string(colInfo.Collation),
						Key:        colInfo.Key,
						Default:    string(colInfo.Default),
						Extra:      colInfo.Extra,
						Comment:    colInfo.Comment,
						Privileges: colInfo.Privileges,
						Null:       colInfo.Null,
					})
				}
			} else {
				// tb.column If it is not empty, you need to fill in the columns used.
				var columns []*domain.Column
				for _, col := range tb.Column {
					for _, colInfo := range td.DescValues {
						if col.Name == colInfo.Field {
							// Complete the column information based on the obtained information
							col.DB = dbName
							col.Table = tb.TableName
							col.DataType = colInfo.Type
							col.Character = string(colInfo.Collation)
							col.Key = colInfo.Key
							col.Default = string(colInfo.Default)
							col.Extra = colInfo.Extra
							col.Comment = colInfo.Comment
							col.Privileges = colInfo.Privileges
							col.Null = colInfo.Null

							columns = append(columns, col)
							break
						}
					}
				}
				tableColumns[dbName][tb.TableName] = columns
			}
		}
	}
	return tableColumns
}

// RealDB get the hashed db from the test environment
func (vEnv *VirtualEnv) RealDB(hash string) string {
	if _, ok := vEnv.Hash2DB[hash]; ok {
		return vEnv.Hash2DB[hash]
	}
	// hash may be real database name not hash
	if strings.HasPrefix(hash, "optimizer_") {
		panic("hash is not optimizer_xxx")
	}
	return hash
}

func (vEnv *VirtualEnv) GenTableColumnsMock(meta domain.Meta) domain.TableColumns {
	dbName := "a"
	tblName := "tbl"

	tableColumns := make(domain.TableColumns)
	if _, ok := tableColumns[dbName]; !ok {
		tableColumns[dbName] = make(map[string][]*domain.Column)
	}

	var columns []*domain.Column
	columns = append(columns, &domain.Column{
		Name:       "id",
		DB:         dbName,
		Table:      tblName,
		DataType:   "varchar",
		Character:  "utf8",
		Key:        "PRI",
		Default:    "",
		Extra:      "",
		Comment:    "",
		Privileges: "",
		Null:       "",
	})
	columns = append(columns, &domain.Column{
		Name:       "name",
		DB:         dbName,
		Table:      tblName,
		DataType:   "int",
		Character:  "utf8",
		Default:    "",
		Extra:      "",
		Comment:    "",
		Privileges: "",
		Null:       "",
	})

	tableColumns[dbName][tblName] = columns

	return tableColumns
}
