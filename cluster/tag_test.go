package cluster

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestTag_FindByKey(t *testing.T) {
	tags := Tags{
		&ec2.Tag{Key: aws.String("kubernetes.io/cluster/hello.k8s.local"), Value: aws.String("owned")},
		&ec2.Tag{Key: aws.String("Name"), Value: aws.String("foo")},
	}
	actual := tags.FindByKey("kubernetes.io/cluster/hello.k8s.local")
	if "owned" != actual {
		t.Errorf("FindByKey wants owned but %s", actual)
	}
}

func TestTag_FindName_NotFound(t *testing.T) {
	tags := Tags{
		&ec2.Tag{Key: aws.String("kubernetes.io/cluster/hello.k8s.local"), Value: aws.String("owned")},
		&ec2.Tag{Key: aws.String("Name"), Value: aws.String("foo")},
	}
	actual := tags.FindByKey("bar")
	if "" != actual {
		t.Errorf("FindByKey wants empty string but %s", actual)
	}
}

func TestTag_FindName(t *testing.T) {
	tags := Tags{
		&ec2.Tag{Key: aws.String("kubernetes.io/cluster/hello.k8s.local"), Value: aws.String("owned")},
		&ec2.Tag{Key: aws.String("Name"), Value: aws.String("foo")},
	}
	actual := tags.FindName()
	if "foo" != actual {
		t.Errorf("FindByKey wants foo but %s", actual)
	}
}
