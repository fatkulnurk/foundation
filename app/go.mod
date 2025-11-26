module github.com/fatkulnurk/foundation/app

go 1.25

require (
	github.com/fatkulnurk/foundation/shared v0.0.0-00010101000000-000000000000
	github.com/fatkulnurk/foundation/support v0.0.0-00010101000000-000000000000
)

replace (
	github.com/fatkulnurk/foundation/shared => ../shared
	github.com/fatkulnurk/foundation/support => ../support
)
