package aws

const (
	BOOTSTRAP_USERDATA = `
#!/bin/bash
adduser USERNAME
echo "USERNAME ALL=(ALL) ALL" >> /etc/sudoers
echo "USERNAME ALL=NOPASSWD: ALL" >> /etc/sudoers
su - USERNAME -c "
  mkdir ~/.ssh
  touch ~/.ssh/authorized_keys
  echo 'PUBLIC_KEY' >> ~/.ssh/authorized_keys
  chmod 700 ~/.ssh
  chmod 600 ~/.ssh/authorized_keys
"
shutdown -h now
`

	SPOT_USERDATA = `
#!/bin/sh
while ! lsblk /dev/xvdf
do
  echo "Running on Spot volume, waiting for EBS attachment"
  sleep 1
done

e2label /dev/xvda1 old/
e2label /dev/xvdf1 /
shutdown -r now
`

	CLOUDFORMATION_SECURITY_GROUP = `
{
    "Description": "Detached Box Security Group",
    "Resources": {
        "InstanceSecurityGroup": {
            "Type": "AWS::EC2::SecurityGroup",
            "Properties": {
                "GroupName": "detached-security-group",
                "GroupDescription": "Enable SSH access via port 22 and Mosh connections",
                "SecurityGroupIngress": [
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "22",
                        "ToPort": "22",
                        "CidrIp": "0.0.0.0/0"
                    },
                    {
                        "IpProtocol": "udp",
                        "FromPort": "60000",
                        "ToPort": "61000",
                        "CidrIp": "0.0.0.0/0"
                    }
                ],
                "Tags": [
                    {
                        "Key" : "source",
                        "Value" : "detached"
                    }
                ]
            }
        }
    }
}
`
)
