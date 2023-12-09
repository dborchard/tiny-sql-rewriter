package ast

//// GetMeta 获取元数据信息，构建到db->table层级。
//// 从 SQL 或 Statement 中获取表信息，并返回。当 meta 不为 nil 时，返回值会将新老 meta 合并去重
//func GetMeta(stmt sqlparser.Statement, meta common.Meta) common.Meta {
//	// 初始化meta
//	if meta == nil {
//		meta = make(map[string]*common.DB)
//	}
//
//	err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
//		switch expr := node.(type) {
//		case *sqlparser.DDL:
//			// 如果 SQL 是一个 DDL，则不需要继续遍历语法树了
//			for _, tb := range expr.FromTables {
//				appendTable(tb, "", meta)
//			}
//
//			for _, tb := range expr.ToTables {
//				appendTable(tb, "", meta)
//			}
//
//			appendTable(expr.Table, "", meta)
//			return false, nil
//		case *sqlparser.AliasedTableExpr:
//			// 非 DDL 情况下处理 TableExpr
//			// 在 sqlparser 中存在三种 TableExpr: AliasedTableExpr，ParenTableExpr 以及 JoinTableExpr。
//			// 其中 AliasedTableExpr 是其他两种 TableExpr 的基础组成，SQL中的 表信息（别名、前缀）在这个结构体中。
//
//			switch table := expr.Expr.(type) {
//
//			// 获取表名、别名与前缀名（数据库名）
//			// 表名存放在 AST 中 TableName 里，包含表名与表前缀名。
//			// 当与 As 相对应的 Expr 为 TableName 的时候，别名才是一张实体表的别名，否则为结果集的别名。
//			case sqlparser.TableName:
//				appendTable(table, expr.As.String(), meta)
//			default:
//				// 如果 AliasedTableExpr 中的 Expr 不是 TableName 结构体，则表示该表为一个查询结果集（子查询或临时表）。
//				// 在这里记录一下别名，但将列名制空，用来保证在其他环节中判断列前缀的时候不会有遗漏
//				// 最终结果为所有的子查询别名都会归于 ""（空） 数据库 ""（空） 表下，对于空数据库，空表后续在索引优化时直接PASS
//				if meta == nil {
//					meta = make(map[string]*common.DB)
//				}
//
//				if meta[""] == nil {
//					meta[""] = common.NewDB("")
//				}
//
//				meta[""].Table[""] = common.NewTable("")
//				meta[""].Table[""].TableAliases = append(meta[""].Table[""].TableAliases, expr.As.String())
//			}
//		}
//		return true, nil
//	}, stmt)
//	common.LogIfWarn(err, "")
//	return meta
//}
