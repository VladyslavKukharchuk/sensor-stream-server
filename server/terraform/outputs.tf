output "service_url" {
  value = module.cloud_run.service_url
}

output "artifact_registry_repo" {
  value = module.artifact_registry.repository_name
}

output "service_account_email" {
  value = module.cloud_run.service_account_email
}
