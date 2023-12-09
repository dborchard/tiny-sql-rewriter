package advisor

import (
	tidb "github.com/pingcap/parser/ast"
	"tiny_rewriter/ast"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Query4Audit The SQL structure to be reviewed consists of the original SQL and its corresponding abstract syntax tree
type Query4Audit struct {
	Query  string              // SQL
	Stmt   sqlparser.Statement // AST Parsed Through Vitess
	TiStmt []tidb.StmtNode     // AST Parsed Through TiDB
}

// NewQuery4Audit return a struct for Query4Audit
func NewQuery4Audit(sql string, options ...string) (*Query4Audit, error) {
	var err, vErr error
	var charset string
	var collation string

	if len(options) > 0 {
		charset = options[0]
	}

	if len(options) > 1 {
		collation = options[1]
	}

	q := &Query4Audit{Query: sql}

	// 1. vitess parser
	q.Stmt, vErr = sqlparser.Parse(sql)
	if vErr != nil {
		return nil, vErr
	}

	// tidb parser
	q.TiStmt, err = ast.TiParse(sql, charset, collation)
	return q, err
}
