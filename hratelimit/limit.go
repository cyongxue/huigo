// @Description: done rate limit
package main

import (
    "github.com/golang/time/rate"
    "time"
    "log"
)

func main()  {
    r := rate.Every(2 * time.Millisecond)
    log.Println(r)

    //time.Sleep(10 * time.Second)

    limit := rate.NewLimiter(r, 10)
    for {
        if limit.Allow() {
            log.Println("allow")
        } else {
            log.Println("not allow")
        }
        time.Sleep(1 * time.Millisecond)
    }
    return
}