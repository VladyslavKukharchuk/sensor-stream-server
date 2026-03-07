resource "google_firestore_index" "measurements_latest" {
  database = google_firestore_database.database.name

  collection = "measurements"

  fields {
    field_path = "device_id"
    order      = "ASCENDING"
  }

  fields {
    field_path = "timestamp"
    order      = "DESCENDING"
  }

  fields {
    field_path = "__name__"
    order      = "DESCENDING"
  }
}
