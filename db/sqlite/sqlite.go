package sqlite

import (
	"database/sql"
	"strings"
)

type FieldInfo struct {
	Cid          int
	Name         string
	FieldType    string
	NotNull      bool
	DefaultValue sql.NullString
	PrimaryKey   bool
}

func GetFieldInfos(db *sql.DB, tableName string) ([]FieldInfo, error) {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fieldInfos []FieldInfo
	for rows.Next() {
		var info FieldInfo
		err = rows.Scan(&info.Cid, &info.Name, &info.FieldType, &info.NotNull, &info.DefaultValue, &info.PrimaryKey)
		if err != nil {
			return nil, err
		}
		fieldInfos = append(fieldInfos, info)
	}
	return fieldInfos, nil
}

type ForeignKeyInfo struct {
	Id		int
	Seq		int
	Table	string
	From	string
	To		string
	OnUpdate	string
	OnDelete	string
	Match	string
}

func GetForeignKeyInfos(db *sql.DB, tableName string) ([]ForeignKeyInfo, error) {
	rows, err := db.Query("PRAGMA foreign_key_list(" + tableName + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fkInfos []ForeignKeyInfo
	for rows.Next() {
		var info ForeignKeyInfo
		err = rows.Scan(&info.Id, &info.Seq, &info.Table, &info.From, &info.To, &info.OnUpdate, &info.OnDelete, &info.Match)
		if err != nil {
			return nil, err
		}
		fkInfos = append(fkInfos, info)
	}
	return fkInfos, nil
}

func MakeSelectSQLWithFkAssociation(db *sql.DB, tableName string) (string, error) {
	basicSQL := "SELECT * FROM " + tableName
	fkInfos, err := GetForeignKeyInfos(db, tableName)
	if err != nil {
		return basicSQL, err
	}
	selectList := []string{ "SELECT " + tableName + ".* \n" }
	fromList := []string{ "FROM " + tableName + "\n" }
	for _, fkInfo := range fkInfos {
		fieldInfos, err := GetFieldInfos(db, fkInfo.Table)
		if err != nil {
			return basicSQL, err
		}
		selectList = append(selectList, "  ")

		for _, fieldInfo := range fieldInfos {
			selectList = append(selectList, ", " + fkInfo.Table + "." + fieldInfo.Name + " AS " + fkInfo.Table + "_" + fieldInfo.Name)
		}
		selectList = append(selectList, "\n")
		fromList = append(fromList, " LEFT JOIN " + fkInfo.Table + " ON " + tableName + "." + fkInfo.From + " = " + fkInfo.Table + "." + fkInfo.To + "\n")	
	}
	return strings.Join(selectList, "") + strings.Join(fromList, ""), nil
}