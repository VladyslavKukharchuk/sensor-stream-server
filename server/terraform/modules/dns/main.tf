variable "project_id" { type = string }
variable "domain_name" { type = string } # Наприклад: "iot.nayebdebilov.com"
variable "service_name" { type = string }
variable "region" { type = string }

locals {
  base_domain = join(".", slice(split(".", var.domain_name), 1, length(split(".", var.domain_name))))
  zone_name   = replace(local.base_domain, ".", "-")
}

resource "google_dns_managed_zone" "zone" {
  name        = "${local.zone_name}-zone"
  dns_name    = "${local.base_domain}."
  description = "Managed zone for ${local.base_domain}"
}

resource "google_cloud_run_domain_mapping" "domain_mapping" {
  location = var.region
  name     = var.domain_name

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = var.service_name
  }
}

output "name_servers" {
  value = google_dns_managed_zone.zone.name_servers
}

output "dns_records" {
  value = google_cloud_run_domain_mapping.domain_mapping.status[0].resource_records
}
