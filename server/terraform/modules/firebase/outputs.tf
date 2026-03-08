output "api_key" {
  value = data.google_firebase_web_app_config.default.api_key
}

output "auth_domain" {
  value = data.google_firebase_web_app_config.default.auth_domain
}
