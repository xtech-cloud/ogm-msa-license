package model

import (
	"time"
    "errors"

	"gorm.io/gorm"
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
	return "ogm_license_key"
}

type KeyDAO struct {
	conn *Conn
}

func NewKeyDAO(_conn *Conn) *KeyDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &KeyDAO{
        conn: conn,
    }
}

func (this *KeyDAO) Insert(_key Key) error {
    db := this.conn.DB
	return db.Create(&_key).Error
}

func (this *KeyDAO) Count(_space string) (int64, error) {
    db := this.conn.DB
	count := int64(0)

	res := db.Model(&Key{}).Where("space = ?", _space).Count(&count)
	return count, res.Error
}

func (this *KeyDAO) Find(_number string) (Key, error) {
    db := this.conn.DB
	var key Key

	res := db.Where("number = ?", _number).First(&key)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Key{}, nil
	}
	return key, res.Error
}

func (this *KeyDAO) Save(_key *Key) error {
    db := this.conn.DB
	return db.Save(_key).Error
}

func (this *KeyDAO) List(_offset int32, _count int32, _space string) ([]Key, error) {
    db := this.conn.DB

	var key []Key
	res := db.Where("space = ?", _space).Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&key)
	return key, res.Error
}
