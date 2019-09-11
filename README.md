# terraform-provider-striketracker

Bootstrap in progress, nothing to say yet

# Secret Management Ideas
https://www.tweag.io/posts/2019-04-03-terraform-provider-secret.html

# notes for future readme

### Import
Importing for striketracker works slightly different

`terraform import striketracker_certificate.cert-resource-name account_hash/certificate_id`

The account hash must be provided in format `account_hash/certificate_id` to find the existing certificate.