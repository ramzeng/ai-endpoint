package database

import "gorm.io/gorm"
import "gorm.io/driver/mysql"

func Initialize(config Config) error {
	var err error

	db, err = gorm.Open(mysql.Open(config.Addr), &gorm.Config{})

	if err != nil {
		return err
	}

	return nil
}
