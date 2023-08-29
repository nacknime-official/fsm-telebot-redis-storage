# FSM Redis storage
Redis storage for [vitaliy-ukiru/fsm-telebot](https://github.com/vitaliy-ukiru/fsm-telebot)

## Install
```
go get github.com/nacknime-official/fsm-telebot-redis-storages
```

## Example

```go
package main

import (
	redisfsm "github.com/nacknime-official/fsm-telebot-redis-storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	redisStorage := redisfsm.NewDefaultStorage(client)
}

```

## Versioning
The module does not have a stable API at this time. A stable version will be released when the stable version policy is defined.

Currently only the latest minor versions of fsm-telebot are supported.

No versioning policy with possibility to support older versions has been defined yet.

Any ideas on versioning policy with possibility to support old versions ? Open issue.