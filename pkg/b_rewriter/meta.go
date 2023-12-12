package rewrite

import (
	"tiny_rewriter/pkg/c_domain"
	"vitess.io/vitess/go/vt/sqlparser"
)

// GetMeta gets metadata information and builds it to the db->table level.
// Get table information from SQL or Statement and return it. When meta is not nil,
// the return value will merge the old and new meta to remove duplicates.
func GetMeta(stmt sqlparser.Statement, meta domain.Meta) domain.Meta {
	// 初始化meta
	if meta == nil {
		meta = make(map[string]*domain.DB)
	}

	err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch expr := node.(type) {
		case *sqlparser.DDL:
			//If SQL is a DDL, there is no need to continue traversing the syntax tree
			for _, tb := range expr.FromTables {
				appendTable(tb, "", meta)
			}

			for _, tb := range expr.ToTables {
				appendTable(tb, "", meta)
			}

			appendTable(expr.Table, "", meta)
			return false, nil
		case *sqlparser.AliasedTableExpr:
			// Process TableExpr in non-DDL situations
			// There are three types of TableExpr in sqlparser:
			// AliasedTableExpr, ParenTableExpr and JoinTableExpr.
			// Among them, AliasedTableExpr is the basic component of the other two TableExprs.
			// The table information (alias, prefix) in SQL is in this structure.

			switch table := expr.Expr.(type) {

			// Get the table name, alias and prefix name (database name)
			// The table name is stored in TableName in the AST, including the table name and table prefix.
			// When the Expr corresponding to As is TableName, the alias is the alias of an entity table,
			// otherwise it is the alias of the result set.
			case sqlparser.TableName:
				appendTable(table, expr.As.String(), meta)
			default:
				// If Expr in AliasedTableExpr is not a TableName structure, it means that the table is a query result set (subquery or temporary table).
				// Record the alias here, but make the column name blank to ensure that there are no omissions when judging column prefixes in other links.
				// The final result is that all subquery aliases will be attributed to the "" (empty) database "" (empty) table.
				// For an empty database, the empty table will be directly PASSed during index optimization.
				if meta == nil {
					meta = make(map[string]*domain.DB)
				}

				if meta[""] == nil {
					meta[""] = domain.NewDB("")
				}

				meta[""].Table[""] = domain.NewTable("")
				meta[""].Table[""].TableAliases = append(meta[""].Table[""].TableAliases, expr.As.String())
			}
		}
		return true, nil
	}, stmt)
	if err != nil {
		panic(err)
	}
	return meta
}

// appendTable extracts the database table information in sqlparser.TableName and puts it into meta
// @tb is the sqlparser.TableName object
// @as is the alias of the table, empty if there is no alias
// @meta is the information collection
func appendTable(tb sqlparser.TableName, as string, meta map[string]*domain.DB) map[string]*domain.DB {
	if meta == nil {
		return meta
	}

	dbName := tb.Qualifier.String()
	tbName := tb.Name.String()
	if tbName == "" {
		return meta
	}

	if meta[dbName] == nil {
		meta[dbName] = domain.NewDB(dbName)
	}

	meta[dbName].Table[tbName] = domain.NewTable(tbName)
	mergeAlias(dbName, tbName, as, meta)

	return meta
}

// mergeAlias merge all table aliases into one table
func mergeAlias(db, tb, as string, meta map[string]*domain.DB) {
	if meta == nil || as == "" {
		return
	}

	aliasExist := false
	for _, existedAlias := range meta[db].Table[tb].TableAliases {
		if existedAlias == as {
			aliasExist = true
		}
	}

	if !aliasExist {
		meta[db].Table[tb].TableAliases = append(meta[db].Table[tb].TableAliases, as)
	}
}
