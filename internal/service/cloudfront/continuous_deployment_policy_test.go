package cloudfront_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/names"

	tfcloudfront "github.com/hashicorp/terraform-provider-aws/internal/service/cloudfront"
)

func TestAccCloudFrontContinuousDeploymentPolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var continuousdeploymentpolicy cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &continuousdeploymentpolicy),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "staging_distribution_dns_names.0.quantity", "true"),
					resource.TestCheckResourceAttr(resourceName, "staging_distribution_dns_names.0.items.0", "d111111abcdef8.cloudfront.net"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "staging_distribution_dns_names.0.items", map[string]string{
						"console_access": "false",
						"groups.#":       "0",
						"username":       "Test",
						"password":       "TestTest1234",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_stagingDistributionDnsNamesUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2 cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContinuousDeploymentPolicyConfig_stagingDistributionDnsNamesUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v2),
					testAccCheckContinuousDeploymentPolicyNotRecreated(&v1, &v2),
					resource.TestCheckResourceAttr(resourceName, "staging_distribution_dns_names.0.items.1", "d222222abcdef8.cloudfront.net"),
				),
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_traficConfigSingleWeightConfig(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var continuousdeploymentpolicy cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleWeightConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &continuousdeploymentpolicy),
					resource.TestCheckResourceAttr(resourceName, "traffic_config.0.single_weight_config.0.weight", "0.15"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "traffic_config.0.single_weight_config.0.session_stick_config.0.*", map[string]string{
						"idle_ttl":    "300",
						"maximum_ttl": "600",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_traficConfigSingleWeightConfigUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2 cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleWeightConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleWeightConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v2),
					testAccCheckContinuousDeploymentPolicyNotRecreated(&v1, &v2),
					resource.TestCheckResourceAttr(resourceName, "traffic_config.0.single_weight_config.0.weight", "0.1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "traffic_config.0.single_weight_config.0.session_stick_config.0.*", map[string]string{
						"idle_ttl":    "100",
						"maximum_ttl": "200",
					}),
				),
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_traficConfigSingleHeaderConfig(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var continuousdeploymentpolicy cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleHeaderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &continuousdeploymentpolicy),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "traffic_config.0.single_header_config.0.*", map[string]string{
						"header": "aws-cf-cd-test",
						"value":  "test",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_traficConfigSingleHeaderConfigUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2 cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleHeaderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContinuousDeploymentPolicyConfig_traficConfigSingleHeaderConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &v2),
					testAccCheckContinuousDeploymentPolicyNotRecreated(&v1, &v2),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "traffic_config.0.single_header_config.0.*", map[string]string{
						"header": "aws-cf-cd-test2",
						"value":  "test2",
					}),
				),
			},
		},
	})
}

func TestAccCloudFrontContinuousDeploymentPolicy_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var continuousdeploymentpolicy cloudfront.GetContinuousDeploymentPolicyOutput
	resourceName := "aws_cloudfront_continuous_deployment_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContinuousDeploymentPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccContinuousDeploymentPolicyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContinuousDeploymentPolicyExists(ctx, resourceName, &continuousdeploymentpolicy),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudfront.ResourceContinuousDeploymentPolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckContinuousDeploymentPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_cloudfront_continuous_deployment_policy" {
				continue
			}

			input := &cloudfront.GetContinuousDeploymentPolicyInput{
				Id: aws.String(rs.Primary.ID),
			}
			_, err := conn.GetContinuousDeploymentPolicyWithContext(ctx, input)
			if tfawserr.ErrCodeEquals(err, cloudfront.ErrCodeNoSuchContinuousDeploymentPolicy) {
				return nil
			}
			if err != nil {
				return nil
			}

			return create.Error(names.CloudFront, create.ErrActionCheckingDestroyed, tfcloudfront.ResNameContinuousDeploymentPolicy, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckContinuousDeploymentPolicyExists(ctx context.Context, name string, continuousdeploymentpolicy *cloudfront.GetContinuousDeploymentPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.CloudFront, create.ErrActionCheckingExistence, tfcloudfront.ResNameContinuousDeploymentPolicy, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.CloudFront, create.ErrActionCheckingExistence, tfcloudfront.ResNameContinuousDeploymentPolicy, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn(ctx)
		resp, err := conn.GetContinuousDeploymentPolicyWithContext(ctx, &cloudfront.GetContinuousDeploymentPolicyInput{
			Id: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.CloudFront, create.ErrActionCheckingExistence, tfcloudfront.ResNameContinuousDeploymentPolicy, rs.Primary.ID, err)
		}

		*continuousdeploymentpolicy = *resp

		return nil
	}
}

func testAccCheckContinuousDeploymentPolicyNotRecreated(before, after *cloudfront.GetContinuousDeploymentPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.ContinuousDeploymentPolicy.Id), aws.StringValue(after.ContinuousDeploymentPolicy.Id); before != after {
			return create.Error(names.CloudFront, create.ErrActionCheckingNotRecreated, tfcloudfront.ResNameContinuousDeploymentPolicy, before, errors.New("recreated"))
		}
		return nil
	}
}

func testAccContinuousDeploymentPolicyConfig_basic() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 1
	items = [
		"d111111abcdef8.cloudfront.net"
	]
  }
}
`
}

func testAccContinuousDeploymentPolicyConfig_stagingDistributionDnsNamesUpdate() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 2
	items = [
		"d111111abcdef8.cloudfront.net",
		"d222222abcdef8.cloudfront.net"
	]
  }
}
`
}

func testAccContinuousDeploymentPolicyConfig_traficConfigSingleWeightConfig() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 1
	items = [
		"d111111abcdef8.cloudfront.net"
	]
  }
  traffic_config {
	type = SingleWeight
	single_weight_config = {
		weight = 0.15
		session_stickiness_config = {
			idle_ttl = 300
			maximum_ttl = 600
		}
	} 
  }
}
`
}

func testAccContinuousDeploymentPolicyConfig_traficConfigSingleWeightConfigUpdate() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 1
	items = [
		"d111111abcdef8.cloudfront.net"
	]
  }
  traffic_config {
	type = SingleWeight
	single_weight_config = {
		weight = 0.1
		session_stickiness_config = {
			idle_ttl = 100
			maximum_ttl = 200
		}
	} 
  }
}
`
}

func testAccContinuousDeploymentPolicyConfig_traficConfigSingleHeaderConfig() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 1
	items = [
		"d111111abcdef8.cloudfront.net"
	]
  }
  traffic_config {
	type = SingleWeight
	single_header_config = {
		header = "aws-cf-cd-test"
		value = "test"
	} 
  }
}
`
}

func testAccContinuousDeploymentPolicyConfig_traficConfigSingleHeaderConfigUpdate() string {
	return `
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_names = {
	quantity = 1
	items = [
		"d111111abcdef8.cloudfront.net"
	]
  }
  traffic_config {
	type = SingleWeight
	single_header_config = {
		header = "aws-cf-cd-test2"
		value = "test2"
	} 
  }
}
`
}
