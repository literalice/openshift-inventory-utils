package node

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// List nodes with the specific type from AWS Tag
func List(clusterName string, role string, roleTag string) (addresses []string, err error) {
	err = eachAwsNodeInstance(clusterName, role, roleTag, func(i ec2.Instance) {
		if *i.State.Name == "running" {
			addresses = append(addresses, *i.PrivateDnsName)
		}
	})
	return
}

func eachAwsNodeInstance(clusterName string, role string, roleTag string, f func(ec2.Instance)) error {
	svc := ec2.New(session.New())
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:kubernetes.io/cluster/" + clusterName),
				Values: []*string{
					aws.String("*"),
				},
			},
			&ec2.Filter{
				Name: aws.String("tag:" + roleTag),
				Values: []*string{
					aws.String("*" + role + "*"),
				},
			},
		},
	}
	res, err := svc.DescribeInstances(input)

	if err != nil {
		return err
	}

	for _, r := range res.Reservations {
		for _, i := range r.Instances {
			f(*i)
		}
	}

	return nil
}
