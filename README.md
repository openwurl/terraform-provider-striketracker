# terraform-provider-striketracker

![wurlwind](static/wurlwind.png)

Bootstrap in progress, nothing to say yet

Backed by the [WurlWind](https://github.com/openwurl/wurlwind) library


# Getting Started

`go get https://github.com/openwurl/terraform-provider-striketracker`

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