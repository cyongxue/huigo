package cache

import "fmt"

type Instance interface {

}

type CreateInstance func() Instance

var adapters = make(map[string]CreateInstance)

func NewCache(adapterName string, config string) (Instance, error) {
    instanceFunc, ok := adapters[adapterName]
    if ok == false {
        return nil, fmt.Errorf("cache: unkown adapter name: %q", adapterName)
    }

    adapter := instanceFunc()
    // todo: 其他操作

    return adapter, nil
}

func Register(name string, adapterFunc CreateInstance)  {
    if adapterFunc == nil {
        panic("cache: Register adapterFunc is nil")
    }
    if _, ok := adapters[name]; ok {
        panic(fmt.Sprintf("cache: dup Register for adapter: %s", name))
    }
    adapters[name] = adapterFunc
    return
}