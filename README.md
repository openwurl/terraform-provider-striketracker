Terraform Provider For Striketracker
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">


Bootstrap in progress, nothing to say yet

Backed by the [WurlWind](https://github.com/openwurl/wurlwind) library

# Maintainers

This terraform provider plugin is maintained by the Engineering team at [Wurl](https://www.wurl.com/).

# Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10.x+
- [Go](https://golang.org/doc/install) 1.11+ (to build the provider plugin)

# Getting Started (Plugin)

This is rough and will be updated

```
mkdir $GOPATH/src/github.com/openwurl; cd $GOPATH/src/github.com/openwurl
git clone git@github.com:openwurl/terraform-provider-striketracker.git
```

Enter the directory and build

```
cd $GOPATH/src/github.com/openwurl/terraform-provider-striketracker.git
make build
```

Install

```
mkdir -p ~/.terraform.d/plugins/<OS>_<ARCH>
mv terraform-provider-striketracker ~/.terraform.d/plugins/<OS>_<ARCH>/
```

# Getting Started (Usage)
`terraform init`
```
provider "striketracker" {
    authorization_header_key = ""
    application_id = ""
}
```



# Resources
Resources define infrastructure at the Striketracker/Highwinds CDN.

Many resources are interdependent, such as Hosts depending on Origins.

---
## Resource `striketracker_origin`
[Definition](resource_origin.go)

Ex.
```
resource "striketracker_origin" "test-origin" {
    account_hash = "${var.account_hash}"
    name = "Yet Another Terraform Test Origin"
    hostname = "not.real.com"
    port = 8080
    path = "/to/thing"
}
```

##### Variables
Variable descriptions can be found in 

* `name`
  * Required
  * String

* `hostname`
  * Required
  * String


* `port`
  * Required
  * Int
  * One of [80, 443, 8080, 1935]


* `account_hash`
  * Required
  * String


* `authentication_type`
  * String
  * One of [NONE, BASIC]


* `certificate_cn`
  * String


* `error_cache_ttl_seconds`
  * Int


* `max_connections_per_edge`
  * Int


* `max_connections_per_edge_enabled`
  * Int


* `maximum_origin_pull_seconds`
  * Int


* `max_retry_count`
  * Int


* `origin_cache_headers`
  * String


* `origin_default_keep_alive`
  * Int


* `origin_pull_headers`
  * String


* `origin_pull_neg_linger`
  * String


* `path`
  * String


* `request_timeout_seconds`
  * Int


* `secure_port`
  * Int


* `verify_certificate`
  * Bool



---
## Resource `striketracker_certificate`
[Definition](resource_certificate.go)

Ex.
```
resource "striketracker_certificate" "cert-name" {
    account_hash = "${var.account_hash}"
    certificate = "${data.local_file.certificate_source.content}"
    key = "${data.local_file.privkey_source.content}"
    ca_bundle = "${data.local_file.bundle_source.content}"
}
```

##### Variables
* `account_hash`
  * Required
  * String
  * The account hash within highwinds/striketracker the certificate will be deployed

* `certificate`
  * Required
  * String
  * Text of the TLS certificate
  
* `key`
  * Required
  * String
  * The text of the private key for the certificate


* `ca_bundle`
  * Technically optional but required for trust chain completion
  * String
  * The text of the FullChain / CABundle for the certificate

##### Available Outputs
* `id`
* `ciphers`
* `common_name`
* `created_date`
* `expiration_date`
* `fingerprint`
* `issuer`
* `requester`
* `trusted`
* `updated_date`


# notes for future readme

### Secret Management Ideas
https://www.tweag.io/posts/2019-04-03-terraform-provider-secret.html

### Import
Importing for striketracker works slightly different

`terraform import striketracker_certificate.cert-resource-name account_hash/certificate_id`

The account hash must be provided in format `account_hash/certificate_id` to find the existing certificate.