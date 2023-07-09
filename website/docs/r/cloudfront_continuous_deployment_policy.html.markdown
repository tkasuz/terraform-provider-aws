---
subcategory: "CloudFront"
layout: "aws"
page_title: "AWS: aws_cloudfront_continuous_deployment_policy"
description: |-
  Terraform resource for managing an AWS CloudFront Continuous Deployment Policy.
---
<!---
TIP: A few guiding principles for writing documentation:
1. Use simple language while avoiding jargon and figures of speech.
2. Focus on brevity and clarity to keep a reader's attention.
3. Use active voice and present tense whenever you can.
4. Document your feature as it exists now; do not mention the future or past if you can help it.
5. Use accessible and inclusive language.
--->`
# Resource: aws_cloudfront_continuous_deployment_policy

Terraform resource for managing an AWS CloudFront Continuous Deployment Policy.

## Example Usage

### Single Weight Config
```terraform
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_name = "d111111abcdef8.cloudfront.net"
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
```

### Single Header Config
```terraform
resource "aws_cloudfront_continuous_deployment_policy" "test" {
  enabled = true		
  staging_distribution_dns_name = "d111111abcdef8.cloudfront.net"
  traffic_config {
	type = SingleWeight
	single_header_config = {
		header = "aws-cf-cd-test"
		value = "test"
	} 
  }
}
```

### Basic Usage

```terraform
resource "aws_cloudfront_continuous_deployment_policy" "example" {
}
```

## Argument Reference

The following arguments are required:

* `enabled` - (Required) A Boolean that indicates whether this continuous deployment policy is enabled (in effect). When this value is `true`, this policy is enabled and in effect. When this value is `false`, this policy is not enabled and has no effect.

* `staging_distribution_dns_names` - (Required) The CloudFront domain name of the staging distribution
  * `qantity` - (Required) The number of CloudFront domain names in your staging distribution.
  * `items` - (Optional) The CloudFront domain name of the staging distribution. For example: `d111111abcdef8.cloudfront.net`.

* `traffic_config` - (Optional) The traffic configuration of your continuous deployment.
  * `type` - (Required) The type of traffic configuration.
  * `single_header_config` - (Optional) Determines which HTTP requests are sent to the staging distribution.
  * `single_weight_config` - (Optional) Contains the percentage of traffic to send to the staging distribution.

### Single Header Config
* `header` - (Required) The request header name that you want CloudFront to send to your staging distribution. The header must contain the prefix `aws-cf-cd-`.
* `value` - (Required) The request header value.

### Single Weight Config
* `weight` - (Required) The percentage of traffic to send to a staging distribution, expressed as a decimal number between 0 and .15.
* `session_stickiness_config` - (Optional) Session stickiness provides the ability to define multiple requests from a single viewer as a single session. 
  * `idle_ttl` - (Required) The amount of time after which you want sessions to cease if no requests are received. Allowed values are 300–3600 seconds (5–60 minutes). The value must be less than or equal to `maximum_ttl`.
  * `maximum_ttl` - (Required) The maximum amount of time to consider requests from the viewer as being part of the same session. Allowed values are 300–3600 seconds (5–60 minutes). The value must be less than or equal to `idle_ttl`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identifier of the continuous deployment policy.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `60m`)
* `update` - (Default `180m`)
* `delete` - (Default `90m`)

## Import

CloudFront Continuous Deployment Policy can be imported using the `id`, e.g.,

```
$ terraform import aws_cloudfront_continuous_deployment_policy.example rft-8012925589
```
