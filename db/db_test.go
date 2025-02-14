package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonify(t *testing.T) {
	fmt.Println(" TestJsonify")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE "testJson" (
		"id"	INTEGER,
		"name"	TEXT,
		"weight"	NUMERIC,
		"height"	REAL,
		"age"	INTEGER,
		"photo"	BLOB,
		PRIMARY KEY("id"))
	`)
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := db.Query("SELECT * FROM testJson")
	if err != nil {
		t.Error(err)
		return
	}
	json, err := Jsonify(rows)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, json == "[]", "json is not [] : " + json)	

	_, err = db.Exec(`INSERT INTO "testJson" VALUES (1, "John", 70.5, 180.5, 30, NULL)`)
	if err != nil { 
		t.Error(err)
		return
	}
	rows, err = db.Query("SELECT * FROM testJson")
	if err != nil {
		t.Error(err)
		return
	}
	json, err = Jsonify(rows)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, json == `[{"id": 1, "name": "John", "weight": 70.5, "height": 180.5, "age": 30, "photo": null}]`, "json error: " + json)

	_, err = db.Exec(`INSERT INTO "testJson" VALUES (2, "Tom", 73.12, 180.34, 31, "aa")`)
	if err != nil {
		t.Error(err)
		return
	}
	rows, err = db.Query("SELECT * FROM testJson")
	if err != nil {
		t.Error(err)
		return
	}
	json, err = Jsonify(rows)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, json == `[{"id": 1, "name": "John", "weight": 70.5, "height": 180.5, "age": 30, "photo": null}, {"id": 2, "name": "Tom", "weight": 73.12, "height": 180.34, "age": 31, "photo": "aa"}]`, 
		"json error: " + json)

}
