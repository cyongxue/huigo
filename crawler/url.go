// @Description：爬虫相关的基础结构
//      本文件用于对于url处理的封装
package crawler

import (
    "strings"
    "net/url"
)

type Url struct {
    StrUrl          string

    UnicodeUrl      string
    Change          bool
    Encoding        string

    URL             *url.URL             // 解析后的url
}

/**
 @Description：初始化
 @Param:
 @Return：
 */
func (u *Url) Init() error {
    // http头处理
    if strings.HasPrefix(u.StrUrl, "https://") == false && strings.HasPrefix(u.StrUrl, "http://") == false {
        u.StrUrl = "http://" + u.StrUrl
    }

    urlRes, err := url.Parse(u.StrUrl)
    if err != nil {
        return err
    }
    u.URL = urlRes
    return nil
}

/**
 @Description：获取url中的文件
 @Param:
 @Return：
 */
func (u *Url) GetFileName() string {
    return u.URL.Path[strings.LastIndex(u.URL.Path, "/") + 1:]
}

/**
 @Description：获取扩展名
 @Param:
 @Return：
 */
func (u *Url) GetFileExt() string {
    fileName := u.GetFileName()
    ext := fileName[strings.LastIndex(fileName, ".") + 1:]
    if ext == fileName {
        return ""
    }
    return ext
}

