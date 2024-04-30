package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type RespData interface {
	any | map[string]any
}

func HttpGet(url string, params map[string]string, token string, headers map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
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

func HttpGetWith[T RespData](url string, params map[string]string, token string, headers map[string]string) (T, error) {
	var res T
	client := resty.New()
	resp, err := client.R().
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

func HttpPost(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
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

func HttpPostWith[T RespData](url string, params interface{}, token string, headers map[string]string) (T, error) {
	var res T
	client := resty.New()
	resp, err := client.R().
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

func HttpPut(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
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

func HttpPutWith[T RespData](url string, params interface{}, token string, headers map[string]string) (T, error) {
	var res T
	client := resty.New()
	resp, err := client.R().
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

func HttpDelete(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
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

func HttpDeleteWith[T RespData](url string, params interface{}, token string, headers map[string]string) (T, error) {
	var res T
	client := resty.New()
	resp, err := client.R().
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
