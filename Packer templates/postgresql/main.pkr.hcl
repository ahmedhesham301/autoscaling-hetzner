packer {
  required_plugins {
    hcloud = {
      source  = "github.com/hetznercloud/hcloud"
      version = "~> 1"
    }
    vagrant = {
      version = "~> 1"
      source  = "github.com/hashicorp/vagrant"
    }
  }
}

variable "hcloud_token" {
  type      = string
  default   = env("HKEY")
  sensitive = true
}

variable "target_env" {
  type      = string
  default   = env("TARGET_ENV")
  sensitive = false
}

variable "ENV" {
  type      = string
  sensitive = false
}

variable "config" {
  type = object({
    app_name         = string
    app_version      = string
    node_exporter    = bool
    service_exporter = bool
  })

  default = {
    app_name         = "postgresql"
    app_version      = "18"
    node_exporter    = false
    service_exporter = true
  }
}


source "hcloud" "postgresql" {
  token           = var.hcloud_token
  image           = "debian-13"
  user_data       = file("cloud-init-default.yml")
  location        = "nbg1"
  server_type     = "cx23"
  ssh_username    = "root"
  snapshot_labels = var.config
  private_ipv4    = var.ENV != "dev" ? true : false
}

source "vagrant" "postgresql" {
  communicator = "ssh"
  source_path  = "cloud-image/debian-13"
  # source_path  = "~/.vagrant.d/boxes/cloud-image-VAGRANTSLASH-debian-13/20260819.2575.0/amd64/virtualbox"
  provider  = "virtualbox"
  add_force = true
  skip_add  = true
}

build {
  sources = var.target_env == "hetzner" ? ["source.hcloud.postgresql"] : ["source.vagrant.postgresql"]

  provisioner "file" {
    source      = "${path.root}/patroni-config.yml"
    destination = "/tmp/patroni-config.yml"
  }
  provisioner "shell" {
    script          = "${path.root}/install.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }} ${var.config["app_version"]} ${var.config["node_exporter"]}"
  }
  provisioner "shell" {
    script          = "${path.root}/cleanup.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }}"
  }
}
