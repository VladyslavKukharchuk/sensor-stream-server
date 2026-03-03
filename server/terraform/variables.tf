variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region to deploy to"
  type        = string
}

variable "service_name" {
  description = "The name of the Cloud Run service"
  type        = string
}

variable "repository_name" {
  description = "The name of the Artifact Registry repository"
  type        = string
}

variable "firestore_database_id" {
  description = "The ID of the Firestore database"
  type        = string
}

variable "github_repository" {
  description = "The GitHub repository in format 'username/repo-name'"
  type        = string
}
