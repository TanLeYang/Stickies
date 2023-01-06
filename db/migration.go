package db

import (
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	"gorm.io/gorm"
)

func PerformMigration(db *gorm.DB) {
	db.AutoMigrate(&stickiesset.StickiesSet{})
}
