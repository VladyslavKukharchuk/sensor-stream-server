variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region to deploy to"
  type        = string
  default     = "europe-west1"
}

variable "service_name" {
  description = "The name of the Cloud Run service"
  type        = string
  default     = "sensor-stream-server"
}

variable "repository_name" {
  description = "The name of the Artifact Registry repository"
  type        = string
  default     = "sensor-stream-repo"
}

variable "firestore_database_id" {
  description = "The ID of the Firestore database (use '(default)' or a specific name)"
  type        = string
  default     = "(default)"
}
