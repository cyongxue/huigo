package hwebdir

import (
    "fmt"
    "strings"
)

// @Description: 404的判断预定义url
var url404 = []string {
    "not_exits_by_wilson_adfxvsdfa",
    "not_exits_by_wilson_iuwertm",
}

// @Description: 备份文件因子
var backFactor = []string {
    "EXT.bak", "EXT~", "EXT.swp", "EXT.log", "EXT.zip", "EXT.tar.gz", "EXT.tar.bz2",
    "EXT.tar",
    "EXT.rar", "EXT.bz2", "EXT.csv", "EXT.gz", "EXT.ini", "EXT.7z", "EXT.access",
    "EXT.cfg",
    "EXT.pass",
    "EXT.log.bak", "EXT.log~", "EXT.log.swp", "EXT.log.zip", "EXT.log.tar.gz",
    "EXT.log.tar.bz2",
    "EXT.log.tar", "EXT.log.rar", "EXT.log.bz2", "EXT.log.gz", "EXT.log.7z",
    "EXT.sql.bak", "EXT.sql~", "EXT.sql.swp", "EXT.sql.zip", "EXT.sql.tar.gz",
    "EXT.sql.tar.bz2",
    "EXT.sql.tar", "EXT.sql.rar", "EXT.sql.bz2", "EXT.sql.gz", "EXT.sql.7z",
    "EXTdb.sql", "EXT_db.sql", "EXT-db.sql",
    "EXTdatabase.sql", "EXT_database.sql", "EXT-database.sql",
    "EXTdump.sql", "EXT_dump.sql", "EXT-dump.sql",
    "EXTbackup.sql", "EXT_backup.sql", "EXT-backup.sql",
    "backupEXT.sql", "backup_EXT.sql", "backup-EXT.sql",
    "EXTdb.sql.zip", "EXT_db.sql,zip", "EXT-db.sql.zip",
    "EXTdatabase.sql.zip", "EXT_database.sql.zip", "EXT-database.sql.zip",
    "EXTdump.sql.zip", "EXT_dump.sql.zip", "EXT-dump.sql.zip",
    "EXTbackup.sql.zip", "EXT_backup.sql.zip", "EXT-backup.sql.zip",
    "backupEXT.sql.zip", "backup_EXT.sql.zip", "backup-EXT.sql.zip",
    "EXTdatabase.sql.rar", "EXT_database.sql.rar", "EXT-database.sql.rar",
    "EXTdump.sql.rar", "EXT_dump.sql.rar", "EXT-dump.sql.rar",
    "EXTbackup.sql.rar", "EXT_backup.sql.rar", "EXT-backup.sql.rar",
    "backupEXT.sql.rar", "backup_EXT.sql.rar", "backup-EXT.sql.rar",
    "EXTdatabase.sql.gz", "EXT_database.sql.gz", "EXT-database.sql.gz",
    "EXTdump.sql.gz", "EXT_dump.sql.gz", "EXT-dump.sql.gz",
    "EXTbackup.sql.gz", "EXT_backup.sql.gz", "EXT-backup.sql.gz",
    "backupEXT.sql.gz", "backup_EXT.sql.gz", "backup-EXT.sql.gz",
    "EXTdatabase.sql.tar", "EXT_database.sql.tar", "EXT-database.sql.tar",
    "EXTdump.sql.tar", "EXT_dump.sql.tar", "EXT-dump.sql.tar",
    "EXTbackup.sql.tar", "EXT_backup.sql.tar", "EXT-backup.sql.tar",
    "backupEXT.sql.tar", "backup_EXT.sql.tar", "backup-EXT.sql.tar",
    "EXTdatabase.sql.tar.gz", "EXT_database.sql.tar.gz", "EXT-database.sql.tar.gz",
    "EXTdump.sql.tar.gz", "EXT_dump.sql.tar.gz", "EXT-dump.sql.tar.gz",
    "EXTbackup.sql.tar.gz", "EXT_backup.sql.tar.gz", "EXT-backup.sql.tar.gz",
    "backupEXT.sql.tar.gz", "backup_EXT.sql.tar.gz", "backup-EXT.sql.tar.gz",
}

// @Description：web dir的类结构
type WebDir struct {
    Domain              string                  // ip或者domain
    Port                string                  // 端口

    Request             HttpRequest

    Status404           bool
    NoFindPageOrigin    string
    DomainDirFiles      []string

    Url200              []string                    // 存在的url的列表
    ChResult            chan BurpResult             // 写入存在的url的channel，并发探测时使用
}

