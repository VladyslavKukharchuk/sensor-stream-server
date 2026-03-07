variable "project_id" { type = string }
variable "domain_name" { type = string }
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

resource "google_dns_record_set" "cloud_run_records" {
  for_each = {
    for record in google_cloud_run_domain_mapping.domain_mapping.status[0].resource_records :
    "${record.name}-${record.type}" => record
  }

  managed_zone = google_dns_managed_zone.zone.name

  name         = each.value.name == "@" ? google_dns_managed_zone.zone.dns_name : "${each.value.name}.${google_dns_managed_zone.zone.dns_name}"
  
  type         = each.value.type
  ttl          = 300
  rrdatas      = [each.value.rrdata]
}

output "name_servers" {
  value = google_dns_managed_zone.zone.name_servers
}
