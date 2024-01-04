package main

import (
	"fmt"
	"strings"
	advisor "tiny_rewriter/pkg/a_advisor"
	rewrite "tiny_rewriter/pkg/b_rewriter"
	env "tiny_rewriter/pkg/d_catalog"
)

func main() {
	sql := "select * from tbl t1 where id < 1000"

	// 1. suggest optimizations
	heuristicSuggest := suggestOptimizations(sql)
	fmt.Println(heuristicSuggest)
	/*
		map[
			ALI.001: {ALI.001 L0 It is recommended to use the AS keyword to explicitly declare an alias In column or table aliases (such as "tbl AS alias"), explicit use of the AS keyword is more understandable than implicit aliases (such as "tbl alias").  select name from tbl t1 where id < 1000 0 0x10327f9e0}
			OK: {OK L0 OK OK OK 0 0x10327f920}
		]
	*/

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
