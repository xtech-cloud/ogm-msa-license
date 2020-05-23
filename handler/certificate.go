package handler

import (
	"context"

	"omo-msa-license/model"

	"github.com/micro/go-micro/v2/logger"

	proto "github.com/xtech-cloud/omo-msp-license/proto/license"
)

type Certificate struct{}

func (this *Certificate) Fetch(_ctx context.Context, _req *proto.CerFetchRequest, _rsp *proto.CerFetchResponse) error {
	logger.Infof("Received Certificate.Fetch, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uid is required"
		return nil
	}

	dao := model.NewCertificateDAO()

	cer, err := dao.Find(_req.Uid)
	// 数据库错误
	if nil != err {
		return err
	}

	if "" == cer.UID {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "certificate not found"
		return nil
	}

	_rsp.Cer = &proto.CertificateEntity{
		Uid:       cer.UID,
		Space:     cer.Space,
		Number:    cer.Key,
		Consumer:  cer.Consumer,
		Content:   cer.Content,
		CreatedAt: cer.GModel.CreatedAt.Unix(),
	}

	return nil
}

func (this *Certificate) List(_ctx context.Context, _req *proto.CerListRequest, _rsp *proto.CerListResponse) error {
	logger.Infof("Received Certificate.List, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
		return nil
	}

	dao := model.NewCertificateDAO()

	count, err := dao.Count(model.CertificateQuery{
		Space: _req.Space,
	})
	// 数据库错误
	if nil != err {
		return err
	}

	cers, err := dao.List(_req.Offset, _req.Count, _req.Space)
	// 数据库错误
	if nil != err {
		return err
	}

	_rsp.Total = count
	_rsp.Cer = make([]*proto.CertificateEntity, len(cers))
	for i, cer := range cers {
		_rsp.Cer[i] = &proto.CertificateEntity{
			Uid:       cer.UID,
			Space:     cer.Space,
			Number:    cer.Key,
			Consumer:  cer.Consumer,
			CreatedAt: cer.GModel.CreatedAt.Unix(),
		}
	}
	return nil
}
