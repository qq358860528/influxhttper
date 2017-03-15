package main

/*
* 358860528@qq.com
* influxdb http写influxdb
* 这样才支持 influxql 的write Data
 */

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main() {
	Post()
	fmt.Println("inflxdb ok")
}

func Post() {

	//influxdb 地址
	host, err := url.Parse("http://localhost:8086")
	if err != nil {
		return
	}

	//批量命令
	var b bytes.Buffer
	//第一条错误 后面的命令依然会写入 只抛弃掉错误的语句
	b.WriteString(`cpu_load_short,host=server02 value=0.67,value1='1'`)
	b.WriteString("\n")
	//库名,key1,key2 field1,field2 time
	b.WriteString("cpu_load_short,host=server02,region=us-west value=0.55,value1=\"1\" 1422568543702900257\n")
	b.WriteString("cpu_load_short,direction=in,host=server01,region=us-west value=2.0,value1=\"xxyy\" 1422568543702900257\n")

	host.Path = "write"
	req, err := http.NewRequest("POST", host.String(), &b)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "")
	req.Header.Set("User-Agent", "influxgoclient")

	//账号密码
	//req.SetBasicAuth(username,password)
	params := req.URL.Query()

	//库名
	params.Set("db", "billog")
	//params.Set("rp", bp.RetentionPolicy())
	params.Set("precision", "ns")
	//params.Set("consistency", bp.WriteConsistency())

	req.URL.RawQuery = params.Encode()

	//设置超时时间 跟 https 证书
	httpClient := &http.Client{Timeout: 5 * time.Second}
	//&http.Client{Timeout: 5 * time.Second,Transport: nil}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	//读完可以保持长链接
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	//失败日志
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))

	}

}
