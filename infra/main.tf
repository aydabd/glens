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

# Workspace-driven environment config.
# Use `terraform workspace select dev` or `terraform workspace select prod`
# to target the correct GCP project. Each workspace maps to an isolated
# GCP project (separate billing account / org folder).
locals {
  env = {
    dev = {
      project       = "glens-dev"
      log_level     = "debug"
      min_instances = 0
      max_instances = 1
      alert_channel = "dev-slack"
      alert_email   = ""
    }
    prod = {
      project       = "glens-prod"
      log_level     = "info"
      min_instances = 1
      max_instances = 10
      alert_channel = "prod-pagerduty"
      alert_email   = var.alert_email
    }
  }

  cfg = local.env[contains(keys(local.env), terraform.workspace) ? terraform.workspace : "dev"]
}

provider "google" {
  project = local.cfg.project
  region  = var.region
}

module "storage" {
  source  = "./modules/storage"
  project = local.cfg.project
  region  = var.region
}

module "firestore" {
  source  = "./modules/firestore"
  project = local.cfg.project
  region  = var.region
}

module "bigquery" {
  source  = "./modules/bigquery"
  project = local.cfg.project
  region  = var.region
}

module "secrets" {
  source  = "./modules/secrets"
  project = local.cfg.project
}

module "pubsub" {
  source  = "./modules/pubsub"
  project = local.cfg.project
}

module "observability" {
  source        = "./modules/observability"
  project       = local.cfg.project
  alert_channel = local.cfg.alert_channel
  alert_email   = local.cfg.alert_email
}

module "api" {
  source        = "./modules/api"
  project       = local.cfg.project
  region        = var.region
  log_level     = local.cfg.log_level
  min_instances = local.cfg.min_instances
  max_instances = local.cfg.max_instances
  api_image_tag = var.api_image_tag

  depends_on = [
    module.firestore,
    module.bigquery,
    module.secrets,
    module.pubsub,
  ]
}
