package hwebdir

import (
    "net/http"
    "strings"
    "io/ioutil"
    "fmt"
    "time"
)

type HttpRequest struct {
    Domain          string          // domain或者ip
    Port            string          // 端口号

    ResponseBody    []byte
    Response        *http.Response
}

func (h *HttpRequest) HttpDo(url string) error {
    client := &http.Client{
        Timeout: time.Second * 5,                   // 设置5秒超时
    }

    var reqUrl string
    if h.Port == "443" {
        reqUrl = fmt.Sprintf("https://%s/%s", h.Domain, url)
    } else {
        reqUrl = fmt.Sprintf("http://%s:%s/%s", h.Domain, h.Port, url)
    }
    req, err := http.NewRequest("GET", reqUrl, strings.NewReader("name=cjb"))
    if err != nil {
        return err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36 (Eagle)")
    req.Header.Set("Accept", "application/json, text/plain, */*")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    h.ResponseBody = body
    h.Response = resp

    return nil
}