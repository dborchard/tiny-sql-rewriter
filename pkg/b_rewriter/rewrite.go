package rewrite

import (
	common "tiny_rewriter/pkg/c_domain"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Rewrite 用于重写SQL
type Rewrite struct {
	SQL     string
	NewSQL  string
	Stmt    sqlparser.Statement
	Columns common.TableColumns
}

// NewRewrite Returns a *Rewrite object. If the SQL cannot be parsed normally,
// the error will be output to the log and a nil will be returned.
func NewRewrite(sql string) *Rewrite {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}

	return &Rewrite{
		SQL:  sql,
		Stmt: stmt,
	}
}

func (rw *Rewrite) Rewrite() *Rewrite {
	for _, rule := range RewriteRules {
		rewriteFn := rule.Func

		if rewriteFn != nil {
			rewriteFn(rw)
		}
	}
	if rw.NewSQL == "" {
		rw.NewSQL = rw.SQL
	}
	rw.Stmt, _ = sqlparser.Parse(rw.NewSQL)
	return rw
}
