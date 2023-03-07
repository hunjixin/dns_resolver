// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	key      = flag.String("key", "", "")
	secret   = flag.String("secret", "", "")
	domain   = flag.String("domain", "", "")
	password = flag.String("password", "", "")
)

type TYGWInfo struct {
	DevType    string
	LANIP      string
	LANIPv6    string
	MAC        string
	ProductCls string
	ProductSN  string
	SWVer      string
	WANIP      string
	WANIPv6    string
	ssid2g     string
	ssid5g     string
	wanAcnt    string
	wifiOnOff  string
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// create a custom error to know if a redirect happened
var RedirectAttemptedError = errors.New("redirect")

func loginTy() ([]*http.Cookie, error) {
	client := &http.Client{}
	// return the error, so client won't attempt redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return RedirectAttemptedError
	}

	req, err := http.NewRequest("POST", "http://192.168.1.1/cgi-bin/luci/", nil)
	if err != nil {
		return nil, err
	}
	/*

		POST http://192.168.1.1/cgi-bin/luci HTTP/1.1
		Host: 192.168.1.1
		Connection: keep-alive
		Content-Length: 28
		Cache-Control: max-age=0
		Upgrade-Insecure-Requests: 1
		Origin: http://192.168.1.1
		Content-Type: application/x-www-form-urlencoded
		User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36
		Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,;q=0.8,application/signed-exchange;v=b3;q=0.7
		Referer: http://192.168.1.1/cgi-bin/luci
		Accept-Encoding: gzip, deflate
		Accept-Language: zh-CN,zh;q=0.9

		username=useradmin&psd=gk4md
	*/
	req.Form = url.Values{}
	req.Form.Add("username", "useradmin")
	req.Form.Add("psd", *password)
	req.Header = http.Header{
		"Accept":                    []string{`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`},
		"Accept-Encoding":           []string{"gzip, deflate"},
		"Accept-Language":           []string{`zh-CN,zh;q=0.9`},
		"Cache-Control":             []string{`max-age=0`},
		"Connection":                []string{`keep-alive`},
		"Content-Length":            []string{`28`},
		"Content-Type":              []string{`application/x-www-form-urlencoded`},
		"Host":                      []string{`192.168.1.1`},
		"Origin":                    []string{`http://192.168.1.1`},
		"Referer":                   []string{`http://192.168.1.1/cgi-bin/luci`},
		"Upgrade-Insecure-Requests": []string{`1`},
		"User-Agent":                []string{`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("login tianyi fail")
	}

	return resp.Cookies(), nil
}

func getTyGatewayIp() (string, error) {
	cookies, err := loginTy()
	if err != nil {
		return "", err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}

	tyUrl, err := url.Parse("http:://192.168.1.1：80")
	if err != nil {
		return "", err
	}
	jar.SetCookies(tyUrl, cookies)

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", "http://192.168.1.1/cgi-bin/luci/admin/settings/gwinfo?get=part", nil)
	if err != nil {
		return "", err
	}

	/*
		Accept-Encoding: gzip, deflate
		Accept-Language: zh-CN,zh;q=0.9
		Connection: keep-alive
		Cookie: sysauth=ff61200a4a4330ad509858173618e1c2
		Host: 192.168.1.1
		Referer: http://192.168.1.1/cgi-bin/luci/admin/settings/info
		User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36
	*/

	req.Header = http.Header{
		"Accept":          []string{`*/*`},
		"Accept-Encoding": []string{"gzip, deflate"},
		"Accept-Language": []string{`zh-CN,zh;q=0.9`},
		"Connection":      []string{`keep-alive`},
		"Host":            []string{`192.168.1.1`},
		"Cookie":          []string{fmt.Sprintf(" sysauth=%s", cookies[0].Value)},
		"Referer":         []string{`http://192.168.1.1/cgi-bin/luci/admin/settings/info`},
		"User-Agent":      []string{`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("get gwinfo fail")
	}

	gwInfo := &TYGWInfo{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, gwInfo)
	if err != nil {
		return "", err
	}
	return gwInfo.LANIP, nil
}

func run(ctx context.Context) error {
	client, err := CreateClient(key, secret)
	if err != nil {
		return err
	}

	records, err := client.DescribeDomainRecords(&alidns20150109.DescribeDomainRecordsRequest{
		DomainName: domain,
	})
	if err != nil {
		return err
	}
	preIp := ""
	if len(records.Body.DomainRecords.Record) > 0 {
		preIp = *records.Body.DomainRecords.Record[0].Value
	}

	for {
		ip, err := getTyGatewayIp()
		if err != nil {
			log.Println("get gateway ", err)
			goto SLEEP
		}

		if ip != preIp {
			_, err = client.AddDomainRecord(&alidns20150109.AddDomainRecordRequest{
				DomainName:   domain,
				Lang:         tea.String("en"),
				RR:           tea.String("@"),
				Type:         tea.String("A"),
				UserClientIp: tea.String(ip),
				Value:        tea.String(ip),
			})
			if err != nil {
				log.Println("request error", err)
				goto SLEEP
			}
			preIp = ip
			goto SLEEP
		}

	SLEEP:
		time.Sleep(time.Second * 10)
	}
}

func main() {
	flag.Parse()
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}
