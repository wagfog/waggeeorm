package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string //字段名 Name
	Type string //类型 Type
	Tag  string //约束条件 Tag
}

type Schema struct {
	Model      interface{}       //被映射的对象 Model
	Name       string            //表名 Name
	Fields     []*Field          //字段 Fields。
	FieldNames []string          //包含所有的字段名(列名)
	fileMap    map[string]*Field //fieldMap 记录字段名和 Field 的映射关系，方便之后直接使用，无需遍历 Fields。
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fileMap[name]
}

// point and Dialect
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	//TypeOf() 和 ValueOf() 是 reflect 包最为基本也是最重要的 2 个方法
	//因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:   dest,             //被映射的对象 Model
		Name:    modelType.Name(), //modelType.Name() 获取到结构体的名称作为表名
		fileMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ { //NumField() 获取实例的字段的个数
		p := modelType.Field(i)                     //通过下标获取到特定字段 p := modelType.Field(i)。
		if !p.Anonymous && ast.IsExported(p.Name) { //判断变量 p 的类型是否是导出的（即对外可见 ast.IsExported
			field := &Field{
				Name: p.Name,                                              //p.Name 即字段名
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), //p.Type 即字段类型. 通过 (Dialect).DataTypeOf() 转换为数据库的字段类
			}

			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fileMap[p.Name] = field
		}
	}
	return schema
}
