// the text templates for shell scripts and cloudformation.
// they are stored in go so we can ship them with the binary.
package aws

const (
	// BOOTSTRAP_USERDATA creates the user in the EBS volume with ssh access
	BOOTSTRAP_USERDATA = `#!/bin/bash
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
shutdown -h now`

	// SPOT_USERDATA is the master trick behind detached
	// It: 1) waits for the EBS volume to be attached
	//     2) swap it to make it the root volume
	//     3) reboot (next boot will pick up the new root without user data)
	SPOT_USERDATA = `#!/bin/sh
while ! lsblk /dev/xvdf
do
  echo "Running on Spot volume, waiting for EBS attachment"
  sleep 1
done

e2label /dev/xvda1 old/
e2label /dev/xvdf1 /
shutdown -r now`

	// CLOUDFORMATION_SECURITY_GROUP is a basic security group that enables ssh
	// and mosh connections. The file is created to enable customization.
	// Ex: ssh port, server port, etc
	CLOUDFORMATION_SECURITY_GROUP = `{
    "Description": "Detached Security Group - DETACHED_ID",
    "Resources": {
        "DetachedSecurityGroup": {
            "Type": "AWS::EC2::SecurityGroup",
            "Properties": {
                "GroupName": "detached-security-group-DETACHED_ID",
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
    },
    "Outputs": {
        "SecurityGroupId": {
            "Description": "Security group ID",
            "Value": { "Ref": "DetachedSecurityGroup" }
        }
    }
}`
)
