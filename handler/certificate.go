package handler

import (
	"context"

	"ogm-msa-license/model"

	"github.com/asim/go-micro/v3/logger"

	proto "github.com/xtech-cloud/ogm-msp-license/proto/license"
)

type Certificate struct{}

func (this *Certificate) Get(_ctx context.Context, _req *proto.CerGetRequest, _rsp *proto.CerGetResponse) error {
	logger.Infof("Received Certificate.Fetch, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewCertificateDAO(nil)

	cer, err := dao.Get(_req.Uuid)
	// 数据库错误
	if nil != err || nil == cer {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "cer not found"
		return nil
	}

	_rsp.Cer = &proto.CertificateEntity{
		Uuid:      cer.UUID,
		Space:     cer.Space,
		Number:    cer.Key,
		Consumer:  cer.Consumer,
		Content:   cer.Content,
		CreatedAt: cer.CreatedAt.Unix(),
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

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}
	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewCertificateDAO(nil)

	total, cers, err := dao.List(offset, count, _req.Space)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = total
	_rsp.Cer = make([]*proto.CertificateEntity, len(cers))
	for i, cer := range cers {
		_rsp.Cer[i] = &proto.CertificateEntity{
			Uuid:      cer.UUID,
			Space:     cer.Space,
			Number:    cer.Key,
			Consumer:  cer.Consumer,
			CreatedAt: cer.CreatedAt.Unix(),
		}
	}
	return nil
}

func (this *Certificate) Search(_ctx context.Context, _req *proto.CerSearchRequest, _rsp *proto.CerListResponse) error {
	logger.Infof("Received Certificate.Search, request is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Space {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "space is required"
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

	dao := model.NewCertificateDAO(nil)

	total, cers, err := dao.Search(offset, count, _req.Space, _req.Number, _req.Consumer)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = total
	_rsp.Cer = make([]*proto.CertificateEntity, len(cers))
	for i, cer := range cers {
		_rsp.Cer[i] = &proto.CertificateEntity{
			Uuid:      cer.UUID,
			Space:     cer.Space,
			Number:    cer.Key,
			Consumer:  cer.Consumer,
			CreatedAt: cer.CreatedAt.Unix(),
		}
	}
	return nil
}
