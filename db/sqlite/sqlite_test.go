package sqlite

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetTableNames(t *testing.T) {
	fmt.Println(" TestGetTableNames")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parent (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = db.Exec("CREATE TABLE child (id INTEGER PRIMARY KEY, parent_id INTEGER, name TEXT)")
	if err != nil {
		t.Error(err)
		return
	}

	tables, err := GetTableNames(db)
	if err != nil { t.Error(err) }
	assert.True(t, len(tables) == 2, "tables length is not 2: " + strconv.Itoa(len(tables)))
	assert.True(t, tables[0] == "child", "tables[0] is not child")
	assert.True(t, tables[1] == "parent", "tables[1] is not parent")
}

func TestGetFieldInfos(t *testing.T) {
	fmt.Println(" TestGetFieldInfos")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)")
	if err != nil {
		t.Error(err)
		return
	}

	fieldInfos, err := GetFieldInfos(db, "test")
	if err != nil { t.Error(err) }
	assert.True(t, len(fieldInfos) == 3, "fieldInfos length is not 3")
	assert.True(t, fieldInfos[0].Name == "id", "fieldInfos[0].Name is not id")
	assert.True(t, fieldInfos[0].FieldType == "INTEGER", "fieldInfos[0].Type is not INTEGER")
	assert.True(t, !fieldInfos[0].NotNull, "fieldInfos[0].NotNull is not true")
	assert.True(t, fieldInfos[0].PrimaryKey, "fieldInfos[0].PrimaryKey is not false")

	assert.True(t, fieldInfos[1].Name == "name", "fieldInfos[1].Name is not name")
	assert.True(t, fieldInfos[1].FieldType == "TEXT", "fieldInfos[1].Type is not TEXT")
	assert.True(t, !fieldInfos[1].NotNull, "fieldInfos[1].NotNull is not true")
	assert.True(t, !fieldInfos[1].PrimaryKey, "fieldInfos[1].PrimaryKey is not false")

	assert.True(t, fieldInfos[2].Name == "age", "fieldInfos[2].Name is not age")
	assert.True(t, fieldInfos[2].FieldType == "INTEGER", "fieldInfos[2].Type is not INTEGER")
	assert.True(t, !fieldInfos[2].NotNull, "fieldInfos[2].NotNull is not true")
	assert.True(t, !fieldInfos[2].PrimaryKey, "fieldInfos[2].PrimaryKey is not false")
}

func TestGetForeignKeyInfos(t *testing.T) {
	fmt.Println(" TestGetForeignKeyInfos")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parent (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = db.Exec("CREATE TABLE child (id INTEGER PRIMARY KEY, parent_id INTEGER, name TEXT, FOREIGN KEY(parent_id) REFERENCES parent(id))")
	if err != nil {
		t.Error(err)
		return
	}

	fkInfos, err := GetForeignKeyInfos(db, "child")
	if err != nil { t.Error(err) }
	assert.True(t, len(fkInfos) == 1, "fkInfos length is not 1")
	assert.True(t, fkInfos[0].Table == "parent", "fkInfos[0].Table is not parent")
	assert.True(t, fkInfos[0].From == "parent_id", "fkInfos[0].From is not parent_id")
	assert.True(t, fkInfos[0].To == "id", "fkInfos[0].To is not id")
}

func TestMakeSelectSQLWithFkAssociation(t *testing.T) {
	fmt.Println(" TestMakeSelectSQLWithFkAssociation")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parent (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.Exec("CREATE TABLE child (id INTEGER PRIMARY KEY, parent_id INTEGER, name TEXT, FOREIGN KEY(parent_id) REFERENCES parent(id))")
	if err != nil {
		t.Error(err)
		return
	}
	sql, err := MakeSelectSQLWithFkAssociation(db, "child")
	if err != nil { t.Error(err) }
	re := regexp.MustCompile("[ \t\n]+")
	sql = re.ReplaceAllString(sql, " ")
	assert.True(t, sql == "SELECT child.* , parent.id AS parent_id, parent.name AS parent_name FROM child LEFT JOIN parent ON child.parent_id = parent.id ", 
		"sql error: [" + sql + "]")
}