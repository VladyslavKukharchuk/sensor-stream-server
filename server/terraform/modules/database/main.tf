resource "google_firestore_database" "database" {
  name        = var.database_id
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
}
