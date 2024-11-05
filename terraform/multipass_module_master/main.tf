
resource "multipass_instance" "multipass_vm_master" {
  count  = var.instance_count
  cloudinit_file  = var.user_data
  name   = "${var.name_prefix}-k8s-master-${count.index + 1}"
  cpus   = var.cpus
  memory = var.memory
  disk   = var.disks
  image  = var.image_name
}

