package tests

import "gorm.io/gorm"

func ClearTables(db *gorm.DB, tables []string) {
	for _, table := range tables {
		db.Exec("delete from " + table)
		// Resetar auto-increment para SQLite
		db.Exec("update sqlite_sequence set seq = 0 where name = ?", table)
	}
}
