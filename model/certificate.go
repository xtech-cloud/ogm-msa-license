package model

import (
	"github.com/jinzhu/gorm"
)

type Certificate struct {
	GModel   gorm.Model `gorm:"embedded"`
	UID      string     `gorm:"column:uid;type:char(32);unique;not null"`
	Space    string     `gorm:"column:space;type:varchar(32);not null"`
	Key      string     `gorm:"column:number;type:char(32);not null"`
	Consumer string     `gorm:"column:consumer;type:varchar(128);not null"`
	Content  string     `gorm:"column:content;type:TEXT;not null"`
}

func (Certificate) TableName() string {
	return "msa_license_certificate"
}

type CertificateQuery struct {
	Space    string
	Consumer string
	Number   string
}

type CertificateDAO struct {
}

func NewCertificateDAO() *CertificateDAO {
	return &CertificateDAO{}
}

func (CertificateDAO) Insert(_cer Certificate) error {
	db, err := openSqlDB()
	if nil != err {
		return err
	}
	defer closeSqlDB(db)

	return db.Create(&_cer).Error
}

func (CertificateDAO) Find(_uid string) (Certificate, error) {
	var cer Certificate
	db, err := openSqlDB()
	if nil != err {
		return cer, err
	}
	defer closeSqlDB(db)

	res := db.Where("uid = ?", _uid).First(&cer)
	if res.RecordNotFound() {
		return Certificate{}, nil
	}
	return cer, err
}

func (CertificateDAO) Query(_query CertificateQuery) ([]*Certificate, error) {
	var cers []*Certificate
	db, err := openSqlDB()
	if nil != err {
		return nil, err
	}
	defer closeSqlDB(db)

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

func (CertificateDAO) Count(_query CertificateQuery) (int64, error) {
	count := int64(0)
	db, err := openSqlDB()
	if nil != err {
		return count, err
	}
	defer closeSqlDB(db)

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

func (CertificateDAO) List(_offset int32, _count int32, _space string) ([]Certificate, error) {
	db, err := openSqlDB()
	if nil != err {
		return nil, err
	}
	defer closeSqlDB(db)

	var cer []Certificate
	res := db.Where("space = ?", _space).Offset(_offset).Limit(_count).Order("created_at desc").Find(&cer)
	return cer, res.Error
}
