output "service_url" {
  value = module.cloud_run.service_url
}

output "artifact_registry_repo" {
  value = module.artifact_registry.repository_name
}

output "service_account_email" {
  value = module.cloud_run.service_account_email
}

output "wif_provider_name" {
  value = module.wif.wif_provider_name
}

output "wif_service_account_email" {
  value = module.wif.deployer_service_account_email
}

output "firebase_api_key" {
  value = module.firebase.api_key
}

output "name_servers" {
  value = var.domain_name != "" ? module.dns[0].name_servers : []
}
