package hwebdir

import "testing"

func TestWebDir_Do(t *testing.T) {
    web := WebDir{
        Domain: "lightless.me",
        Port: "443",
        ChResult: make(chan BurpResult),
    }
    request := HttpRequest{
        Domain: web.Domain,
        Port: web.Port,
    }
    web.Request = request

    web.Do()
}
