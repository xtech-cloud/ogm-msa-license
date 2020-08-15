package handler

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"omo-msa-license/crypto"
	"omo-msa-license/model"

	"github.com/micro/go-micro/v2/logger"
	uuid "github.com/satori/go.uuid"

	proto "github.com/xtech-cloud/omo-msp-license/proto/license"
)

var base64Coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

type Key struct{}

func (this *Key) Generate(_ctx context.Context, _req *proto.KeyGenerateRequest, _rsp *proto.KeyGenerateResponse) error {
	logger.Infof("Received Key.Generate, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
		return nil
	}

	capacity := int32(1)
	if _req.Capacity > 0 {
		capacity = _req.Capacity
	}

	count := int32(1)
	if _req.Count > 0 {
		count = _req.Count
	}

	expiry := int32(0)
	if _req.Expiry > 0 {
		expiry = _req.Expiry
	}

	daoSpace := model.NewSpaceDAO()
	daoKey := model.NewKeyDAO()

	space, err := daoSpace.Find(_req.Space)
	if nil != err {
		return err
	}

	if "" == space.Name {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "space not found"
		return nil
	}

	_rsp.Number = make([]string, 0)
	for i := int32(0); i < count; i++ {
		number, err := newNumber()
		if nil != err {
			continue
		}
		key := model.Key{
			Number:   number,
			Space:    space.Name,
			Capacity: capacity,
			Expiry:   expiry,
			Ban:      0,
			Storage:  _req.Storage,
			Profile:  _req.Profile,
		}
		err = daoKey.Insert(key)
		if nil != err {
			logger.Error(err.Error())
			continue
		}
		_rsp.Number = append(_rsp.Number, key.Number)
	}
	return nil
}

func (this *Key) Query(_ctx context.Context, _req *proto.KeyQueryRequest, _rsp *proto.KeyQueryResponse) error {
	logger.Infof("Received Key.Query, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Number {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "number is required"
		return nil
	}

	dao := model.NewKeyDAO()

	key, err := dao.Find(_req.Number)
	if nil != err {
		return err
	}

	if "" == key.Number {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "key not found"
		return nil
	}

	_rsp.Key = &proto.KeyEntity{
		Number:      key.Number,
		Space:       key.Space,
		Capacity:    key.Capacity,
		Expiry:      key.Expiry,
		Storage:     key.Storage,
		Profile:     key.Profile,
		Ban:         key.Ban,
		CreatedAt:   key.GModel.CreatedAt.Unix(),
		UpdatedAt:   key.GModel.UpdatedAt.Unix(),
		ActivatedAt: key.ActivatedAt.Unix(),
	}
	if _rsp.Key.ActivatedAt < _rsp.Key.CreatedAt {
		_rsp.Key.ActivatedAt = 0
	}

	daoCer := model.NewCertificateDAO()
	// 获取已激活的证书
	cers, err := daoCer.Query(model.CertificateQuery{
		Space:  key.Space,
		Number: key.Number,
	})
	if nil != err {
		return err
	}
	_rsp.Key.Consumer = make([]string, len(cers))
	for i, key := range cers {
		_rsp.Key.Consumer[i] = key.Consumer
	}
	return nil
}

func (this *Key) List(_ctx context.Context, _req *proto.KeyListRequest, _rsp *proto.KeyListResponse) error {
	logger.Infof("Received Key.List, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
		return nil
	}

	dao := model.NewKeyDAO()

	count, err := dao.Count(_req.Space)
	// 数据库错误
	if nil != err {
		return err
	}

	keys, err := dao.List(_req.Offset, _req.Count, _req.Space)
	// 数据库错误
	if nil != err {
		return err
	}

	daoCer := model.NewCertificateDAO()
	_rsp.Total = count
	_rsp.Key = make([]*proto.KeyEntity, len(keys))
	for i, key := range keys {
		_rsp.Key[i] = &proto.KeyEntity{
			Number:      key.Number,
			Space:       key.Space,
			Capacity:    key.Capacity,
			Expiry:      key.Expiry,
			Storage:     key.Storage,
			Profile:     key.Profile,
			Ban:         key.Ban,
			CreatedAt:   key.GModel.CreatedAt.Unix(),
			UpdatedAt:   key.GModel.UpdatedAt.Unix(),
			ActivatedAt: key.ActivatedAt.Unix(),
		}
		if _rsp.Key[i].ActivatedAt < _rsp.Key[i].CreatedAt {
			_rsp.Key[i].ActivatedAt = 0
		}
		// 获取已激活的消费者
		consumers, err := daoCer.Query(model.CertificateQuery{
			Space:  key.Space,
			Number: key.Number,
		})
		if nil != err {
			continue
		}
		_rsp.Key[i].Consumer = make([]string, len(consumers))
		for j, c := range consumers {
			_rsp.Key[i].Consumer[j] = c.Consumer
		}
	}

	return nil
}

