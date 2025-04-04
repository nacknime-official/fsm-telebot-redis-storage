module github.com/nacknime-official/fsm-telebot-redis-storage/v3

go 1.21.0

require (
	github.com/redis/go-redis/v9 v9.7.3
	github.com/stretchr/testify v1.10.0
	// sync main repository under version policy
	// for some bracnches may as pseudo version
	github.com/vitaliy-ukiru/fsm-telebot/v2 v2.0.0-beta.2
	gopkg.in/telebot.v4 v4.0.0-beta.4 // indirect

)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/vitaliy-ukiru/telebot-filter/v2 v2.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/vitaliy-ukiru/fsm-telebot/v2 v2.0.0-beta.2 => ../../fsm-telebot
