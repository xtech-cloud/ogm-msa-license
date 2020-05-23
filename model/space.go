package model

import (
	"github.com/jinzhu/gorm"
)

type Space struct {
	GModel      gorm.Model `gorm:"embedded"`
	Name        string     `gorm:"column:name;type:varchar(32);unique;not null"`
	SpaceKey    string     `gorm:"column:space_key;type:char(32);unique;not null"`
	SpaceSecret string     `gorm:"column:space_secret;type:char(32);unique;not null"`
	PublicKey   string     `gorm:"column:key_public;type:TEXT;not null"`
	PrivateKey  string     `gorm:"column:key_private;type:TEXT;not null"`
	Profile     string     `gorm:"column:profile;type:TEXT"`
}

func (Space) TableName() string {
	return "msa_license_space"
}

type SpaceQuery struct {
	Name        string
	SpaceKey    string
	SpaceSecret string
}

type SpaceDAO struct {
}

func NewSpaceDAO() *SpaceDAO {
	return &SpaceDAO{}
}

func (SpaceDAO) Exists(_name string) (bool, error) {
	db, err := openSqlDB()
	if nil != err {
		return false, err
	}
	defer closeSqlDB(db)

	var space Space
	result := db.Where("name = ?", _name).First(&space)
	if result.RecordNotFound() {
		return false, nil
	}

	return "" != space.Name, result.Error
}

func (SpaceDAO) Insert(_space Space) error {
	db, err := openSqlDB()
	if nil != err {
		return err
	}
	defer closeSqlDB(db)

	return db.Create(&_space).Error
}

func (SpaceDAO) Find(_name string) (Space, error) {
	var space Space
	db, err := openSqlDB()
	if nil != err {
		return space, err
	}
	defer closeSqlDB(db)

	res := db.Where("name = ?", _name).First(&space)
	if res.RecordNotFound() {
		return Space{}, nil
	}
	return space, err
}

func (SpaceDAO) Count() (int64, error) {
	count := int64(0)
	db, err := openSqlDB()
	if nil != err {
		return count, err
	}
	defer closeSqlDB(db)

	res := db.Model(&Space{}).Count(&count)
	return count, res.Error
}

func (SpaceDAO) Fetch(_key string, _secret string) (Space, error) {
	var space Space
	db, err := openSqlDB()
	if nil != err {
		return space, err
	}
	defer closeSqlDB(db)

	res := db.Where("space_key = ? AND space_secret = ?", _key, _secret).First(&space)
	if res.RecordNotFound() {
		return Space{}, nil
	}
	return space, err
}

func (SpaceDAO) List(_offset int32, _count int32) ([]Space, error) {
	db, err := openSqlDB()
	if nil != err {
		return nil, err
	}
	defer closeSqlDB(db)

	var space []Space
	res := db.Offset(_offset).Limit(_count).Order("created_at desc").Find(&space)
	return space, res.Error
}
