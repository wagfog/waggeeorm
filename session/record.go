package session

import (
	"fmt"
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames) //次调用 clause.Set() 构造好每一个子句。
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES) //调用一次 clause.Build() 按照传入的顺序构造出最终的 SQL 语句。
	fmt.Println(sql)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	//reflect.Value代表一个具体的值，并提供了一系列方法来对这个值进行读取和操作。
	//通过reflect.Value，我们可以获取值的类型、值的字段值和方法，以及对值进行修改等操作
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	//reflect.Type则代表一个类型，不具体代表某个具体的值。
	//reflect.Type提供了一些方法来获取类型的信息，比如类型的名称、包路径、字段、方法等。它可以用于获取类型的静态信息，而不是具体的值。
	//destSlice.Type().Elem() 获取切片的单个元素的类型 destType
	destType := destSlice.Type().Elem()
	//使用 reflect.New() 方法创建一个 destType 的实例
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	//根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	//遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		//调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
