variable "project_id" {
  type = string
}

variable "github_repository" {
  type        = string
  description = "Format: username/repository-name"
}

variable "cloud_run_sa_id" {
  type        = string
  description = "Resource ID of the Cloud Run service account"
}
