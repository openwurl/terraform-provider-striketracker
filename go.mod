module github.com/openwurl/terraform-provider-striketracker

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

// Dev - Uncomment to do local development
// Then run go mod tidy
// replace github.com/openwurl/wurlwind v0.0.0-20190913072758-e7ad7bcb913b => ../wurlwind

require (
	github.com/hashicorp/terraform v0.12.8
	github.com/openwurl/wurlwind v0.0.0-20190913072758-e7ad7bcb913b
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

go 1.13
