// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	query "h12.io/html-query"
)

var (
	key      = flag.String("key", "", "")
	secret   = flag.String("secret", "", "")
	domain   = flag.String("domain", "", "")
	interval = flag.String("interval", "10m", "")

	iprex = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
)

func createClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// resolvePublicIp resolve public ip by ident me
func resolvePublicIpByIdentMe(ctx context.Context) (string, error) {
	resp, err := http.Get("https://ident.me")
	if err != nil {
		return "", err
	}

	root, err := query.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	return root.Body().InternalNode().FirstChild.Data, nil
}

// resolvePublicIp resolve public ip by ident me
func resolvePublicIpByNetCN(ctx context.Context) (string, error) {
	resp, err := http.Get("http://www.net.cn/static/customercare/yourip.asp")
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	submatchall := iprex.FindAllString(string(data), 1)
	if len(submatchall) == 0 {
		return "", fmt.Errorf("www.net.cn maybe fail")
	}
	return submatchall[0], nil
}

// run get ip and set to aliyun dns server
func run(ctx context.Context) error {
	client, err := createClient(key, secret)
	if err != nil {
		return err
	}

	dur, err := time.ParseDuration(*interval)
	if err != nil {
		return err
	}

	records, err := client.DescribeDomainRecords(&alidns20150109.DescribeDomainRecordsRequest{
		DomainName: domain,
	})
	if err != nil {
		return err
	}
	recordId := ""
	preIP := ""
	if len(records.Body.DomainRecords.Record) > 0 {
		preIP = *records.Body.DomainRecords.Record[0].Value
		recordId = *records.Body.DomainRecords.Record[0].RecordId
	}

	for {
		ip, err := resolvePublicIpByNetCN(ctx)
		if err != nil {
			log.Println("get public ips ", err)
			goto SLEEP
		}

		if ip != preIP {
			log.Printf("ip changed old(%s) new(%s) try to change dns", preIP, ip)
			if len(recordId) == 0 {
				adrResp, err := client.AddDomainRecord(&alidns20150109.AddDomainRecordRequest{
					DomainName:   domain,
					Lang:         tea.String("en"),
					RR:           tea.String("@"),
					Type:         tea.String("A"),
					UserClientIp: tea.String(ip),
					Value:        tea.String(ip),
				})
				if err != nil {
					log.Println("add record error", err)
					goto SLEEP
				}
				recordId = *adrResp.Body.RecordId
			} else {
				_, err = client.UpdateDomainRecord(&alidns20150109.UpdateDomainRecordRequest{
					Lang:         tea.String("en"),
					RR:           tea.String("@"),
					RecordId:     tea.String(recordId),
					Type:         tea.String("A"),
					UserClientIp: tea.String(ip),
					Value:        tea.String(ip),
				})
				if err != nil {
					log.Println("update record error", err)
					goto SLEEP
				}
			}

			preIP = ip
			log.Println("dns updated successfully")
			goto SLEEP
		}

	SLEEP:
		time.Sleep(dur)
	}
}

func main() {
	os.Setenv("HTTPS_PROXY", "")
	os.Setenv("HTTP_PROXY", "")
	flag.Parse()
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}
