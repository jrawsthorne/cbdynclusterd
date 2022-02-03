packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "download_username" {
  type    = string
  default = "couchbase"
}

variable "download_password" {
  type      = string
  sensitive = true
}

variable "build_pkg" {
  type = string
}

variable "base_url" {
  type = string
}

variable "version" {
  type = string
}

variable "arch" {
  type    = string
  default = "x86_64"
}

variable "ami_name" {
  type = string
}

variable "source_ami_filter" {
  type = string
}

variable "ssh_username" {
  type = string
}

variable "ami_owner" {
  type = string
}

variable "device_name" {
  type = string
}

source "amazon-ebs" "yum" {
  ami_name      = var.ami_name
  instance_type = var.arch == "aarch64" ? "t4g.micro" : "t2.micro"
  source_ami_filter {
    filters = {
      name                = var.source_ami_filter
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = [var.ami_owner]
  }
  ssh_username          = var.ssh_username
  force_deregister      = true
  force_delete_snapshot = true
  launch_block_device_mappings {
    device_name           = var.device_name
    volume_size           = 40
    volume_type           = "gp3"
    delete_on_termination = true
  }
  snapshot_tags = {
    Name = var.ami_name
  }
}

build {
  sources = [
    "source.amazon-ebs.yum"
  ]
  provisioner "shell" {
    inline = [
      "curl -u ${var.download_username}:${var.download_password} -o ${var.build_pkg} -s ${local.build_url}",
      "sudo yum install -y ${var.build_pkg}",
      "rm ${var.build_pkg}",
      "sudo usermod -a -G couchbase ${var.ssh_username}"
    ]
  }
}

locals {
  build_url = "${var.base_url}/${var.build_pkg}"
}
