package cloudfront

import (
	"context"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_cloudfront_continuous_deployment_policy", name="Continuous Deployment Policy")
func ResourceContinuousDeploymentPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceContinuousDeploymentPolicyCreate,
		ReadWithoutTimeout:   resourceContinuousDeploymentPolicyRead,
		UpdateWithoutTimeout: resourceContinuousDeploymentPolicyUpdate,
		DeleteWithoutTimeout: resourceContinuousDeploymentPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"staging_distribution_dns_names": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quantity": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"items": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-z0-9]+.cloudfront.net`), "Must CloudFront domain name"),
							},
						},
					},
				},
			},
			"traffic_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SingleWeight",
								"SingleHeader",
							}, false),
						},
						"single_header_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: map[string]*schema.Schema{
								"header": {
									Type:         schema.TypeString,
									Required:     true,
									ValidateFunc: validation.StringMatch(regexp.MustCompile(`^aws-cf-cd-.*`), "Must start with `aws-cf-cd-`."),
								},
								"value": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
						"single_weight_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: map[string]*schema.Schema{
								"weight": {
									Type:         schema.TypeFloat,
									Required:     true,
									ValidateFunc: validation.FloatBetween(0, 0.15),
								},
								"session_stickiness_config": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: map[string]*schema.Schema{
										"idle_ttl": {
											Type:         schema.TypeInt,
											Required:     true,
											ValidateFunc: validation.IntBetween(300, 3600),
										},
										"maximum_ttl": {
											Type:         schema.TypeInt,
											Required:     true,
											ValidateFunc: validation.IntBetween(300, 3600),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameContinuousDeploymentPolicy = "Continuous Deployment Policy"
)

func resourceContinuousDeploymentPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFrontConn(ctx)

	in := &cloudfront.CreateContinuousDeploymentPolicyInput{
		ContinuousDeploymentPolicyConfig: &cloudfront.ContinuousDeploymentPolicyConfig{
			Enabled: aws.Bool((d.Get("enabled").(bool))),
		},
	}

	if v, ok := d.GetOk("staging_distribution_dns_names"); ok {
		in.ContinuousDeploymentPolicyConfig.StagingDistributionDnsNames = expandStagingDistributionDnsNames(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("traffic_config"); ok {
		in.ContinuousDeploymentPolicyConfig.TrafficConfig = expandTrafficConfig(v.([]interface{})[0].(map[string]interface{}))
	}

	out, err := conn.CreateContinuousDeploymentPolicyWithContext(ctx, in)
	if err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionCreating, ResNameContinuousDeploymentPolicy, d.Get("name").(string), err)...)
	}

	if out == nil || out.ContinuousDeploymentPolicy == nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionCreating, ResNameContinuousDeploymentPolicy, d.Get("name").(string), errors.New("empty output"))...)
	}

	d.SetId(aws.StringValue(out.ContinuousDeploymentPolicy.Id))

	return append(diags, resourceContinuousDeploymentPolicyRead(ctx, d, meta)...)
}

func resourceContinuousDeploymentPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFrontConn(ctx)

	out, err := findContinuousDeploymentPolicyByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudFront ContinuousDeploymentPolicy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionReading, ResNameContinuousDeploymentPolicy, d.Id(), err)...)
	}

	d.Set("enabled", out.ContinuousDeploymentPolicyConfig.Enabled)

	if err := d.Set("staging_distribution_dns_names", flattenStagingDistributionDnsNames(out.ContinuousDeploymentPolicyConfig.StagingDistributionDnsNames)); err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionSetting, ResNameContinuousDeploymentPolicy, d.Id(), err)...)
	}
	if err := d.Set("traffic_config", flattenTrafficConfig(out.ContinuousDeploymentPolicyConfig.TrafficConfig)); err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionSetting, ResNameContinuousDeploymentPolicy, d.Id(), err)...)
	}

	return diags
}

func resourceContinuousDeploymentPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFrontConn(ctx)

	update := false

	in := &cloudfront.UpdateContinuousDeploymentPolicyInput{
		Id: aws.String(d.Id()),
	}

	if d.HasChanges("enabled") {
		in.ContinuousDeploymentPolicyConfig.Enabled = aws.Bool(d.Get("enabled").(bool))
		update = true
	}

	if d.HasChange("staging_distribution_dns_names") {
		in.ContinuousDeploymentPolicyConfig.StagingDistributionDnsNames = expandStagingDistributionDnsNames(d.Get("staging_distribution_dns_names").([]interface{})[0].(map[string]interface{}))
	}

	if d.HasChange("traffic_config") {
		in.ContinuousDeploymentPolicyConfig.TrafficConfig = expandTrafficConfig(d.Get("traffic_config").([]interface{})[0].(map[string]interface{}))
	}

	if !update {
		return diags
	}

	log.Printf("[DEBUG] Updating CloudFront ContinuousDeploymentPolicy (%s): %#v", d.Id(), in)
	_, err := conn.UpdateContinuousDeploymentPolicyWithContext(ctx, in)
	if err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionUpdating, ResNameContinuousDeploymentPolicy, d.Id(), err)...)
	}

	return append(diags, resourceContinuousDeploymentPolicyRead(ctx, d, meta)...)
}

func resourceContinuousDeploymentPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFrontConn(ctx)

	log.Printf("[INFO] Deleting CloudFront ContinuousDeploymentPolicy %s", d.Id())

	_, err := conn.DeleteContinuousDeploymentPolicyWithContext(ctx, &cloudfront.DeleteContinuousDeploymentPolicyInput{
		Id: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, cloudfront.ErrCodeNoSuchContinuousDeploymentPolicy) {
		return diags
	}
	if err != nil {
		return append(diags, create.DiagError(names.CloudFront, create.ErrActionDeleting, ResNameContinuousDeploymentPolicy, d.Id(), err)...)
	}

	return diags
}

const (
	statusChangePending = "Pending"
	statusDeleting      = "Deleting"
	statusNormal        = "Normal"
	statusUpdated       = "Updated"
)

func findContinuousDeploymentPolicyByID(ctx context.Context, conn *cloudfront.CloudFront, id string) (*cloudfront.ContinuousDeploymentPolicy, error) {
	in := &cloudfront.GetContinuousDeploymentPolicyInput{
		Id: aws.String(id),
	}
	out, err := conn.GetContinuousDeploymentPolicyWithContext(ctx, in)
	if tfawserr.ErrCodeEquals(err, cloudfront.ErrCodeNoSuchContinuousDeploymentPolicy) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}
	if err != nil {
		return nil, err
	}

	if out == nil || out.ContinuousDeploymentPolicy == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.ContinuousDeploymentPolicy, nil
}

func flattenStagingDistributionDnsNames(sddn *cloudfront.StagingDistributionDnsNames) map[string]interface{} {
	if sddn == nil {
		return nil
	}

	m := map[string]interface{}{
		"quantity": aws.Int64Value(sddn.Quantity),
	}

	if v := sddn.Items; v != nil {
		m["items"] = flex.FlattenStringList(v)
	}

	return m
}

func flattenTrafficConfig(tf *cloudfront.TrafficConfig) map[string]interface{} {
	if tf == nil {
		return nil
	}

	m := map[string]interface{}{
		"type": aws.StringValue(tf.Type),
	}

	if v := tf.SingleHeaderConfig; v != nil {
		m["single_header_config"] = flattenSingleHeaderConfig(v)
	}

	if v := tf.SingleWeightConfig; v != nil {
		m["single_weight_config"] = flattenSingleWeightConfig(v)
	}

	return m
}

func flattenSingleHeaderConfig(cdshc *cloudfront.ContinuousDeploymentSingleHeaderConfig) []interface{} {
	if cdshc == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"header": aws.StringValue(cdshc.Header),
		"value":  aws.StringValue(cdshc.Value),
	}

	return []interface{}{m}
}

func flattenSingleWeightConfig(cdswc *cloudfront.ContinuousDeploymentSingleWeightConfig) []interface{} {
	if cdswc == nil {
		return nil
	}

	m := map[string]interface{}{
		"weight": aws.Float64Value(cdswc.Weight),
	}

	if v := cdswc.SessionStickinessConfig; v != nil {
		m["session_stickiness_config"] = flattenSessionStickinessConfig(v)
	}

	return []interface{}{m}
}

func flattenSessionStickinessConfig(ssc *cloudfront.SessionStickinessConfig) []interface{} {
	if ssc == nil {
		return nil
	}

	m := map[string]interface{}{
		"idle_ttl":    aws.Int64Value(ssc.IdleTTL),
		"maximum_ttl": aws.Int64Value(ssc.MaximumTTL),
	}

	return []interface{}{m}
}

func expandStagingDistributionDnsNames(m map[string]interface{}) *cloudfront.StagingDistributionDnsNames {
	if m == nil {
		return nil
	}

	sd := &cloudfront.StagingDistributionDnsNames{
		Quantity: aws.Int64(m["quantity"].(int64)),
	}

	if v, ok := m["items"]; ok {
		sd.Items = flex.ExpandStringList(v.([]interface{}))
	}

	return sd
}

func expandTrafficConfig(m map[string]interface{}) *cloudfront.TrafficConfig {
	if m == nil {
		return nil
	}

	tf := &cloudfront.TrafficConfig{
		Type: aws.String(m["type"].(string)),
	}

	if v, ok := m["single_head_config"]; ok {
		tf.SingleHeaderConfig = expandSingleHeaderConfig(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := m["single_weight_config"]; ok {
		tf.SingleWeightConfig = expandSingleWeightConfig(v.([]interface{})[0].(map[string]interface{}))
	}

	return tf
}

func expandSingleHeaderConfig(m map[string]interface{}) *cloudfront.ContinuousDeploymentSingleHeaderConfig {
	if m == nil {
		return nil
	}
	hc := &cloudfront.ContinuousDeploymentSingleHeaderConfig{
		Header: aws.String(m["header"].(string)),
		Value:  aws.String(m["value"].(string)),
	}
	return hc
}

func expandSingleWeightConfig(m map[string]interface{}) *cloudfront.ContinuousDeploymentSingleWeightConfig {
	if m == nil {
		return nil
	}
	wc := &cloudfront.ContinuousDeploymentSingleWeightConfig{
		Weight: aws.Float64(m["weight"].(float64)),
	}

	if v, ok := m["single_weight_config"]; ok {
		wc.SessionStickinessConfig = expandSessionStickinessConfig(v.([]interface{})[0].(map[string]interface{}))
	}
	return wc
}

func expandSessionStickinessConfig(m map[string]interface{}) *cloudfront.SessionStickinessConfig {
	if m == nil {
		return nil
	}
	sc := &cloudfront.SessionStickinessConfig{
		IdleTTL:    aws.Int64(m["idle_ttl"].(int64)),
		MaximumTTL: aws.Int64(m["maximum_ttl"].(int64)),
	}
	return sc
}
