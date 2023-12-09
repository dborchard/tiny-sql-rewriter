package advisor

type Rule struct {
	Item     string                  `json:"Item"`     // ruleCode
	Severity string                  `json:"Severity"` // Hazard level. L0: OK, L1: Warning, L2: Critical
	Summary  string                  `json:"Summary"`  // summary of rules
	Content  string                  `json:"Content"`  // rule explanation
	Case     string                  `json:"Case"`     // sql example
	Position int                     `json:"Position"` // The SQL character position of the suggestion. The default value is 0, indicating global suggestion.
	Func     func(*Query4Audit) Rule `json:"-"`        // function name
}

var HeuristicRules map[string]Rule

func init() {
	InitHeuristicRules()
}

// InitHeuristicRules Initialize the heuristic rules
func InitHeuristicRules() {
	HeuristicRules = map[string]Rule{
		"OK": {
			Item:     "OK",
			Severity: "L0",
			Summary:  "OK",
			Content:  `OK`,
			Case:     "OK",
			Func:     (*Query4Audit).RuleOK,
		},
		"ALI.001": {
			Item:     "ALI.001",
			Severity: "L0",
			Summary:  "It is recommended to use the AS keyword to explicitly declare an alias",
			Content:  `In column or table aliases (such as "tbl AS alias"), explicit use of the AS keyword is more understandable than implicit aliases (such as "tbl alias"). `,
			Case:     "select name from tbl t1 where id < 1000",
			Func:     (*Query4Audit).RuleImplicitAlias,
		},
	}
}
