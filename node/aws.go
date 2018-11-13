package node

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (n *Node) importAwsTags(tags []*ec2.Tag) {
	for _, tag := range tags {
		if strings.HasPrefix(*tag.Key, "openshift_") {
			n.Vars[*tag.Key] = *tag.Value
		}
	}
}

// List nodes with the specific type from AWS Tag
func List(clusterName string, role string, roleTag string) (nodes []*Node, err error) {
	err = eachAwsNodeInstance(clusterName, role, roleTag, func(i *ec2.Instance, tags []*ec2.Tag) {
		if *i.State.Name == "running" {
			node := &Node{
				Host: *i.PrivateDnsName,
				Vars: make(map[string]interface{}),
			}
			node.importAwsTags(tags)
			nodes = append(nodes, node)
		}
	})
	return
}

func eachAwsNodeInstance(clusterName string, role string, roleTag string, f func(*ec2.Instance, []*ec2.Tag)) error {
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
			f(i, i.Tags)
		}
	}

	return nil
}
