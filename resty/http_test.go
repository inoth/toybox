package resty

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ca = `-----BEGIN CERTIFICATE-----
MIICzDCCAbQCCQDA+rLymNnfJzANBgkqhkiG9w0BAQsFADAoMSYwJAYDVQQKDB1x
dWljLWdvIENlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMDA4MTgwOTIxMzVaFw0z
MDA4MTYwOTIxMzVaMCgxJjAkBgNVBAoMHXF1aWMtZ28gQ2VydGlmaWNhdGUgQXV0
aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1OcsYrVaSDfh
iDppl6oteVspOY3yFb96T9Y/biaGPJAkBO9VGKcqwOUPmUeiWpedRAUB9LE7Srs6
qBX4mnl90Icjp8jbIs5cPgIWLkIu8Qm549RghFzB3bn+EmCQSe4cxvyDMN3ndClp
3YMXpZgXWgJGiPOylVi/OwHDdWDBorw4hvry+6yDtpQo2TuI2A/xtxXPT7BgsEJD
WGffdgZOYXChcFA0c1XVLIYlu2w2JhxS8c2TUF6uSDlmcoONNKVoiNCuu1Z9MorS
Qmg7a2G7dSPu123KcTcSQFcmJrt+1G81gOBtHB69kacD8xDmgksj09h/ODPL/gIU
1ZcU2ci1/QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQB0Tb1JbLXp/BvWovSAhO/j
wG7UEaUA1rCtkDB+fV2HS9bxCbV5eErdg8AMHKgB51ygUrq95vm/baZmUILr84XK
uTEoxxrw5S9Z7SrhtbOpKCumoSeTsCPjDvCcwFExHv4XHFk+CPqZwbMHueVIMT0+
nGWss/KecCPdJLdnUgMRz0tIuXzkoRuOiUiZfUeyBNVNbDFSrLigYshTeAPGaYjX
CypoHxkeS93nWfOMUu8FTYLYkvGMU5i076zDoFGKJiEtbjSiNW+Hei7u2aSEuCzp
qyTKzYPWYffAq3MM2MKJgZdL04e9GEGeuce/qhM1o3q77aI/XJImwEDdut2LDec1
-----END CERTIFICATE-----`
)

func TestHttpGet(t *testing.T) {
	type resp struct {
		Ret  int    `json:"ret"`
		Msg  string `json:"msg"`
		Data struct {
			Id    int    `json:"id"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		} `json:"data"`
	}
	type args struct {
		url    string
		params map[string]string
		opts   []RequestOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HttpGet",
			args: args{
				url: "http://localhost:8080/test",
				params: map[string]string{
					"id":    "3",
					"phone": "18581619978",
					"email": "aaa@aaa.com",
				},
			},
			want: `{"data":{"id":3,"phone":"18581619978","email":"aaa@aaa.com"},"msg":"ok","ret":0}`,
		},
		{
			name: "HttpGetWith",
			args: args{
				url: "http://localhost:8080/test",
				params: map[string]string{
					"id":    "3",
					"phone": "18581619978",
					"email": "aaa@aaa.com",
				},
			},
			want: `id:3, phone:18581619978, email:aaa@aaa.com`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "HttpGet":
				got, err := HttpGet(tt.args.url, tt.args.params, tt.args.opts...)
				if err != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, string(got), "HttpGet() = %v, want %v", string(got), tt.want)
			case "HttpGetWith":
				got, err := HttpGetWith[resp](tt.args.url, tt.args.params, tt.args.opts...)
				if err != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				str := fmt.Sprintf("id:%d, phone:%s, email:%s", got.Data.Id, got.Data.Phone, got.Data.Email)
				require.Equal(t, tt.want, str, "HttpGet() = %v, want %v", str, tt.want)
				// case "HttpPost":
				// case "HttpPostWith":
			}
		})
	}
}

// func TestHttp3Get(t *testing.T) {
// 	type resp struct {
// 		Ret  int    `json:"ret"`
// 		Msg  string `json:"msg"`
// 		Data any    `json:"data"`
// 	}
// 	type args struct {
// 		url    string
// 		params map[string]string
// 		opts   []resty.RequestOption
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "HttpGetWith",
// 			args: args{
// 				url: "https://localhost:9062/api/sayhi/httpget",
// 				opts: []resty.RequestOption{
// 					{CaCertRaw: []byte(ca)},
// 				},
// 			},
// 			want: "hello httpget",
// 		},
// 		{
// 			name: "HttpPostWith",
// 			args: args{
// 				url: "https://localhost:9062/api/hi/httpget1",
// 				opts: []resty.RequestOption{
// 					{CaCertRaw: []byte(ca)},
// 				},
// 			},
// 			want: "hello httpget1",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			switch tt.name {
// 			case "HttpGetWith":
// 				got, err := resty.HttpGetWith[resp](tt.args.url, tt.args.params, tt.args.opts...)
// 				if err != nil {
// 					require.Error(t, err)
// 				} else {
// 					require.NoError(t, err)
// 				}
// 				require.Equal(t, tt.want, got.Msg, "HttpGet() = %v, want %v", got.Msg, tt.want)

// 			case "HttpPostWith":
// 				got, err := resty.HttpPostWith[resp](tt.args.url, tt.args.params, tt.args.opts...)
// 				if err != nil {
// 					require.Error(t, err)
// 				} else {
// 					require.NoError(t, err)
// 				}
// 				require.Equal(t, tt.want, got.Msg, "HttpGet() = %v, want %v", got.Msg, tt.want)
// 			}
// 		})
// 	}
// }
