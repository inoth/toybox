package resty

import (
	"testing"

	"github.com/inoth/toybox/resty/v3"
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
		Data any    `json:"data"`
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
				url: "http://localhost:9060/api/sayhi/httpget",
			},
			want: `{"trace_id":"7f8a87d5b827491db9b9e7000089b2ad","ret":0,"msg":"hello httpget"}`,
		},
		{
			name: "HttpGetWith",
			args: args{
				url: "http://localhost:9060/api/sayhi/httpget",
			},
			want: "hello httpget",
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
				require.Equal(t, len(tt.want), len(got), "HttpGet() = %v, want %v", string(got), tt.want)
			case "HttpGetWith":
				got, err := HttpGetWith[resp](tt.args.url, tt.args.params, tt.args.opts...)
				if err != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, got.Msg, "HttpGet() = %v, want %v", got.Msg, tt.want)
				// case "HttpPost":
				// case "HttpPostWith":
			}
		})
	}
}

func TestHttp3Get(t *testing.T) {
	type resp struct {
		Ret  int    `json:"ret"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	type args struct {
		url    string
		params map[string]string
		opts   []resty.RequestOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HttpGet",
			args: args{
				url: "https://localhost:9062/api/sayhi/httpget",
				opts: []resty.RequestOption{
					{CaCertRaw: []byte(ca)},
				},
			},
			want: `{"trace_id":"7f8a87d5b827491db9b9e7000089b2ad","ret":0,"msg":"hello httpget"}`,
		},
		{
			name: "HttpGetWith",
			args: args{
				url: "https://localhost:9062/api/sayhi/httpget",
				opts: []resty.RequestOption{
					{CaCertRaw: []byte(ca)},
				},
			},
			want: "hello httpget",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "HttpGet":
				got, err := resty.HttpGet(tt.args.url, tt.args.params, tt.args.opts...)
				if err != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, len(tt.want), len(got), "HttpGet() = %v, want %v", string(got), tt.want)
			case "HttpGetWith":
				got, err := resty.HttpGetWith[resp](tt.args.url, tt.args.params, tt.args.opts...)
				if err != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, got.Msg, "HttpGet() = %v, want %v", got.Msg, tt.want)
				// case "HttpPost":
				// case "HttpPostWith":
			}
		})
	}
}
