package main

import (
    "huigo/hwebdir"
)

func main()  {
    web := hwebdir.WebDir{
        Domain: "lightless.me",
        Port: "443",
        ChResult: make(chan hwebdir.BurpResult),
    }
    request := hwebdir.HttpRequest{
        Domain: web.Domain,
        Port: web.Port,
    }
    web.Request = request

    web.Do()
}