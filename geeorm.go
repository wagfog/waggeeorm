package geeorm

//那交互前的准备工作（比如连接/测试数据库），交互后的收尾工作（关闭连接）等就交给 Engine 来负责了
import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// Engine 的逻辑非常简单，最重要的方法是 NewEngine
// NewEngine 主要做了两件事。
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source) //1...连接数据库，返回 *sql.DB。
	if err != nil {
		log.Error(err)
		return
	}
	// send a ping to make sure the database connection is alive
	if err = db.Ping(); err != nil { //2...调用 db.Ping()，检查数据库是否能够正常连接。
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("CLose database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
