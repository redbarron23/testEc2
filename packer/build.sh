#!/bin/bash

# aws ec2 describe-images \
# 	--owners 309956199498 --query 'Images[*].[CreationDate,Name,ImageId]' \
# 	--filters "Name=name,Values=RHEL-7.?*GA*" 
#       --region eu-west-2 --output table | sort -r


#==> amazon-ebs: 2019-05-11 03:56:39 (114 MB/s) - ‘go1.12.5.linux-amd64.tar.gz’ saved [127938445/127938445]
#==> amazon-ebs: --2019-05-11 03:56:43--  ftp://https/
#==> amazon-ebs:            => ‘.listing’
#==> amazon-ebs: Resolving https (https)... failed: Name or service not known.
#==> amazon-ebs: wget: unable to resolve host address ‘https’
#==> amazon-ebs: //raw.githubusercontent.com/golang/dep/master/install.sh: Scheme missing.
#==> amazon-ebs: mkdir: cannot create directory ‘/home/ec2-user/go/bin’: No such file or directory
#==> amazon-ebs: Terminating the source AWS instance...

SOURCE_AMI="ami-0188c0c5eddd2d032"

packer validate packer.json
packer build packer.json
