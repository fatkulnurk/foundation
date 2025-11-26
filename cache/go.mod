module github.com/fatkulnurk/foundation/cache

go 1.25

require (
	github.com/fatkulnurk/foundation/support v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.17.0
)

replace github.com/fatkulnurk/foundation/support => ../support

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
