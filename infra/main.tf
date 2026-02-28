terraform {
  required_version = ">= 1.10"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  backend "gcs" {
    bucket = "glens-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project
  region  = var.region
}

module "storage" {
  source  = "./modules/storage"
  project = var.project
  region  = var.region
}

module "firestore" {
  source  = "./modules/firestore"
  project = var.project
  region  = var.region
}

module "bigquery" {
  source  = "./modules/bigquery"
  project = var.project
  region  = var.region
}

module "secrets" {
  source  = "./modules/secrets"
  project = var.project
}

module "pubsub" {
  source  = "./modules/pubsub"
  project = var.project
}

module "observability" {
  source        = "./modules/observability"
  project       = var.project
  alert_channel = var.alert_channel
  alert_email   = var.alert_email
}

module "api" {
  source        = "./modules/api"
  project       = var.project
  region        = var.region
  log_level     = var.log_level
  min_instances = var.min_instances
  max_instances = var.max_instances
  api_image_tag = var.api_image_tag

  depends_on = [
    module.firestore,
    module.bigquery,
    module.secrets,
    module.pubsub,
  ]
}
