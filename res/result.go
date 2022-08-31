package res

var (
	// 成功
	SUCCESS = 200
	// 已经处于处理状态，无需处理
	SUCCESSED = 201
	// 参数错误
	PARAMETERERR = 400
	// 中断执行
	PROHIBITED = 402
	// 无效授权
	UNAUTHORIZATION = 403
	// 空
	NOTFOUND = 404
	// 失败
	FAILED = 500
)

type IResult interface {
	R()
}

type ApiResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ApiResultWithData struct {
	ApiResult
	Data interface{} `json:"data"`
}

func (r *ApiResultWithData) R() {}

func (r *ApiResult) R() {}

func result(code int, msg string, data ...interface{}) IResult {
	if len(data) > 0 {
		r := &ApiResultWithData{
			ApiResult: ApiResult{
				Code: code,
				Msg:  msg},
			Data: data[0]}
		return r
	}
	r := &ApiResult{Code: code, Msg: msg}
	return r
}

func resultOK(msg ...string) IResult {
	var r *ApiResult
	if len(msg) > 0 {
		r = &ApiResult{
			Code: SUCCESS,
			Msg:  msg[0],
		}
		return r
	}
	r = &ApiResult{Code: SUCCESS, Msg: "ok"}
	return r
}

func ok(msg string, data ...interface{}) IResult {
	if len(data) > 0 {
		r := &ApiResultWithData{
			ApiResult: ApiResult{
				Code: SUCCESS,
				Msg:  msg},
			Data: data[0]}
		return r
	}
	r := &ApiResult{Code: SUCCESS, Msg: msg}
	return r
}

func err(msg string) IResult {
	return &ApiResult{Code: FAILED, Msg: msg}
}

func notFound(msg string) IResult {
	return &ApiResult{Code: NOTFOUND, Msg: msg}
}

func paramErr(msg string) IResult {
	return &ApiResult{Code: PARAMETERERR, Msg: msg}
}

func unauthrized(msg string) IResult {
	return &ApiResult{Code: UNAUTHORIZATION, Msg: msg}
}
