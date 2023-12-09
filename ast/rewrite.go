package ast

//// Rewrite 用于重写SQL
//type Rewrite struct {
//	SQL     string
//	NewSQL  string
//	Stmt    sqlparser.Statement
//	Columns common.TableColumns
//}
//
//// NewRewrite 返回一个*Rewrite对象，如果SQL无法被正常解析，将错误输出到日志中，返回一个nil
//func NewRewrite(sql string) *Rewrite {
//	stmt, err := sqlparser.Parse(sql)
//	if err != nil {
//		panic(err)
//	}
//
//	return &Rewrite{
//		SQL:  sql,
//		Stmt: stmt,
//	}
//}
