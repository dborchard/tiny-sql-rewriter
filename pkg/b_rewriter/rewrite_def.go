package rewrite

// Rule sqlRewriteRules
type Rule struct {
	Name        string                  `json:"Name"`
	Description string                  `json:"Description"`
	Original    string                  `json:"Original"` // Error demonstration. Empty or "not supported yet" will not appear in list-rewrite-rules
	Suggest     string                  `json:"Suggest"`  // demonstrate correctly。
	Func        func(*Rewrite) *Rewrite `json:"-"`        // If Func is not defined, multiple SQL statements need to be rewritten in conjunction
}

// RewriteRules SQL重写规则，注意这个规则是有序的，先后顺序不能乱
var RewriteRules []Rule

func init() {
	RewriteRules = []Rule{
		{
			Name:        "star2columns",
			Description: "complete table column information for select",
			Original:    "SELECT * FROM film",
			Suggest:     "select film.film_id, film.title from film",
			Func:        (*Rewrite).RewriteStar2Columns,
		},
		{
			Name:        "orderbynull",
			Description: "If the GROUP BY statement does not specify the ORDER BY condition, it will cause unnecessary sorting. If sorting is not required, it is recommended to add ORDER BY NULL.",
			Original:    "SELECT sum(col1) FROM tbl GROUP BY col",
			Suggest:     "select sum(col1) from tbl group by col order by null",
			Func:        (*Rewrite).RewriteAddOrderByNull,
		},
		{
			Name:        "dmlorderby",
			Description: "Remove meaningless ORDER BY in DML update operations",
			Original:    "DELETE FROM tbl WHERE col1=1 ORDER BY col",
			Suggest:     "delete from tbl where col1 = 1",
			Func:        (*Rewrite).RewriteRemoveDMLOrderBy,
		},
	}
}
