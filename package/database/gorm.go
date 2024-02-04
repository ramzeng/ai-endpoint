package database

import "gorm.io/gorm"

func Where(query interface{}, args ...interface{}) *gorm.DB {
	return db.Where(query, args...)
}
