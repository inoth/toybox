package httptool

import (
	"fmt"
	"io"

	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/utils"

	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpGET(traceId, urlString string, urlParams url.Values, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	urlString = AddGetDataToUrl(urlString, urlParams)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, traceId)
	resp, err := client.Do(req)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"result":    utils.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	logger.Zap.Info(fmt.Sprintf("%+v", map[string]interface{}{
		"trace_id":  traceId,
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "GET",
		"args":      urlParams,
		"result":    utils.Substr(string(body), 0, 1024),
	}))
	return resp, body, nil
}

func HttpPOST(traceId, urlString string, urlParams url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}
	urlParamEncode := urlParams.Encode()
	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParamEncode))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, traceId)
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      utils.Substr(urlParamEncode, 0, 1024),
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      utils.Substr(urlParamEncode, 0, 1024),
			"result":    utils.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	logger.Zap.Info(fmt.Sprintf("%+v", map[string]interface{}{
		"trace_id":  traceId,
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      utils.Substr(urlParamEncode, 0, 1024),
		"result":    utils.Substr(string(body), 0, 1024),
	}))
	return resp, body, nil
}

func HttpJSON(traceId, urlString string, jsonContent string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(jsonContent))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, traceId)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      utils.Substr(jsonContent, 0, 1024),
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Zap.Warn(fmt.Sprintf("%+v", map[string]interface{}{
			"trace_id":  traceId,
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      utils.Substr(jsonContent, 0, 1024),
			"result":    utils.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		}))
		return nil, nil, err
	}
	logger.Zap.Info(fmt.Sprintf("%+v", map[string]interface{}{
		"trace_id":  traceId,
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      utils.Substr(jsonContent, 0, 1024),
		"result":    utils.Substr(string(body), 0, 1024),
	}))
	return resp, body, nil
}

func AddGetDataToUrl(urlString string, data url.Values) string {
	if strings.Contains(urlString, "?") {
		urlString = urlString + "&"
	} else {
		urlString = urlString + "?"
	}
	return fmt.Sprintf("%s%s", urlString, data.Encode())
}

func addTrace2Header(request *http.Request, traceId string) *http.Request {
	if traceId != "" {
		request.Header.Set("didi-header-rid", traceId)
	}
	return request
}
