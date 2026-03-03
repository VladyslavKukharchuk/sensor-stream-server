output "service_url" {
  value = google_cloud_run_v2_service.default.uri
}

output "service_account_email" {
  value = google_service_account.cloud_run_sa.email
}

output "service_account_id" {
  value = google_service_account.cloud_run_sa.name
}
