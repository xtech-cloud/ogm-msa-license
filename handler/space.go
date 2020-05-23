package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"omo-msa-license/crypto"
	"omo-msa-license/model"

	"github.com/micro/go-micro/v2/logger"

	proto "github.com/xtech-cloud/omo-msp-license/proto/license"
)

type Space struct{}

func (this *Space) Create(_ctx context.Context, _req *proto.SpaceCreateRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Space.Create, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewSpaceDAO()

	// 账号存在检测
	exists, err := dao.Exists(_req.Name)
	// 数据库错误
	if nil != err {
		return err
	}

	if exists {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "space exists"
		return nil
	}

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
		return err
	}

	space := model.Space{
		Name:        _req.Name,
		SpaceKey:    spaceKey,
		SpaceSecret: spaceSecret,
		PublicKey:   string(publicKey),
		PrivateKey:  string(privateKey),
	}

	err = dao.Insert(space)
	if nil != err {
		return err
	}
	return nil
}

func (this *Space) Query(_ctx context.Context, _req *proto.SpaceQueryRequest, _rsp *proto.SpaceQueryResponse) error {
	logger.Infof("Received Space.Query , request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewSpaceDAO()

	space, err := dao.Find(_req.Name)
	// 数据库错误
	if nil != err {
		return err
	}

	if "" == space.Name {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "space not found"
		return nil
	}

	_rsp.SpaceKey = space.SpaceKey
	_rsp.SpaceSecret = space.SpaceSecret
	_rsp.PublicKey = space.PublicKey
	_rsp.PrivateKey = space.PrivateKey
	_rsp.Profile = space.Profile
	_rsp.CreatedAt = space.GModel.CreatedAt.Unix()
	_rsp.UpdatedAt = space.GModel.UpdatedAt.Unix()

	return nil
}

func (this *Space) List(_ctx context.Context, _req *proto.SpaceListRequest, _rsp *proto.SpaceListResponse) error {
	logger.Infof("Received Space.List, request is %v", _req)
	_rsp.Status = &proto.Status{}

	dao := model.NewSpaceDAO()

	spaces, err := dao.List()
	// 数据库错误
	if nil != err {
		return err
	}

	_rsp.Name = make([]string, len(spaces))
	for i, space := range spaces {
		_rsp.Name[i] = space.Name
	}
	return nil
}
