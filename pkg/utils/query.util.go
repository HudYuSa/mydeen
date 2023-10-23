package utils

import "gorm.io/gorm"

func SelectColumnDB(column ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(column)
	}
}
