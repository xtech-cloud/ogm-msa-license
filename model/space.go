package model

import (
    "errors"
	"gorm.io/gorm"
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
	return "ogm_license_space"
}

type SpaceQuery struct {
	Name        string
	SpaceKey    string
	SpaceSecret string
}

type SpaceDAO struct {
	conn *Conn
}

func NewSpaceDAO(_conn *Conn) *SpaceDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &SpaceDAO{
        conn: conn,
    }
}

func (this *SpaceDAO) Exists(_name string) (bool, error) {
    db := this.conn.DB

	var space Space
	result := db.Where("name = ?", _name).First(&space)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return "" != space.Name, result.Error
}

func (this* SpaceDAO) Insert(_space Space) error {
    db := this.conn.DB

	return db.Create(&_space).Error
}

func (this* SpaceDAO) Find(_name string) (Space, error) {
    db := this.conn.DB

	var space Space

	res := db.Where("name = ?", _name).First(&space)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Space{}, nil
	}
	return space, res.Error
}

func (this* SpaceDAO) Count() (int64, error) {
    db := this.conn.DB

	count := int64(0)

	res := db.Model(&Space{}).Count(&count)
	return count, res.Error
}

func (this* SpaceDAO) Fetch(_key string, _secret string) (Space, error) {
    db := this.conn.DB

	var space Space

	res := db.Where("space_key = ? AND space_secret = ?", _key, _secret).First(&space)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Space{}, nil
	}
	return space, res.Error
}

func (this* SpaceDAO) List(_offset int32, _count int32) ([]Space, error) {
    db := this.conn.DB

	var space []Space
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&space)
	return space, res.Error
}
