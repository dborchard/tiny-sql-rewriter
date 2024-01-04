package advisor

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	tidb "github.com/pingcap/parser/ast"
	"vitess.io/vitess/go/vt/sqlparser"
)
import _ "github.com/pingcap/tidb/types/parser_driver"

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
	q.TiStmt, err = TiParse(sql, charset, collation)
	return q, err
}

func (q *Query4Audit) Advise() map[string]Rule {
	heuristicSuggest := make(map[string]Rule)

	for item, rule := range HeuristicRules {
		okFunc := (*Query4Audit).RuleOK

		heuristicFn := rule.Func
		if &heuristicFn != &okFunc {
			// NOTE: This is the key point of this snippet.
			//By calling r := heuristicFn(q), you're invoking the function with q as its receiver.
			//This means that if heuristicFn is set to (*Query4Audit).RuleImplicitAlias, then heuristicFn(q) is equivalent to q.RuleImplicitAlias().
			r := heuristicFn(q)

			if r.Item == item {
				heuristicSuggest[item] = r
			}
		}
	}
	return heuristicSuggest
}

// TiParse TiDB grammar analysis
func TiParse(sql, charset, collation string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmt, _, err := p.Parse(sql, charset, collation)
	return stmt, err
}
