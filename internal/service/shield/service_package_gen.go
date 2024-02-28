// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package shield

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newProactiveEngagementResource,
			Name:    "Proactive Engagement",
		},
		{
			Factory: newResourceApplicationLayerAutomaticResponse,
			Name:    "Application Layer Automatic Response",
		},
		{
			Factory: newResourceDRTAccessLogBucketAssociation,
			Name:    "DRT Access Log Bucket Association",
		},
		{
			Factory: newResourceDRTAccessRoleARNAssociation,
			Name:    "DRT Access Role ARN Association",
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceProtection,
			TypeName: "aws_shield_protection",
			Name:     "Protection",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceProtectionGroup,
			TypeName: "aws_shield_protection_group",
			Name:     "Protection Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "protection_group_arn",
			},
		},
		{
			Factory:  ResourceProtectionHealthCheckAssociation,
			TypeName: "aws_shield_protection_health_check_association",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Shield
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
