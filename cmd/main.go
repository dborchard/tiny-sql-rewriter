package main

import (
	"fmt"
	"strings"
	advisor "tiny_rewriter/pkg/a_advisor"
	rewrite "tiny_rewriter/pkg/b_rewriter"
	env "tiny_rewriter/pkg/d_catalog"
)

func main() {
	sql := "select * from t where id = 1"

	// 1. suggest optimizations
	heuristicSuggest := suggestOptimizations(sql)
	fmt.Println(heuristicSuggest)

	// 2. rewrite sql
	newSql := rewriteSql(sql)
	fmt.Println(newSql)
}

func suggestOptimizations(sql string) map[string]advisor.Rule {
	q, _ := advisor.NewQuery4Audit(sql)
	return q.Advise()
}

func rewriteSql(sql string) string {
	//Environment initialization, connection check online environment + build test environment
	vEnv, _ := env.BuildEnv()

	rw := rewrite.NewRewrite(sql)
	rw.Columns = vEnv.GenTableColumns(rewrite.GetMeta(rw.Stmt, nil))
	rw.Rewrite()
	return strings.TrimSpace(rw.NewSQL)
}
