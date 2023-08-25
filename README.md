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
	redisFsm "github.com/nacknime-official/fsm-telebot-redis-storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	redisStorage := redisFsm.NewStorage(client, redisFsm.StorageSettings{})
}

```