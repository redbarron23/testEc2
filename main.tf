provider "aws" {
  region = "${var.region}"
}

terraform {
  backend "s3" {
    bucket = "joshj-lseg-terraform-state"
    key    = "terraform"
    region = "us-east-1"
  }
}

variable "name" {
  description = "The name that will be appended to all resources so we know who it belongs to"
}

variable "region" {
  default = "us-east-1"
}

variable "ami" {
  default = "ami-a4c7edb2"
}

variable "public_key_path" {
  description = <<DESCRIPTION
Path to the SSH public key to be used for authentication.
Ensure this keypair is added to your local SSH agent so provisioners can
connect.
Example: ~/.ssh/terraform.pub
DESCRIPTION
}

variable "private_key_path" {
  description = <<DESCRIPTION
Path to the SSH public key to be used for authentication.
Ensure this keypair is added to your local SSH agent so provisioners can
connect.
Example: ~/.ssh/terraform.pem
DESCRIPTION
}

variable "key_name" {
  description = "Desired name of AWS key pair"
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = "map"
  default     = {
    Name        = "joshj"
    Owner       = "Josh Johnston"
    Lifespan    = 5
    Client      = "LSEG"
    Description = "terratest needs an env"
  }
}

resource "aws_vpc" "default_vpc" {
  cidr_block = "10.0.0.0/16"
  tags       = "${var.tags}"
}

resource "aws_subnet" "default_subnet" {
  availability_zone       = "us-east-1c"
  cidr_block              = "10.0.1.0/24"
  vpc_id                  = "${aws_vpc.default_vpc.id}"
  map_public_ip_on_launch = true
  tags                    = "${var.tags}"
}

resource "aws_internet_gateway" "default_gateway" {
  vpc_id = "${aws_vpc.default_vpc.id}"
  tags   = "${var.tags}"
}

# Grant the VPC internet access on its main route table
resource "aws_route" "internet_access" {
  route_table_id         = "${aws_vpc.default_vpc.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.default_gateway.id}"
}

resource "aws_security_group" "default_sg" {
  name        = "${var.name}-default-sg"
  description = "Used in the terraform"
  vpc_id      = "${aws_vpc.default_vpc.id}"
  tags        = "${var.tags}"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from the VPC
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "auth" {
  key_name   = "${var.name}-${var.key_name}"
  public_key = "${file(var.public_key_path)}"
}

data "aws_iam_policy_document" "lseg_instance_policy_assume_role" {
  statement {
    sid     = "1"
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      identifiers = ["ec2.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "lseg_instance_role" {
  name               = "${var.name}-lseg-instance-role"
  assume_role_policy = "${data.aws_iam_policy_document.lseg_instance_policy_assume_role.json}"
}

resource "aws_iam_role_policy_attachment" "lseg_instance_role_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
  role       = "${aws_iam_role.lseg_instance_role.name}"
}

resource "aws_iam_instance_profile" "lseg_instance_profile" {
  name = "${var.name}-lseg-instance-profile"
  role = "${aws_iam_role.lseg_instance_role.name}"
}

resource "aws_instance" "lseg_instance" {
  ami                    = "${var.ami}"
  instance_type          = "t2.micro"
  tags                   = "${var.tags}"
  key_name               = "${aws_key_pair.auth.id}"
  iam_instance_profile   = "${aws_iam_instance_profile.lseg_instance_profile.name}"
  subnet_id              = "${aws_subnet.default_subnet.id}"
  vpc_security_group_ids = ["${aws_security_group.default_sg.id}"]

  connection {
    user        = "ec2-user"
    private_key = "${file("${var.private_key_path}")}"
  }

  root_block_device {
    volume_size = 100
  }

  provisioner "remote-exec" {
    inline = [
      "echo '${var.name}' > $HOME/name-file"
    ]
  }
}

data "aws_iam_policy_document" "lseg_instance_policy_s3" {
  statement {
    sid       = "2"
    actions   = ["s3:*"]
    effect    = "Allow"
    resources = [
      "${aws_s3_bucket.lseg_instance_bucket.arn}",
      "${aws_s3_bucket.lseg_instance_bucket.arn}/*"
    ]
    principals {
      identifiers = ["ec2.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_s3_bucket" "lseg_instance_bucket" {
  bucket = "${var.name}-lseg-instance-bucket"
}

resource "aws_s3_bucket_policy" "lseg_instance_bucket_policy" {
  bucket = "${aws_s3_bucket.lseg_instance_bucket.id}"
  policy = "${data.aws_iam_policy_document.lseg_instance_policy_s3.json}"
}

output "instance_ip" {
  value = "${aws_instance.lseg_instance.public_ip}"
}

output "instance_subnet_id" {
  value = "${aws_subnet.default_subnet.id}"
}

output "instance_sg_name" {
  value = "${aws_security_group.default_sg.name}"
}

output "instance_sg_id" {
  value = "${aws_security_group.default_sg.id}"
}

output "s3_bucket" {
  value = "${aws_s3_bucket.lseg_instance_bucket.bucket}"
}

output "key_name" {
  value = "${aws_key_pair.auth.key_name}"
}

output "region" {
  value = "${var.region}"
}
