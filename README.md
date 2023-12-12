## Tiny SQL Rewriter

This project is based out of [XiaoMi Soar](https://github.com/XiaoMi/soar.git)

### Components
- Advisor : This code is gives suggestions on how to improve your SQL query. It does not rewrite the query. The rules
 are taken from [here](https://github.com/XiaoMi/soar/blob/fab04633b12ba1e4f35456112360150a6d0d1421/advisor/rules.go#L119)
- Rewriter : This code rewrites the SQL query based on a set of Rewrite rules. The rules are taken from
 [here](https://github.com/XiaoMi/soar/blob/fab04633b12ba1e4f35456112360150a6d0d1421/ast/rewrite.go#L47)
- Domain: Contains domain object required by rewriter to interface with underlying database to fetch details like table
 schema, indexes, column names etc.
- Catalog API: Contains API interface to interact with underlying database.
- Catalog MySQL: Contains implementation of Catalog API for MySQL database.