// @Description: burp的结果，type：1-->失败，错误信息；2-->成功，url信息
type BurpResult struct {
    Type            int
    Result          string
}

/**
 @Description：检查404
 @Param:
 @Return：
 */
func (w *WebDir) check404() error {
    for _, one := range url404 {
        err := w.Request.HttpDo(one)
        if err != nil {
            w.Status404 = false
            return err
        }

        if w.Request.Response.StatusCode == 404 {
            w.Status404 = true
            return nil
        }

        // 判断是不是无法识别的404
        notFindPage := strings.Replace(string(w.Request.ResponseBody), one, "", -1)
        if w.NoFindPageOrigin != "" {
            if w.NoFindPageOrigin != notFindPage {
                w.Status404 = false
                return fmt.Errorf("can't find NoFindPage")
            }
        }
        w.NoFindPageOrigin = notFindPage
    }

    w.Status404 = true
    return nil
}

/**
 @Description：基于一定的规则拼接处所有可能的domain dir file
 @Param:
 @Return：
 */
func (w *WebDir) getDomainDirFile() error {
    // 域名因子
    domainFactor := []string{
        w.Domain,
    }
    elems := strings.Split(w.Domain, ".")
    if len(elems) > 1 {
        domainFactor = append(domainFactor, elems[0])
    }

    // 备份文件因子
    for _, oneDomain := range domainFactor {
        for _, oneBack := range backFactor {
            w.DomainDirFiles = append(w.DomainDirFiles, strings.Replace(oneBack, "EXT", oneDomain, -1))
        }
    }

    return nil
}

/**
 @Description：执行
 @Param:
 @Return：
 */
func (w *WebDir) Do() error {
    fmt.Println("begin...")
    // 首先域名访问，保证可访问
    err := w.Request.HttpDo("")
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    // 检查是否有定义404
    err = w.check404()
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    if w.Status404 == false {
        fmt.Println(fmt.Sprintf("'%s' can't define 404", w.Domain))
        return fmt.Errorf("'%s' can't define 404", w.Domain)
    }

    // 拼装要检查的文件file
    err = w.getDomainDirFile()
    if err != nil {
        fmt.Println(fmt.Sprintf("ready domain dir files error: %s", err.Error()))
        return fmt.Errorf("ready domain dir files error: %s", err.Error())
    }

    // 开始执行爆破
    cnt := len(w.DomainDirFiles)
    fmt.Println(fmt.Sprintf("begin add '%d' task to queue.", cnt))
    // 将burp爆破的item加入到queue中
    for _, one := range w.DomainDirFiles {
        item := &BurpItem{
            Domain: w.Domain,
            Port: w.Port,
            Status404: w.Status404,
            File: one,
            ChResult: w.ChResult,
        }
        TaskQueueInstance().MsgQ.Push(item)
    }

    fmt.Println(fmt.Sprintf("wait for '%d' result.", cnt))
    for i := 0; i < cnt; i++ {
        ret := <- w.ChResult
        if ret.Type == 2 {
            w.Url200 = append(w.Url200, ret.Result)
        } else {
            fmt.Println(ret.Result)
        }
    }
    close(w.ChResult)
    fmt.Println("end....")

    return nil
}

// *********************************************************************
// 执行burp，并将结果通过channel传回
// @Description: 进去queue中执行的task item
type BurpItem struct {
    Domain          string
    Port            string
    Status404       bool

    File            string
    ChResult        chan BurpResult             // 写入存在的url的channel，并发探测时使用
}

/**
 @Description：执行爆破
 @Param:
 @Return：
 */
func (b *BurpItem) Done() {
    httpRequest := &HttpRequest{
        Domain: b.Domain,
        Port: b.Port,
    }
    err := httpRequest.HttpDo(b.File)
    if err != nil {
        burpResult := BurpResult{
            Type: 1,
            Result: err.Error(),
        }
        b.ChResult <- burpResult
        return
    }

    if b.Status404 == true {
        // 有404标记页面的
        if httpRequest.Response.StatusCode == 200 {
            // url存在
            burpResult := BurpResult{
                Type: 2,
                Result:  httpRequest.Response.Request.URL.Path,
            }
            b.ChResult <- burpResult
            return
        }
    } else {
        // 没有404标记的，需要自己判断
        // todo: 其他的条件
        if httpRequest.Response.StatusCode == 200 {
            burpResult := BurpResult{
                Type: 2,
                Result:  httpRequest.Response.Request.URL.Path,
            }
            b.ChResult <- burpResult
            return
        }
    }

    burpResult := BurpResult{
        Type: 1,
        Result: "no get url",
    }
    b.ChResult <- burpResult

    return
}