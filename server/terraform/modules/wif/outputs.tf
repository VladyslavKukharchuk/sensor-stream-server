output "wif_provider_name" {
  value = google_iam_workload_identity_pool_provider.github_provider.name
}

output "deployer_service_account_email" {
  value = google_service_account.github_deployer.email
}
