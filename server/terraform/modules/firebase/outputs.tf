output "api_key" {
  value = data.google_firebase_web_app_config.default.api_key
}

output "auth_domain" {
  value = data.google_firebase_web_app_config.default.auth_domain
}

output "service_account_email" {
  value = google_service_account.firebase_auth_sa.email
}

output "service_account_name" {
  value = google_service_account.firebase_auth_sa.name
}
