module "multipass_vm_master" {
  source = "./multipass_module_master"

  instance_count = 1
  user_data      = "${path.module}/user_data.cfg"
  name_prefix    = "arvan-challenge"
  name           = var.name
  image_name     = var.image_name
  cpus           = var.cpus
  memory         = var.memory
  disks          = var.disks
}

module "multipass_vm_worker" {
  source = "./multipass_module_worker"

  instance_count = 2
  user_data      = "${path.module}/user_data.cfg"
  name_prefix    = "arvan-challenge"
  name           = var.name
  image_name     = var.image_name
  cpus           = var.cpus
  memory         = var.memory
  disks          = var.disks
}
