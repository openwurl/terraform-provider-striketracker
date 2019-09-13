module github.com/openwurl/terraform-provider-striketracker

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

// Dev
replace github.com/openwurl/wurlwind/striketracker => /Users/ccorbett/go/src/github.com/openwurl/wurlwind/striketracker

require (
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/hashicorp/terraform v0.12.8
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/openwurl/wurlwind v0.0.0-20190909185455-f9f0a76e8d0a
	github.com/wurlinc/hls-config v0.0.0-20190511002729-10a112b30d34 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1 // indirect
)

go 1.13
