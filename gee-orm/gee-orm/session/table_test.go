package session

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type User struct {
	Name string `geeorm: "PRIMARY KEY"`
	Age int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("failed to create table User")
	}
	_ = s.DropTable()
	if s.HasTable() {
		t.Fatal("failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	s := NewSession().Model(&User{})
	table := s.RefTable()

	s.Model(&Session{})
	if table.Name != "User" || s.RefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}
