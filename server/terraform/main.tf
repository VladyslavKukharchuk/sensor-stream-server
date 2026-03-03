terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

module "services" {
  source = "./modules/services"
}

module "artifact_registry" {
  source          = "./modules/artifact_registry"
  region          = var.region
  repository_name = var.repository_name

  depends_on = [module.services]
}

module "database" {
  source      = "./modules/database"
  region      = var.region
  database_id = var.firestore_database_id

  depends_on = [module.services]
}

module "cloud_run" {
  source       = "./modules/cloud_run"
  project_id   = var.project_id
  region       = var.region
  service_name = var.service_name
  database_id  = var.firestore_database_id

  depends_on = [module.services, module.database]
}
