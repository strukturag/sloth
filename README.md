## Sloth

#### A tiny REST framework for Go

Sloth is a micro-framework for building RESTful APIs. It was forked from [Sleepy](https://github.com/dougblack/sleepy) to support the [Gorilla web toolkit](http://www.gorillatoolkit.org/pkg/mux).

```go
package main

import (
    "net/http"
    "github.com/strukturag/sloth"
)

type Item struct {}

func (item Item) Get(request *http.Request) (int, interface{}, http.Header) {
    items := []string{"item1", "item2"}
    data := map[string][]string{"items": items}
    return 200, data, http.Header{"Content-type": {"application/json"}}
}

func main() {
    item := new(Item)

    api := sloth.NewAPI()
    api.AddResource(item, "/items")
    api.Start(3000)
}
```

Now if we curl that endpoint:

```bash
$ curl localhost:3000/items
{"items": ["item1", "item2"]}
```

## Docs

Documentation lives [here](http://godoc.org/github.com/strukturag/sloth).

## License

`Sloth` is released under the [MIT License](http://opensource.org/licenses/MIT).
