packer {
  required_plugins {
    hcloud = {
      source  = "github.com/hetznercloud/hcloud"
      version = "~> 1"
    }
    # vagrant = {
    #   version = "~> 1"
    #   source  = "github.com/hashicorp/vagrant"
    # }
  }
}

variable "hcloud_token" {
  type      = string
  default   = env("HKEY")
  sensitive = true
}

variable "build_target" {
  type      = string
  default   = env("BUILD_TARGET")
  sensitive = false
}

variable "env" {
  type      = string
  sensitive = false
}

variable "networkID" {
  type      = number
  sensitive = false
}

variable "config" {
  type = object({
    type             = string
    engine           = string
    version          = string
    node_exporter    = bool
    service_exporter = bool
  })

  default = {
    type             = "database"
    engine           = "postgresql"
    version          = "18"
    node_exporter    = false
    service_exporter = true
  }
}


source "hcloud" "postgresql" {
  token           = var.hcloud_token
  image           = "debian-13"
  user_data       = file("${path.root}/../../common/cloud-init/cloud-init-default.yml")
  location        = "nbg1"
  server_type     = "cx23"
  ssh_username    = "root"
  snapshot_labels = var.config
  networks        = var.networkID != null ? [var.networkID] : []
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
  sources = var.build_target == "hetzner" ? ["source.hcloud.postgresql"] : ["source.vagrant.postgresql"]

  provisioner "shell" {
    script          = "${path.root}/../../common/scripts/update-upgrade.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }}"
  }

  provisioner "shell" {
    script          = "${path.root}/../../common/scripts/disable-service-autostart.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }}"
  }

  provisioner "shell" {
    script          = "${path.root}/../../common/scripts/install-node-exporter.sh"
    execute_command = var.config["node_exporter"] ? "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }}" : "true"
  }

  provisioner "file" {
    source      = "${path.root}/patroni-config.yml"
    destination = "/tmp/patroni-config.yml"
  }

  provisioner "shell" {
    script          = "${path.root}/install.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }} ${var.config["version"]}"
  }

  provisioner "shell" {
    script          = "${path.root}/../../common/scripts/cleanup.sh"
    execute_command = "chmod +x {{ .Path }}; sudo env {{ .Vars }} {{ .Path }}"
  }
}
