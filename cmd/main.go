package main

import (
	"tiny_rewriter/advisor"
)

func main() {
	sql := "select * from t where id = 1"
	q, syntaxErr := advisor.NewQuery4Audit(sql)
	if syntaxErr != nil {
		panic(syntaxErr)
	}

	heuristicSuggest := make(map[string]advisor.Rule)

	for item, rule := range advisor.HeuristicRules {
		okFunc := (*advisor.Query4Audit).RuleOK

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

	//rw := ast.NewRewrite(sql)
	//meta := ast.GetMeta(rw.Stmt, nil)
	//rw.Columns = vEnv.GenTableColumns(meta)
	//// 执行定义好的 SQL 重写规则
	//rw.Rewrite()
	//fmt.Println(strings.TrimSpace(rw.NewSQL))

}
