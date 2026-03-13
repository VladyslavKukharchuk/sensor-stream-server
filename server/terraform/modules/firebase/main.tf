resource "google_firebase_project" "default" {
  provider = google-beta
  project  = var.project_id
}

resource "google_firebase_web_app" "default" {
  provider     = google-beta
  project      = var.project_id
  display_name = "Sensor Stream Web App"

  depends_on = [google_firebase_project.default]
}

# Service Account for Firebase Auth (token signing)
resource "google_service_account" "firebase_auth_sa" {
  account_id   = "sensor-stream-firebase-auth"
  display_name = "Firebase Auth Service Account for Sensor Stream"
}

# Allow the Firebase Auth service account to sign tokens for itself
resource "google_service_account_iam_member" "firebase_auth_sa_token_creator" {
  service_account_id = google_service_account.firebase_auth_sa.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.firebase_auth_sa.email}"
}

data "google_firebase_web_app_config" "default" {
  provider   = google-beta
  project    = var.project_id
  web_app_id = google_firebase_web_app.default.app_id
}
