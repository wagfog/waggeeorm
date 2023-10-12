package clause

import "strings"

type Type int

// 实现结构体 Clause 拼接各个独立的子句。
type Clause struct {
	//组成总sql语句的子sql语句
	sql map[Type]string
	//对应子sql语句的变量
	sqlVars map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

// 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句。
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
