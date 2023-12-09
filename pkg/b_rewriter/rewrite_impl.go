package rewrite

import (
	"vitess.io/vitess/go/vt/sqlparser"
)

// RewriteStar2Columns star2columns: corresponding to COL.001, SELECT completion * refers to the column name
func (rw *Rewrite) RewriteStar2Columns() *Rewrite {

	// Single table select * does not complete the table name to avoid the SQL being too long.
	// Select tb1.*, tb2.* of multiple tables needs to complete the table name.
	var multiTable bool
	if len(rw.Columns) > 1 {
		multiTable = true
	} else {
		for db := range rw.Columns {
			if len(rw.Columns[db]) > 1 {
				multiTable = true
			}
		}
	}

	err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.Select:

			// select possible situations
			// 1. select * from tb;
			// 2. select * from tb1,tb2;
			// 3. select tb1.* from tb1;
			// 4. select tb1.*,tb2.col from tb1,tb2;
			// 5. select db.tb1.* from tb1;
			// 6. select db.tb1.*,db.tb2.col from db.tb1,db.tb2;

			newSelectExprs := make(sqlparser.SelectExprs, 0)
			for _, expr := range n.SelectExprs {
				switch e := expr.(type) {
				case *sqlparser.StarExpr:
					// Under normal circumstances, the outermost loop will not exceed two layers.
					for _, tables := range rw.Columns {
						for _, cols := range tables {
							for _, col := range cols {
								var table string
								if multiTable {
									table = col.Table
								}
								newExpr := &sqlparser.AliasedExpr{
									Expr: &sqlparser.ColName{
										Metadata: nil,
										Name:     sqlparser.NewColIdent(col.Name),
										Qualifier: sqlparser.TableName{
											Name: sqlparser.NewTableIdent(table),
											// Because cross-DB queries are not recommended,
											// the db prefix here will not be completed.
											Qualifier: sqlparser.TableIdent{},
										},
									},
									As: sqlparser.ColIdent{},
								}

								if e.TableName.Name.IsEmpty() {
									// CASE 1 2
									newSelectExprs = append(newSelectExprs, newExpr)
								} else {
									// In other cases, replacement will only occur when the table name matches
									if e.TableName.Name.String() == col.Table {
										newSelectExprs = append(newSelectExprs, newExpr)
									}
								}
							}
						}
					}
				default:
					newSelectExprs = append(newSelectExprs, e)
				}
			}

			n.SelectExprs = newSelectExprs
		}
		return true, nil
	}, rw.Stmt)
	if err != nil {
		panic(err)
	}
	rw.NewSQL = sqlparser.String(rw.Stmt)
	return rw
}

// RewriteAddOrderByNull corresponding to CLA.008, add ORDER BY NULL when GROUP BY has no sorting requirements
func (rw *Rewrite) RewriteAddOrderByNull() *Rewrite {
	err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.Select:
			if n.GroupBy != nil && n.OrderBy == nil {
				n.OrderBy = sqlparser.OrderBy{
					&sqlparser.Order{
						Expr:      &sqlparser.NullVal{},
						Direction: "asc",
					},
				}
			}
		}
		return true, nil
	}, rw.Stmt)
	if err != nil {
		panic(err)
	}
	rw.NewSQL = sqlparser.String(rw.Stmt)
	return rw
}

// RewriteRemoveDMLOrderBy Corresponds to RES.004, deletes the ORDER BY contained in UPDATE and DELETE
// when there is no LIMIT condition
func (rw *Rewrite) RewriteRemoveDMLOrderBy() *Rewrite {
	switch st := rw.Stmt.(type) {
	case *sqlparser.Update:
		err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			switch n := node.(type) {
			case *sqlparser.Select:
				if n.OrderBy != nil && n.Limit == nil {
					n.OrderBy = nil
				}
				return false, nil
			}
			return true, nil
		}, rw.Stmt)
		if err != nil {
			panic(err)
		}
		if st.OrderBy != nil && st.Limit == nil {
			st.OrderBy = nil
		}
	case *sqlparser.Delete:
		err := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			switch n := node.(type) {
			case *sqlparser.Select:
				if n.OrderBy != nil && n.Limit == nil {
					n.OrderBy = nil
				}
				return false, nil
			}
			return true, nil
		}, rw.Stmt)
		if err != nil {
			panic(err)
		}
		if st.OrderBy != nil && st.Limit == nil {
			st.OrderBy = nil
		}
	}
	rw.NewSQL = sqlparser.String(rw.Stmt)
	return rw
}
