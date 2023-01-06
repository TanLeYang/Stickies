package stickiesset

import "gorm.io/gorm"

type StickiesSet struct {
	gorm.Model
	Owner			 int64 
	TgStickerSetName string
	UniqueCode       string
}

type StickiesSetRepository interface {
	Create(s StickiesSet) error
	GetByOwner(owner int64) ([]StickiesSet, error)
	GetByUniqueCode(code string) (*StickiesSet, error)
}

type StickiesSetDb struct {
	Db *gorm.DB
}

func (r StickiesSetDb) Create(s StickiesSet) error {
	result := r.Db.Create(&s)
	return result.Error
}

func (r StickiesSetDb) GetByOwner(owner int64) ([]StickiesSet, error) {
	stickiesSets := []StickiesSet{}
	result := r.Db.Where(&StickiesSet{Owner: owner}).Find(&stickiesSets)
	return stickiesSets, result.Error
}

func (r StickiesSetDb) GetByUniqueCode(code string) (*StickiesSet, error) {
	var stickiesSet StickiesSet;
	result := r.Db.Where(&StickiesSet{UniqueCode: code}).First(&stickiesSet)
	return &stickiesSet, result.Error
}
