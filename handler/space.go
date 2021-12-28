package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"ogm-msa-license/crypto"
	"ogm-msa-license/model"

	"github.com/asim/go-micro/v3/logger"

	proto "github.com/xtech-cloud/ogm-msp-license/proto/license"
)

type Space struct{}

func (this *Space) Create(_ctx context.Context, _req *proto.SpaceCreateRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Space.Create, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewSpaceDAO(nil)

	now := time.Now().Unix()
	keyCode := fmt.Sprintf("%v-%v-key", _req.Name, now)
	secretCode := fmt.Sprintf("%v-%v-secret", _req.Name, now)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(keyCode))
	spaceKey := hex.EncodeToString(md5Ctx.Sum(nil))
	md5Ctx.Write([]byte(secretCode))
	spaceSecret := hex.EncodeToString(md5Ctx.Sum(nil))
	publicKey, privateKey, err := crypto.GenerateKeyRSA()
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	space := model.Space{
		UUID:        model.ToUUID(_req.Name),
		Name:        _req.Name,
		SpaceKey:    spaceKey,
		SpaceSecret: spaceSecret,
		PublicKey:   string(publicKey),
		PrivateKey:  string(privateKey),
		Profile:     _req.Profile,
	}

	err = dao.Insert(space)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Uuid = space.UUID
	return nil
}

func (this *Space) Get(_ctx context.Context, _req *proto.SpaceGetRequest, _rsp *proto.SpaceGetResponse) error {
	logger.Infof("Received Space.Get, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewSpaceDAO(nil)

	space, err := dao.Get(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Space = &proto.SpaceEntity{
		Uuid:        space.UUID,
		Name:        space.Name,
		SpaceKey:    space.SpaceKey,
		SpaceSecret: space.SpaceSecret,
		PublicKey:   space.PublicKey,
		PrivateKey:  space.PrivateKey,
		Profile:     space.Profile,
		CreatedAt:   space.CreatedAt.Unix(),
	}
	return nil
}

func (this *Space) Update(_ctx context.Context, _req *proto.SpaceUpdateRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Space.Update, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewSpaceDAO(nil)

	space := &model.Space{
		UUID:    _req.Uuid,
		Profile: _req.Profile,
	}

	err := dao.Update(space)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Uuid = _req.Uuid
	return nil
}

func (this *Space) Search(_ctx context.Context, _req *proto.SpaceSearchRequest, _rsp *proto.SpaceListResponse) error {
	logger.Infof("Received Space.Search, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}

	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewSpaceDAO(nil)

	total, space, err := dao.Search(offset, count, _req.Name)
	// 数据库错误
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = total
	_rsp.Space = make([]*proto.SpaceEntity, len(space))
	for i, e := range space {
		_rsp.Space[i] = &proto.SpaceEntity{
			Uuid:        e.UUID,
			Name:        e.Name,
			SpaceKey:    e.SpaceKey,
			SpaceSecret: e.SpaceSecret,
			PublicKey:   e.PublicKey,
			PrivateKey:  e.PrivateKey,
			Profile:     e.Profile,
			CreatedAt:   e.CreatedAt.Unix(),
		}
	}

	return nil
}

func (this *Space) List(_ctx context.Context, _req *proto.SpaceListRequest, _rsp *proto.SpaceListResponse) error {
	logger.Infof("Received Space.List, request is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}

	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewSpaceDAO(nil)
	total, space, err := dao.List(offset, count)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = total
	_rsp.Space = make([]*proto.SpaceEntity, len(space))
	for i, e := range space {
		_rsp.Space[i] = &proto.SpaceEntity{
			Uuid:        e.UUID,
			Name:        e.Name,
			SpaceKey:    e.SpaceKey,
			SpaceSecret: e.SpaceSecret,
			PublicKey:   e.PublicKey,
			PrivateKey:  e.PrivateKey,
			Profile:     e.Profile,
			CreatedAt:   e.CreatedAt.Unix(),
		}
	}
	return nil
}
