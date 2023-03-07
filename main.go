// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"context"
	"flag"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"time"
)

var (
	key    = flag.String("key", "", "")
	secret = flag.String("secret", "", "")
	domain = flag.String("domain", "", "")
)

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
func getTyGatewayIp() (string, error) {
	return "1.1.1.1", nil
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
