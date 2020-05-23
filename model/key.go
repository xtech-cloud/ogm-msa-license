package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Key struct {
	GModel      gorm.Model `gorm:"embedded"`
	Number      string     `gorm:"column:number;type:char(32);unique;not null"`
	Space       string     `gorm:"column:space;type:varchar(32);not null"`
	Capacity    int32      `gorm:"column:capacity;not null;default:1"`
	Expiry      int32      `gorm:"column:expiry;not null;default:0"`
	Ban         int32      `gorm:"column:ban;not null;default:0"`
	Storage     string     `gorm:"column:storage"`
	Profile     string     `gorm:"column:profile;type:TEXT"`
	ActivatedAt time.Time  `gorm:"column:activated_at;`
}

func (Key) TableName() string {
	return "msa_license_key"
}

type KeyDAO struct {
}

func NewKeyDAO() *KeyDAO {
	return &KeyDAO{}
}

func (KeyDAO) Insert(_key Key) error {
	db, err := openSqlDB()
	if nil != err {
		return err
	}
	defer closeSqlDB(db)

	return db.Create(&_key).Error
}

func (KeyDAO) Find(_number string) (Key, error) {
	var key Key
	db, err := openSqlDB()
	if nil != err {
		return key, err
	}
	defer closeSqlDB(db)

	res := db.Where("number = ?", _number).First(&key)
	if res.RecordNotFound() {
		return Key{}, nil
	}
	return key, err
}

func (KeyDAO) Save(_key *Key) error {
	db, err := openSqlDB()
	if nil != err {
		return err
	}
	defer closeSqlDB(db)

	return db.Save(_key).Error
}
