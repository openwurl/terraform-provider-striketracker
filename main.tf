
provider "striketracker" {
}

#resource "striketracker_certificate" "integrations-cert" {
#    
#}

resource "striketracker_origin" "test-origin" {
    name = "Yet Another Terraform Test Origin"
    hostname = "not.real.com"
    account_hash = "z3d5t6j7"
    port = 8080
    path = "/to/thing"
}

