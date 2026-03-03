output "service_url" {
  value = google_cloud_run_v2_service.default.uri
}

output "artifact_registry_repo" {
  value = google_artifact_registry_repository.repo.name
}

output "service_account_email" {
  value = google_service_account.cloud_run_sa.email
}
