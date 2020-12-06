package oss

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/collapsar/src/services/qiniu"
	"github.com/linshenqi/sptty"
)

var endpoint = "ashibro"

func getService() *Service {
	oss := Service{
		cfg: Config{Endpoints: map[string]base.Endpoint{
			endpoint: {
				Provider:  base.Qiniu,
				AppKey:    "yvgHtEuRewtH09atitGpYaq1oIW2NIPQYyzKgaqv",
				AppSecret: "WZaqBuONRJv8y4y-ooSxozc5XiRtLsj5Akwooohs",
				Bucket:    "ashibro",
				Zone:      "huadong",
			},
		}},

		providers: map[string]base.IOss{
			base.Qiniu: &qiniu.Oss{},
		},
	}

	oss.initProviders()
	return &oss
}

func getUrlImg(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func TestService(t *testing.T) {
	oss := getService()

	file := sptty.RandomFilename("wef")
	content, err := getUrlImg("http://thirdwx.qlogo.cn/mmopen/vi_32/E7XLbDS0gRJibYGpzxcEwXibyTwQAAHX9Koia7oln1821c8Djkibtpf6O20J3nacpnb0pg1UmtpdfDznHYBZvL78kw/132")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := oss.Upload(endpoint, file, content); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := oss.Delete(endpoint, file); err != nil {
		fmt.Println(err.Error())
		return
	}
}