func (this *Key) Activate(_ctx context.Context, _req *proto.KeyActivateRequest, _rsp *proto.KeyActivateResponse) error {
	logger.Infof("Received Key.Activate, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Number {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "number is required"
		return nil
	}

	if "" == _req.Consumer {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "consumer is required"
		return nil
	}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
		return nil
	}

	daoSpace := model.NewSpaceDAO()
	space, err := daoSpace.Find(_req.Space)
	if nil != err {
		return err
	}

	if "" == space.Name {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "space not found"
		return nil
	}

	daoKey := model.NewKeyDAO()
	key, err := daoKey.Find(_req.Number)
	if nil != err {
		return err
	}
	if "" == key.Number {
		_rsp.Status.Code = 3
		_rsp.Status.Message = "key not found"
		return nil
	}

	if key.Space != space.Name {
		_rsp.Status.Code = 4
		_rsp.Status.Message = "space not matched"
		return nil
	}

	if key.Ban > 0 {
		_rsp.Status.Code = 5
		_rsp.Status.Message = "ban > 0"
		return nil
	}

	// 如果存在已激活的有效证书，则直接返回
	uid := model.ToUUID(fmt.Sprintf("%s%s%s", space.Name, key.Number, _req.Consumer))
	daoCer := model.NewCertificateDAO()
	certificate, err := daoCer.Find(uid)
	if nil != err {
		return nil
	}
	if "" != certificate.UID {
		_rsp.CerContent = certificate.Content
		return nil
	}

	// 获取已激活的数量
	count, err := daoCer.Count(model.CertificateQuery{
		Space:  _req.Space,
		Number: _req.Number,
	})
	if nil != err {
		return err
	}

	// 已达到激活码的激活能力
	if int32(count) >= key.Capacity {
		_rsp.Status.Code = 6
		_rsp.Status.Message = "out of capacity"
		return nil
	}

	// 新建证书
	cer, err := makeCertificate(space.SpaceKey, space.SpaceSecret, _req.Consumer, key.Storage, key.Expiry, space.PublicKey, space.PrivateKey)
	if nil != err {
		return err
	}

	newCer := model.Certificate{
		UID:      uid,
		Space:    space.Name,
		Consumer: _req.Consumer,
		Key:      key.Number,
		Content:  cer,
	}

	// 保存证书
	err = daoCer.Insert(newCer)
	if nil != err {
		return err
	}

	if key.ActivatedAt.Unix() < key.GModel.CreatedAt.Unix() {
		key.ActivatedAt = time.Now()
	}
	daoKey.Save(&key)

	_rsp.CerUID = newCer.UID
	_rsp.CerContent = newCer.Content
	return nil
}

func (this *Key) Suspend(_ctx context.Context, _req *proto.KeySuspendRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Key.Suspend, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
		return nil
	}

	if "" == _req.Number {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "number is required"
		return nil
	}

	if "" == _req.Reason {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "reason is required"
		return nil
	}

	dao := model.NewKeyDAO()

	key, err := dao.Find(_req.Number)
	if nil != err {
		return err
	}

	if "" == key.Number {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "key not found"
		return nil
	}

	key.Ban = _req.Ban
	err = dao.Save(&key)
	if nil != err {
		return nil
	}
	return nil
}

func newNumber() (string, error) {
	id := uuid.NewV4()
	h := md5.New()
	h.Write(id.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}

func makeCertificate(_spaceKey string, _spaceSecret string, _consumer string, _storage string, _expiry int32, _publicKey string, _privateKey string) (string, error) {
	now := time.Now().Unix()

	passwd := toPassword(_spaceKey, _spaceSecret)

	pub_ciphertext, err := crypto.EncryptAES([]byte(_publicKey), []byte(passwd))
	if nil != err {
		return "", err
	}
	pub := base64Coder.EncodeToString(pub_ciphertext)

	//generate payload
	payload := fmt.Sprintf("spacekey:\n%s\nconsumer:\n%s\ntimestamp:\n%d\nexpiry:\n%d\nstorage:\n%s\npubkey:\n%s",
		_spaceKey, _consumer, now, _expiry, _storage, pub)
	identity_ciphertext, err := crypto.EncryptAES([]byte(payload), []byte(passwd))
	identity := toMD5(identity_ciphertext)
	sig_ciphertext, err := crypto.SignRSA([]byte(_privateKey), []byte(identity))
	if nil != err {
		return "", err
	}
	sig := base64Coder.EncodeToString(sig_ciphertext)
	license := fmt.Sprintf("%s\nsig:\n%s", payload, sig)
	return license, nil
}

func toMD5(_val []byte) string {
	hash := md5.New()
	hash.Write(_val)
	return hex.EncodeToString(hash.Sum(nil))
}

func toPassword(_key string, _secret string) string {
	hash := md5.New()
	hash.Write([]byte(_key + _secret))
	pwd := hex.EncodeToString(hash.Sum(nil))
	return strings.ToUpper(pwd)
}
