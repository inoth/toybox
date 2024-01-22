package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func httpGet(url string, params map[string]string, token string, headers map[string]string) ([]byte, error) {
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

func httpGetWith[T interface{} | map[string]interface{}](url string, params map[string]string, token string, headers map[string]string) (T, error) {
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

func httpPost(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
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

func httpPostWith[T interface{} | map[string]interface{}](url string, params interface{}, token string, headers map[string]string) (T, error) {
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

func httpPut(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
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

func httpPutWith[T interface{} | map[string]interface{}](url string, params interface{}, token string, headers map[string]string) (T, error) {
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

func httpDelete(url string, params interface{}, token string, headers map[string]string) ([]byte, error) {
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

func httpDeleteWith[T interface{} | map[string]interface{}](url string, params interface{}, token string, headers map[string]string) (T, error) {
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
