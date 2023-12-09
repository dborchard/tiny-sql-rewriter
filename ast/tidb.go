package ast

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

// TiParse TiDB 语法解析
func TiParse(sql, charset, collation string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmt, _, err := p.Parse(sql, charset, collation)
	return stmt, err
}
