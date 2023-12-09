package advisor

import (
	"tiny_rewriter/ast"
	"vitess.io/vitess/go/vt/sqlparser"
)

// RuleOK OK
func (q *Query4Audit) RuleOK() Rule {
	return HeuristicRules["OK"]
}

// RuleImplicitAlias ALI.001
func (q *Query4Audit) RuleImplicitAlias() Rule {
	var rule = q.RuleOK()
	tkns := ast.Tokenizer(q.Query)
	if len(tkns) == 0 {
		return rule
	}
	if tkns[0].Type != sqlparser.SELECT {
		return rule
	}
	for i, tkn := range tkns {
		if tkn.Type == sqlparser.ID &&
			i+1 < len(tkns) && tkn.Type == tkns[i+1].Type {

			rule = HeuristicRules["ALI.001"]
			break
		}
	}
	return rule
}
