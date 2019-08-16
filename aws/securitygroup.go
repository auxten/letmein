package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

type SecurityGroup struct {
	Region string              `yaml:"Region"`
	SgName string              `yaml:"SgName"`
	SgId   string              `yaml:"SgId"`
	sess   *awsSession.Session `yaml:"-"`
}

func (sg *SecurityGroup) Init() (err error) {
	//select Region to use.
	conf := aws.Config{Region: aws.String(sg.Region)}
	sg.sess, err = awsSession.NewSession(&conf)
	if err != nil {
		logrus.Errorf("init aws session failed: %s", sg.Region)
		return
	}

	return
}

func (sg *SecurityGroup) ListSg(groupIds []string) (sgs []*ec2.SecurityGroup, err error) {
	// Create an EC2 service client.
	svc := ec2.New(sg.sess)

	// Retrieve the security group descriptions
	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: aws.StringSlice(groupIds),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fallthrough
			case "InvalidGroup.NotFound":
				err = fmt.Errorf("group not found %s", aerr.Message())
				return
			}
		}
		err = fmt.Errorf("unable to get descriptions for security groups, %v", err)
		return
	}

	logrus.Debug("Security Group:")
	for _, group := range result.SecurityGroups {
		logrus.Debugln(group.String())
	}
	sgs = result.SecurityGroups
	return
}

func (sg *SecurityGroup) AuthSgIngress(ip string, desc string) (err error) {
	// Create an EC2 service client.
	svc := ec2.New(sg.sess)

	cidr := fmt.Sprintf("%s/32", ip)
	req := &ec2.AuthorizeSecurityGroupIngressInput{}
	req.SetDryRun(false)
	req.SetGroupId(sg.SgId)
	req.SetGroupName(sg.SgName)
	ipranges := make([]*ec2.IpRange, 1)
	ipranges[0] = &ec2.IpRange{
		CidrIp:      &cidr,
		Description: &desc,
	}
	req.IpPermissions = make([]*ec2.IpPermission, 1)
	ipperm := &ec2.IpPermission{}
	ipperm.SetFromPort(-1)
	ipperm.SetIpProtocol("-1")
	ipperm.SetIpRanges(ipranges)
	ipperm.SetToPort(-1)

	req.IpPermissions[0] = ipperm

	output, err := svc.AuthorizeSecurityGroupIngress(req)
	logrus.Infof("Auth Security Group for %s:\n%s", cidr, output.String())

	return
}

func (sg *SecurityGroup) RevokeSgIngress(ip string) (err error) {
	// Create an EC2 service client.
	svc := ec2.New(sg.sess)

	cidr := fmt.Sprintf("%s/32", ip)
	req := &ec2.RevokeSecurityGroupIngressInput{}
	req.SetCidrIp(cidr)
	req.SetDryRun(false)
	req.SetFromPort(-1)
	req.SetGroupId(sg.SgId)
	req.SetGroupName(sg.SgName)
	req.SetIpProtocol("-1")
	req.SetToPort(-1)

	output, err := svc.RevokeSecurityGroupIngress(req)
	logrus.Infof("Revoke Security Group for %s:\n%s", cidr, output.String())

	return
}
