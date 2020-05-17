package session

import (
	"database/sql"
	"geeorm/dialect"
	"os"
	"testing"
)

var (
	TestDB *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../../gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(Name) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	rows := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := rows.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
