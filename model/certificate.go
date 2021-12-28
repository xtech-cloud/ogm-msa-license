package model

import (
	"time"
    "errors"

	"gorm.io/gorm"
)

type Certificate struct {
	UUID      string `gorm:"column:uuid;type:char(32);unique;not null;primaryKey"`
	Space     string `gorm:"column:space;type:varchar(32);not null"`
	Key       string `gorm:"column:number;type:char(32);not null"`
	Consumer  string `gorm:"column:consumer;type:varchar(128);not null"`
	Content   string `gorm:"column:content;type:TEXT;not null"`
	CreatedAt time.Time
}

func (Certificate) TableName() string {
	return "ogm_license_certificate"
}

type CertificateQuery struct {
	Space    string
	Consumer string
	Number   string
}

type CertificateDAO struct {
	conn *Conn
}

func NewCertificateDAO(_conn *Conn) *CertificateDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &CertificateDAO{
		conn: conn,
	}
}

func (this *CertificateDAO) Insert(_cer Certificate) error {
	db := this.conn.DB
	return db.Create(&_cer).Error
}

func (this *CertificateDAO) Get(_uuid string) (*Certificate, error) {
	db := this.conn.DB

	var cer Certificate
	res := db.Where("uuid = ?", _uuid).First(&cer)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cer, res.Error
}

func (this *CertificateDAO) Query(_query CertificateQuery) ([]*Certificate, error) {
	db := this.conn.DB

	var cers []*Certificate

	db = db.Model(&Certificate{}).Order("created_at desc")
	blankQuery := true
	if "" != _query.Space {
		db = db.Where("space = ?", _query.Space)
		blankQuery = false
	}
	if "" != _query.Consumer {
		db = db.Where("consumer = ?", _query.Consumer)
		blankQuery = false
	}
	if "" != _query.Number {
		db = db.Where("number = ?", _query.Number)
		blankQuery = false
	}

	if blankQuery {
		return make([]*Certificate, 0), nil
	}

	res := db.Find(&cers)
	return cers, res.Error
}

func (this *CertificateDAO) Count(_query CertificateQuery) (int64, error) {
	db := this.conn.DB

	count := int64(0)

	db = db.Model(&Certificate{})

	if "" != _query.Space {
		db = db.Where("space = ?", _query.Space)
	}
	if "" != _query.Consumer {
		db = db.Where("consumer = ?", _query.Consumer)
	}
	if "" != _query.Number {
		db = db.Where("number = ?", _query.Number)
	}

	res := db.Count(&count)
	return count, res.Error
}

func (this *CertificateDAO) List(_offset int64, _count int64, _space string) (int64, []*Certificate, error) {
	db := this.conn.DB.Model(&Certificate{}).Where("`space` = ?", _space)
	count := int64(0)
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}

	var cer []*Certificate
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&cer)
	return count, cer, res.Error
}

func (this *CertificateDAO) Search(_offset int64, _count int64, _space string, _number string, _consumer string) (int64, []*Certificate, error) {
	db := this.conn.DB.Model(&Certificate{}).Where("`space` = ?", _space)
	if "" != _number {
		db = db.Where("`number` LIKE ?", "%"+_number+"%")
	}
	if "" != _consumer {
		db = db.Where("`consumer` LIKE ?", "%"+_consumer+"%")
	}
	count := int64(0)
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}

	var cer []*Certificate
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&cer)
	return count, cer, res.Error
}
