// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
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
func resolvePublicIp(ctx context.Context) (string, error) {
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
	preIP := ""
	if len(records.Body.DomainRecords.Record) > 0 {
		preIP = *records.Body.DomainRecords.Record[0].Value
	}

	for {
		ip, err := resolvePublicIp(ctx)
		if err != nil {
			log.Println("get public ips ", err)
			goto SLEEP
		}

		if ip != preIP {
			log.Printf("ip changed old(%s) new(%s) try to change dns", preIP, ip)
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
			preIP = ip
			log.Println("dns updated successfully")
			goto SLEEP
		}

	SLEEP:
		time.Sleep(dur)
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
