packer {
  required_plugins {
    proxmox = {
      version = "1.1.8"
      source  = "github.com/hashicorp/proxmox"
    }
  }
}

source "proxmox-iso" "ubuntu2404" {
  proxmox_url              = "https://localhost:8006/api2/json"
  insecure_skip_tls_verify = true
  username                 = "root@pam"
  password                 = "vagrant"

  iso_file     = "local:iso/ubuntu-24.04.1-live-server-amd64.iso"
  iso_checksum = "sha256:e240e4b801f7bb68c20d1356b60968ad0c33a41d00d828e74ceb3364a0317be9"
  unmount_iso  = true

  node    = "pve"
  sockets = 1
  cores   = 4
  memory  = 4096

  network_adapters {
    model  = "e1000"
    bridge = "vmbr0"
  }

  disks {
    disk_size    = "10G"
    storage_pool = "local-lvm"
    type         = "scsi"
  }

  cloud_init              = true
  cloud_init_storage_pool = "local-lvm"

  additional_iso_files {
    iso_storage_pool = "local"
    cd_files = ["./cloud-init/meta-data", "./cloud-init/user-data"]
    cd_label         = "cidata"
    unmount          = true
  }

  boot_wait = "1s"
  boot_command = [
    "<spacebar><wait><spacebar><wait><spacebar><wait><spacebar><wait><spacebar><wait>",
    "e<wait>",
    "<down><down><down><end><left><left><left><left><wait5>",
    "autoinstall console=ttyS0,115200",
    "<wait>",
    "<f10>",
  ]

  ssh_host     = "127.0.0.1"
  ssh_port     = 2223
  ssh_timeout  = "20m"
  ssh_username = "ansible"
  ssh_password = var.ssh_password

  serials = ["socket"]
}

build {
  sources = ["source.proxmox-iso.ubuntu2404"]

  provisioner "shell" {
    script = "scripts/non-interactive-front-end.sh"
  }

  provisioner "shell" {
    script = "scripts/update.sh"
  }
}