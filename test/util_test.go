package test

import (
	"crypto"
	"fmt"
	"testing"
	"time"

	"github.com/inoth/ino-toybox/components/license"
	"github.com/inoth/ino-toybox/util"
	"github.com/inoth/ino-toybox/util/encryption"
)

const (
	public = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAnhxzyi5OXN76nf+UtsApquhF1Ju40FvTV/IlnfwmzXaWnXoDeeLn
oEJ6I+AmsEMqmWtX2YQEEl9H1JwHXAsGhqHbAo9BRQcefdZHG+pvHonp9+c7uJbU
oFEWdSrlI2yD8E4HJCijsgWvrrUmEsqopV8M/8HpCJl3dQGPJOoRTTU45Vq1muRa
7vYSpmzjbJYyxTqful9hLtj8BW/IS6Sd2EjkbBWSp0mFc2thY4osDYKoHqp0rjIq
+Ou/bKIoYTlyHEdmqQmSg9GEfbp/S+kwyxNAuvkLzfD/mPrhzrb7gxP4LX9LdrlP
qvA3pNRGLed7XgVujwf4QYaP+GLanG2XsQIDAQAB
-----END RSA PUBLIC KEY-----`
	private = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAnhxzyi5OXN76nf+UtsApquhF1Ju40FvTV/IlnfwmzXaWnXoD
eeLnoEJ6I+AmsEMqmWtX2YQEEl9H1JwHXAsGhqHbAo9BRQcefdZHG+pvHonp9+c7
uJbUoFEWdSrlI2yD8E4HJCijsgWvrrUmEsqopV8M/8HpCJl3dQGPJOoRTTU45Vq1
muRa7vYSpmzjbJYyxTqful9hLtj8BW/IS6Sd2EjkbBWSp0mFc2thY4osDYKoHqp0
rjIq+Ou/bKIoYTlyHEdmqQmSg9GEfbp/S+kwyxNAuvkLzfD/mPrhzrb7gxP4LX9L
drlPqvA3pNRGLed7XgVujwf4QYaP+GLanG2XsQIDAQABAoIBAEs/8V/dUBBlU1PV
KxsbM/mSWIHKp1gLC/gEWCDrvC/3a9GBG5xsAJ9GZEkkymUDYofoDcSJT0LLNC2d
IOeOm8mByPFb6s2GiN2NGQTRO5eGPeNtmv/MUWAHl6+l/a4xXE4HZOCxss7sY+O7
dWZmK4OhTFeBG36rQ65dUppYCGn81g3XUEkopF5lLyfCnEyKqZWvfdcedD9nDHTs
0TvQbjm5rSAllnpWKC6cuyqhjNi7HsAtP3axRAb+iXi0RHSf8srIQSv9Z5NyREuc
7Y8H024YXGoUtGaQeL/g0+b+1zLhrvzTU0K86Wtl+LuybTePfxSpS6gxJe/2YsqZ
CWJyzxECgYEAzYe/WzdHvkvU40fIpCpgqcwP7R2r8pHAXYi8j36B73T2gFUeAsB5
oOietLf5+TBPxCD8oIvS2LR4LffBegjRRZAnwt11xPMlgXklBvlU+JyGjoSHPWJj
93Y5c8tZ/c6GjmG7hX9fHQVHmyhrqk9ifMVl8ZMj2Bg+2GJfTCzU6ZUCgYEAxO/L
tepMH/T0LcckL2A5aG+ZrvA3zLBsIarwWZn7HM9yvOmlySQsn/69QzvUKTnk/1xJ
beb14+H09J4W8pujsPxYI/6KilBqrqdvZvjmlgscgPKZtQxtVPR+T4YT83Ci5szW
XNHJA5+yYIqITHhSdsLGC18txTrAUfHKoRy5Rq0CgYBTJoJCUwERefhs4xPHZuWo
jEg9M+3muxTKQpGWtCW5TOaVUNpNXrVWZgYfMvdM20DKJlZOVYM97PVaE4wQ5RRV
QlbzvUjyHzSjRvG1+pVn51uAuRlFulKbQRdJQ5Hq3u0NGXkWL0u5n/MyUI4OXwOH
Ww09SLwNpvF19YZ8eP7CaQKBgFEsWfYYpdoCOGdqDuMsMV13qovt3cIT8e4KrFjy
XAvrAesWD0ySCYbFFDPTREbd4yLSYj3XlgChETuGsgS73EPGL3pen7IVJXPp9cQm
0byExfHsjSiP/7ylri6PIEgWZD7nrW/C1K0WtQqP71A9xBfJfqIPUClcZwsfs5qm
4UNdAoGABoeWU8GU7Xvj1enK6jqkWuPHtWnmczY2w0Egqp53FLt2e/7clHAbe9Q1
ZM4p2RdR4StQotenspVLZos6WPWvDzfV3DTtEpHx+urOjsPL1/mLWPvPh0RDklZU
jfntFYom1CHuOYzT/cETKs/xOpcYz2Luw1gfOnAf0+PscOVyZmM=
-----END RSA PRIVATE KEY-----`
)

func TestRandString(t *testing.T) {
	i := 255
	h := fmt.Sprintf("%x", i)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i, h)
	h = fmt.Sprintf("%X", i)
	fmt.Printf("HEX conv of '%d' is '%s'\n", i, h)
	// println(util.RandString(18))
}

func TestAes(t *testing.T) {
	// data := `{"appId": "cn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemocn.fdm.offlicensedemo","issuedTime": 1595951714,"notBefore": 1538671712,"notAfter": 1640966400,"customerInfo": "亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a亚马逊公司a"}`
	lic := license.License{
		LegalMachine: map[string]struct{}{
			"machineA": {},
			"machineB": {},
			"machineC": {},
		},
		Subscribes: map[string]license.Subscribe{
			"nginx": {
				AppName:   "ngxin",
				Name:      "nginx",
				NodeLimit: 10,
				State:     true,
				Expire:    time.Now().Add(time.Hour * 720),
			},
		},
	}
	key := util.RandString(16)
	// randStr := util.RandString(32)
	encrypt, err := encryption.AesEcrypt(lic.String(), []byte(key))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Printf("data: %v\n", string(lic.String()))
	fmt.Printf("key: %v\n", key)
	fmt.Printf("aes: %s\n", encrypt)
	fmt.Printf("len: %v\n", len(encrypt))
	fmt.Printf("hex: %s\n", fmt.Sprintf("%x", len(encrypt)))

	// f, err := os.Create(key + ".license")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer f.Close()
	// f.WriteString(encrypt)

	// sign
	sign, err := encryption.SignWithRSA([]byte(encrypt), []byte(private), crypto.SHA256)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("sign: %v\n", sign)

	// verify
	check := encryption.VerifyWithRSA([]byte(encrypt), []byte(public), sign, crypto.SHA256)
	fmt.Printf("verify: %v\n", check)

	res, _ := encryption.AesDeCrypt(encrypt, []byte(key))
	fmt.Printf("data: %v\n", string(res))
}

func TestMakeLicense(t *testing.T) {
	lic := license.License{
		LegalMachine: map[string]struct{}{
			"machineA": {},
			"machineB": {},
			"machineC": {},
		},
		Subscribes: map[string]license.Subscribe{
			"nginx": {
				AppName:   "ngxin",
				Name:      "nginx",
				NodeLimit: 10,
				State:     true,
				Expire:    time.Now().Add(time.Hour * 720),
			},
		},
	}
	key := util.RandString(16)
	license.MakeLicense(key, []byte(private), lic)
}
