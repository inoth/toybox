package http3

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

type RespData interface {
	any | map[string]any
}

type RequestOption struct {
	Token      string
	CaCertPath string
	Headers    map[string]string
}

func getHttp3RoundTripper(caCertPath string) (*http3.RoundTripper, error) {
	if caCertPath == "" {
		return nil, fmt.Errorf("the ca file address is empty")
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	caCertRaw, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "load ca path failed.")
	}
	if ok := pool.AppendCertsFromPEM(caCertRaw); !ok {
		return nil, errors.Wrap(err, "Could not add root ceritificate to pool.")
	}

	return &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}, nil
}

func HttpGet(url string, params map[string]string, opts ...RequestOption) ([]byte, error) {
	client := resty.New()
	opt := util.First(RequestOption{}, opts)

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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

	roundTripper, err := getHttp3RoundTripper(opt.CaCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()
	client.SetTransport(roundTripper)

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
