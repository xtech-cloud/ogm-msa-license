package model

import (
	"time"
)

type Key struct {
	UUID        string `gorm:"column:uuid;type:char(32);unique;not null;primaryKey"`
	Number      string `gorm:"column:number;type:char(32);unique;not null"`
	Space       string `gorm:"column:space;type:varchar(32);not null"`
	Capacity    int32  `gorm:"column:capacity;not null;default:1"`
	Expiry      int32  `gorm:"column:expiry;not null;default:0"`
	Ban         int32  `gorm:"column:ban;not null;default:0"`
	Reason      string `gorm:"column:reason;type:varchar(256);not null;default:''"`
	Storage     string `gorm:"column:storage"`
	Profile     string `gorm:"column:profile;type:TEXT"`
	CreatedAt   time.Time
	ActivatedAt time.Time `gorm:"column:activated_at;`
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

func (this *KeyDAO) Get(_uuid string) (*Key, error) {
	db := this.conn.DB
	var key Key

	res := db.Where("uuid = ?", _uuid).First(&key)
	return &key, res.Error
}

func (this *KeyDAO) Find(_number string) (*Key, error) {
	db := this.conn.DB
	var key Key

	res := db.Where("number = ?", _number).First(&key)
	return &key, res.Error
}

func (this *KeyDAO) Save(_key *Key) error {
	db := this.conn.DB
	return db.Save(_key).Error
}

func (this *KeyDAO) List(_offset int64, _count int64, _space string) (int64, []*Key, error) {
	db := this.conn.DB.Model(&Key{}).Where("`space` = ?", _space)
	var count int64
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}

	var key []*Key
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&key)
	return count, key, res.Error
}

func (this *KeyDAO) Search(_offset int64, _count int64, _space string, _number string, _capacity int32, _expiry int32, _storage string, _profile string, _ban int32) (int64, []*Key, error) {
	db := this.conn.DB.Model(&Key{}).Where("`space` = ?", _space)
	if "" != _number {
		db = db.Where("`number` LIKE ?", "%"+_number+"%")
	}
	if 0 != _capacity {
		db = db.Where("`capacity` = ?", _capacity)
	}
	if _expiry > 0 {
		db = db.Where("`expiry` = ?", _expiry)
	}
	if _ban > 0 {
		db = db.Where("`ban` = ?", _ban)
	}
	if "" != _storage {
		db = db.Where("`storage` LIKE ?", "%"+_storage+"%")
	}
	if "" != _profile {
		db = db.Where("`profile` LIKE ?", "%"+_profile+"%")
	}
	var count int64
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}
	var keys []*Key
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&keys)
	return count, keys, res.Error
}
