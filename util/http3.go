package util

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

func getHttp3RoundTripper(caCertPath string) (*http3.RoundTripper, error) {
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

func Http3Get(url string, params map[string]string, token string, headers map[string]string, caCertPath string) ([]byte, error) {
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func Http3GetWith[T RespData](url string, params map[string]string, token string, headers map[string]string, caCertPath string) (T, error) {
	var res T
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		SetResult(&res).
		Get(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func Http3Post(url string, params any, token string, headers map[string]string, caCertPath string) ([]byte, error) {
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		Post(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func Http3PostWith[T RespData](url string, params any, token string, headers map[string]string, caCertPath string) (T, error) {
	var res T
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		SetResult(&res).
		Post(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func Http3Put(url string, params any, token string, headers map[string]string, caCertPath string) ([]byte, error) {
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		Put(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func Http3PutWith[T RespData](url string, params any, token string, headers map[string]string, caCertPath string) (T, error) {
	var res T
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		SetResult(&res).
		Put(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}

func Http3Delete(url string, params any, token string, headers map[string]string, caCertPath string) ([]byte, error) {
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return nil, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		Delete(url)
	if err != nil {
		return nil, errors.Wrap(err, resp.Status())
	}
	return resp.Body(), nil
}

func Http3DeleteWith[T RespData](url string, params any, token string, headers map[string]string, caCertPath string) (T, error) {
	var res T
	client := resty.New()
	roundTripper, err := getHttp3RoundTripper(caCertPath)
	if err != nil {
		return res, err
	}
	defer roundTripper.Close()

	resp, err := client.SetTransport(roundTripper).R().
		EnableTrace().
		SetBody(params).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetAuthToken(token).
		SetResult(&res).
		Delete(url)
	if err != nil {
		return res, errors.Wrap(err, resp.Status())
	}
	return res, nil
}
