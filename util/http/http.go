package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
)

type RespData interface {
	any | map[string]any
}

type RequestOption struct {
	Token      string
	CaCertPath string
	Headers    map[string]string
}

func HttpGet(url string, params map[string]string, opts ...RequestOption) ([]byte, error) {
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		Get(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func HttpGetWith[T RespData](url string, params map[string]string, opts ...RequestOption) (T, error) {
	var res T
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		SetResult(&res).
		Get(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func HttpPost(url string, params any, opts ...RequestOption) ([]byte, error) {
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		Post(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func HttpPostWith[T RespData](url string, params any, opts ...RequestOption) (T, error) {
	var res T
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		SetResult(&res).
		Post(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func HttpPut(url string, params any, opts ...RequestOption) ([]byte, error) {
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		Put(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func HttpPutWith[T RespData](url string, params any, opts ...RequestOption) (T, error) {
	var res T
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		SetResult(&res).
		Put(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func HttpDelete(url string, params any, opts ...RequestOption) ([]byte, error) {
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		Delete(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func HttpDeleteWith[T RespData](url string, params any, opts ...RequestOption) (T, error) {
	var res T
	client := resty.New()
	opt := util.First(RequestOption{}, opts)
	if opt.CaCertPath != "" {
		client.SetRootCertificate(opt.CaCertPath)
	}
	if opt.Token != "" {
		client.SetAuthToken(opt.Token)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(opt.Headers).
		SetResult(&res).
		Delete(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}
