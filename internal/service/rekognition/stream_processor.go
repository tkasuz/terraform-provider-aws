package rekognition

// **PLEASE DELETE THIS AND ALL TIP COMMENTS BEFORE SUBMITTING A PR FOR REVIEW!**
//
// TIP: ==== INTRODUCTION ====
// Thank you for trying the skaff tool!
//
// You have opted to include these helpful comments. They all include "TIP:"
// to help you find and remove them when you're done with them.
//
// While some aspects of this file are customized to your input, the
// scaffold tool does *not* look at the AWS API and ensure it has correct
// function, structure, and variable names. It makes guesses based on
// commonalities. You will need to make significant adjustments.
//
// In other words, as generated, this is a rough outline of the work you will
// need to do. If something doesn't make sense for your situation, get rid of
// it.
//
// Remember to register this new resource in the provider
// (internal/provider/provider.go) once you finish. Otherwise, Terraform won't
// know about it.

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	//
	// The provider linter wants your imports to be in two groups: first,
	// standard library (i.e., "fmt" or "strings"), second, everything else.
	//
	// Also, AWS Go SDK v2 may handle nested structures differently than v1,
	// using the services/rekognition/types package. If so, you'll
	// need to import types and reference the nested types, e.g., as
	// types.<Type Name>.
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
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

// TIP: ==== FILE STRUCTURE ====
// All resources should follow this basic outline. Improve this resource's
// maintainability by sticking to it.
//
// 1. Package declaration
// 2. Imports
// 3. Main resource function with schema
// 4. Create, read, update, delete functions (in that order)
// 5. Other functions (flatteners, expanders, waiters, finders, etc.)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @SDKResource("aws_rekognition_stream_processor", name="Stream Processor")
// Tagging annotations are used for "transparent tagging".
// Change the "identifierAttribute" value to the name of the attribute used in ListTags and UpdateTags calls (e.g. "arn").
// @Tags(identifierAttribute="id")

func boundingBoxSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"height": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"left": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"top": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"width": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
			},
		},
	}
}

func polygonSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 10,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"x": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"y": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
			},
		},
	}
}

func connectedHomeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"labels": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					MaxItems: 128,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"min_confidence": {
					Type:         schema.TypeFloat,
					Optional:     true,
					ValidateFunc: validation.FloatBetween(0, 100),
				},
			},
		},
	}
}

func faceSearchSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"collection_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"face_match_treashold": {
					Type:         schema.TypeFloat,
					Optional:     true,
					Default:      80,
					ValidateFunc: validation.FloatBetween(0, 100),
				},
			},
		},
	}
}

func kinesisDataStreamSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"arn": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: verify.ValidARN,
				},
			},
		},
	}
}

func s3DestinationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bucket": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(3, 255),
				},
				"key_prefix": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(0, 1024),
				},
			},
		},
	}
}

// @SDKResource("aws_rekognition_stream_processor", name="Stream Processor")
func ResourceStreamProcessor() *schema.Resource {
	return &schema.Resource{
		// TIP: ==== ASSIGN CRUD FUNCTIONS ====
		// These 4 functions handle CRUD responsibilities below.
		CreateWithoutTimeout: resourceStreamProcessorCreate,
		ReadWithoutTimeout:   resourceStreamProcessorRead,
		UpdateWithoutTimeout: resourceStreamProcessorUpdate,
		DeleteWithoutTimeout: resourceStreamProcessorDelete,

		// TIP: ==== TERRAFORM IMPORTING ====
		// If Read can get all the information it needs from the Identifier
		// (i.e., d.Id()), you can use the Passthrough importer. Otherwise,
		// you'll need a custom import function.
		//
		// See more:
		// https://hashicorp.github.io/terraform-provider-aws/add-import-support/
		// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#implicit-state-passthrough
		// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#virtual-attributes
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// TIP: ==== CONFIGURABLE TIMEOUTS ====
		// Users can configure timeout lengths but you need to use the times they
		// provide. Access the timeout they configure (or the defaults) using,
		// e.g., d.Timeout(schema.TimeoutCreate) (see below). The times here are
		// the defaults if they don't configure timeouts.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		// TIP: ==== SCHEMA ====
		// In the schema, add each of the attributes in snake case (e.g.,
		// delete_automated_backups).
		//
		// Formatting rules:
		// * Alphabetize attributes to make them easier to find.
		// * Do not add a blank line between attributes.
		//
		// Attribute basics:
		// * If a user can provide a value ("configure a value") for an
		//   attribute (e.g., instances = 5), we call the attribute an
		//   "argument."
		// * You change the way users interact with attributes using:
		//     - Required
		//     - Optional
		//     - Computed
		// * There are only four valid combinations:
		//
		// 1. Required only - the user must provide a value
		// Required: true,
		//
		// 2. Optional only - the user can configure or omit a value; do not
		//    use Default or DefaultFunc
		// Optional: true,
		//
		// 3. Computed only - the provider can provide a value but the user
		//    cannot, i.e., read-only
		// Computed: true,
		//
		// 4. Optional AND Computed - the provider or user can provide a value;
		//    use this combination if you are using Default or DefaultFunc
		// Optional: true,
		// Computed: true,
		//
		// You will typically find arguments in the input struct
		// (e.g., CreateDBInstanceInput) for the create operation. Sometimes
		// they are only in the input struct (e.g., ModifyDBInstanceInput) for
		// the modify operation.
		//
		// For more about schema options, visit
		// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"arn": { // TIP: Many, but not all, resources have an `arn` attribute.
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"data_sharing_preference": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"opt_in": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"last_update_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_channel": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sns_topic_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"output": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kinesis_data_stream": kinesisDataStreamSchema(),
						"s3_destination":      s3DestinationSchema(),
					},
				},
			},
			"regions_of_interest": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bounding_box": boundingBoxSchema(),
						"polygon":      polygonSchema(),
					},
				},
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connected_home": connectedHomeSchema(),
						"face_search":    faceSearchSchema(),
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameStreamProcessor = "Stream Processor"
)

func resourceStreamProcessorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE CREATE ====
	// Generally, the Create function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a create input structure
	// 3. Call the AWS create/put function
	// 4. Using the output from the create function, set the minimum arguments
	//    and attributes for the Read function to work. At a minimum, set the
	//    resource ID. E.g., d.SetId(<Identifier, such as AWS ID or ARN>)
	// 5. Use a waiter to wait for create to complete
	// 6. Call the Read function in the Create return

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).RekognitionClient()

	name := d.Get("name").(string)

	// TIP: -- 2. Populate a create input structure
	in := &rekognition.CreateStreamProcessorInput{
		// TIP: Mandatory or fields that will always be present can be set when
		// you create the Input structure. (Replace these with real fields.)
		Input:    extractInput(d.Get("input").([]interface{})),
		Name:     aws.String(name),
		Output:   extractOutput(d.Get("output").([]interface{})),
		RoleArn:  aws.String(d.Get("role_arn").(string)),
		Settings: &types.StreamProcessorSettings{},

		// TIP: Not all resources support tags and tags don't always make sense. If
		// your resource doesn't need tags, you can remove the tags lines here and
		// below. Many resources do include tags so this a reminder to include them
		// where possible.
		Tags: GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("data_sharing_preference"); ok {
		// TIP: Optional fields should be set based on whether or not they are
		// used.
		in.DataSharingPreference = extractDataSharingPreference(v.([]interface{}))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		// TIP: Optional fields should be set based on whether or not they are
		// used.
		in.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("notification_channel"); ok {
		in.NotificationChannel = extractNotificationChannel(v.([]interface{}))
	}

	if v, ok := d.GetOk("output"); ok {
		in.Output = extractOutput(v.([]interface{}))
	}

	if v, ok := d.GetOk("regions_of_interest"); ok {
		in.RegionsOfInterest = extractRegionsOfInterest(v.([]interface{}))
	}

	if v, ok := d.GetOk("settings"); ok {
		in.Settings = extractSettings(v.([]interface{}))
	}

	// TIP: -- 3. Call the AWS create function
	out, err := conn.CreateStreamProcessor(ctx, in)
	if err != nil {
		// TIP: Since d.SetId() has not been called yet, you cannot use d.Id()
		// in error messages at this point.
		return create.DiagError(names.Rekognition, create.ErrActionCreating, ResNameStreamProcessor, d.Get("name").(string), err)
	}

	if out == nil || out.StreamProcessorArn == nil {
		return create.DiagError(names.Rekognition, create.ErrActionCreating, ResNameStreamProcessor, d.Get("name").(string), errors.New("empty output"))
	}

	// TIP: -- 4. Set the minimum arguments and/or attributes for the Read function to
	// work.
	arn := *out.StreamProcessorArn
	d.SetId(arn[strings.LastIndex(arn, "/")+1:])

	// TIP: -- 6. Call the Read function in the Create return
	return resourceStreamProcessorRead(ctx, d, meta)
}

func resourceStreamProcessorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE READ ====
	// Generally, the Read function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Get the resource from AWS
	// 3. Set ID to empty where resource is not new and not found
	// 4. Set the arguments and attributes
	// 5. Set the tags
	// 6. Return nil

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).RekognitionClient()

	// TIP: -- 2. Get the resource from AWS using an API Get, List, or Describe-
	// type function, or, better yet, using a finder.
	out, err := findStreamProcessorByID(ctx, conn, d.Id())

	// TIP: -- 3. Set ID to empty where resource is not new and not found
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Rekognition StreamProcessor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionReading, ResNameStreamProcessor, d.Id(), err)
	}

	// TIP: -- 4. Set the arguments and attributes
	//
	// For simple data types (i.e., schema.TypeString, schema.TypeBool,
	// schema.TypeInt, and schema.TypeFloat), a simple Set call (e.g.,
	// d.Set("arn", out.Arn) is sufficient. No error or nil checking is
	// necessary.
	//
	// However, there are some situations where more handling is needed.
	// a. Complex data types (e.g., schema.TypeList, schema.TypeSet)
	// b. Where errorneous diffs occur. For example, a schema.TypeString may be
	//    a JSON. AWS may return the JSON in a slightly different order but it
	//    is equivalent to what is already set. In that case, you may check if
	//    it is equivalent before setting the different JSON.
	d.Set("arn", out.StreamProcessorArn)
	d.Set("name", out.Name)
	d.Set("creation_timestamp", aws.ToTime(out.CreationTimestamp).Format(time.RFC3339))
	d.Set("kms_key_id", out.KmsKeyId)
	d.Set("last_update_timestamp", aws.ToTime(out.LastUpdateTimestamp).Format(time.RFC3339))
	d.Set("role_arn", out.RoleArn)


	if err := d.Set(("data_sharing_preference"), flattenDataSharingPreference(out.DataSharingPreference)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}

	if err := d.Set(("input"), flattenInput(out.Input)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}

	if err := d.Set(("notification_channel"), flattenNotificationChannel(out.NotificationChannel)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}

	if err := d.Set(("output"), flattenOutput(out.Output)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}

	if err := d.Set(("regios_of_interest"), flattenRegionsOfInterest(out.RegionsOfInterest)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}

	if err := d.Set(("settings"), flattenSettings(out.Settings)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionSetting, ResNameStreamProcessor, d.Id(), err)
	}
	// TIP: Setting a complex type.
	// For more information, see:
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#data-handling-and-conversion
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#flatten-functions-for-blocks
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#root-typeset-of-resource-and-aws-list-of-structure

	// TIP: -- 6. Return nil
	return nil
}

func resourceStreamProcessorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE UPDATE ====
	// Not all resources have Update functions. There are a few reasons:
	// a. The AWS API does not support changing a resource
	// b. All arguments have ForceNew: true, set
	// c. The AWS API uses a create call to modify an existing resource
	//
	// In the cases of a. and b., the main resource function will not have a
	// UpdateWithoutTimeout defined. In the case of c., Update and Create are
	// the same.
	//
	// The rest of the time, there should be an Update function and it should
	// do the following things. Make sure there is a good reason if you don't
	// do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a modify input structure and check for changes
	// 3. Call the AWS modify/update function
	// 4. Use a waiter to wait for update to complete
	// 5. Call the Read function in the Update return

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).RekognitionClient()

	// TIP: -- 2. Populate a modify input structure and check for changes
	//
	// When creating the input structure, only include mandatory fields. Other
	// fields are set as needed. You can use a flag, such as update below, to
	// determine if a certain portion of arguments have been changed and
	// whether to call the AWS update function.
	update := false

	in := &rekognition.UpdateStreamProcessorInput{
		Name: aws.String(d.Id()),
	}

	if d.HasChanges("data_sharing_preference") {
		in.DataSharingPreferenceForUpdate = extractDataSharingPreference(d.Get("data_sharing_preference").([]interface{}))
		if in.DataSharingPreferenceForUpdate != nil {
			update = true
		}
	}

	if d.HasChanges("regions_of_interest") {
		in.RegionsOfInterestForUpdate = extractRegionsOfInterest(d.Get("regions_of_interest").([]interface{}))
		if in.RegionsOfInterestForUpdate != nil {
			in.ParametersToDelete = append(in.ParametersToDelete, "RegionsOfInterest")
			update = true
		}
	}

	if d.HasChanges("settings") {
		in.SettingsForUpdate = extractSettingsForUpdate(d.Get("settings").([]interface{}))
		if in.SettingsForUpdate.ConnectedHomeForUpdate != nil {
			in.ParametersToDelete = append(in.ParametersToDelete, "ConnectedHomeMinConfidence")
			update = true
		}
	}

	if !update {
		return nil
	}

	// TIP: -- 3. Call the AWS modify/update function
	log.Printf("[DEBUG] Updating Rekognition StreamProcessor (%s): %#v", d.Id(), in)
	_, err := conn.UpdateStreamProcessor(ctx, in)
	if err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionUpdating, ResNameStreamProcessor, d.Id(), err)
	}

	// TIP: -- 4. Use a waiter to wait for update to complete
	if _, err := waitStreamProcessorUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return create.DiagError(names.Rekognition, create.ErrActionWaitingForUpdate, ResNameStreamProcessor, d.Id(), err)
	}

	// TIP: -- 5. Call the Read function in the Update return
	return resourceStreamProcessorRead(ctx, d, meta)
}

func resourceStreamProcessorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE DELETE ====
	// Most resources have Delete functions. There are rare situations
	// where you might not need a delete:
	// a. The AWS API does not provide a way to delete the resource
	// b. The point of your resource is to perform an action (e.g., reboot a
	//    server) and deleting serves no purpose.
	//
	// The Delete function should do the following things. Make sure there
	// is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a delete input structure
	// 3. Call the AWS delete function
	// 4. Use a waiter to wait for delete to complete
	// 5. Return nil

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).RekognitionClient()

	// TIP: -- 2. Populate a delete input structure
	log.Printf("[INFO] Deleting Rekognition StreamProcessor %s", d.Id())

	// TIP: -- 3. Call the AWS delete function
	_, err := conn.DeleteStreamProcessor(ctx, &rekognition.DeleteStreamProcessorInput{
		Name: aws.String(d.Id()),
	})

	// TIP: On rare occassions, the API returns a not found error after deleting a
	// resource. If that happens, we don't want it to show up as an error.
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.Rekognition, create.ErrActionDeleting, ResNameStreamProcessor, d.Id(), err)
	}
	return nil
}

// TIP: ==== STATUS CONSTANTS ====
// Create constants for states and statuses if the service does not
// already have suitable constants. We prefer that you use the constants
// provided in the service if available (e.g., amp.WorkspaceStatusCodeActive).
const (
	statusStopped        = "STOPPED"
	statusUpdating       = "UPDATING"
)

// TIP: ==== WAITERS ====
// Some resources of some services have waiters provided by the AWS API.
// Unless they do not work properly, use them rather than defining new ones
// here.
//
// Sometimes we define the wait, status, and find functions in separate
// files, wait.go, status.go, and find.go. Follow the pattern set out in the
// service and define these where it makes the most sense.
//
// If these functions are used in the _test.go file, they will need to be
// exported (i.e., capitalized).
//
// You will need to adjust the parameters and names to fit the service.

// TIP: It is easier to determine whether a resource is updated for some
// resources than others. The best case is a status flag that tells you when
// the update has been fully realized. Other times, you can check to see if a
// key resource argument is updated to a new value or not.

func waitStreamProcessorUpdated(ctx context.Context, conn *rekognition.Client, id string, timeout time.Duration) (*rekognition.DescribeStreamProcessorOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{statusUpdating},
		Target:                    []string{statusStopped},
		Refresh:                   statusStreamProcessor(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*rekognition.DescribeStreamProcessorOutput); ok {
		return out, err
	}

	return nil, err
}

// TIP: ==== STATUS ====
// The status function can return an actual status when that field is
// available from the API (e.g., out.Status). Otherwise, you can use custom
// statuses to communicate the states of the resource.
//
// Waiters consume the values returned by status functions. Design status so
// that it can be reused by a create, update, and delete waiter, if possible.

func statusStreamProcessor(ctx context.Context, conn *rekognition.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findStreamProcessorByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.Status), nil
	}
}

// TIP: ==== FINDERS ====
// The find function is not strictly necessary. You could do the API
// request from the status function. However, we have found that find often
// comes in handy in other places besides the status function. As a result, it
// is good practice to define it separately.

func findStreamProcessorByID(ctx context.Context, conn *rekognition.Client, id string) (*rekognition.DescribeStreamProcessorOutput, error) {
	in := &rekognition.DescribeStreamProcessorInput{
		Name: aws.String(id),
	}
	out, err := conn.DescribeStreamProcessor(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}
		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

// TIP: ==== FLEX ====
// Flatteners and expanders ("flex" functions) help handle complex data
// types. Flatteners take an API data type and return something you can use in
// a d.Set() call. In other words, flatteners translate from AWS -> Terraform.
//
// On the other hand, expanders take a Terraform data structure and return
// something that you can send to the AWS API. In other words, expanders
// translate from Terraform -> AWS.
//
// See more:
// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/

// TIP: Often the AWS API will return a slice of structures in response to a
// request for information. Sometimes you will have set criteria (e.g., the ID)
// that means you'll get back a one-length slice. This plural function works
// brilliantly for that situation too.

func flattenDataSharingPreference(apiObject *types.StreamProcessorDataSharingPreference) map[string]bool {
	if apiObject == nil {
		return nil
	}
	return map[string]bool{
		"opt_in": apiObject.OptIn,
	}
}


func flattenInput(apiObject *types.StreamProcessorInput) map[string]interface{} {
	if apiObject == nil {
		return nil
	}
	m := map[string]interface{}{}
	if v := apiObject.KinesisVideoStream.Arn; v != nil {
		m["kinesis_video_stream"] = map[string]string{
			"arn": aws.ToString(v),
		}
	}
	return m
}

func flattenNotificationChannel(apiObject *types.StreamProcessorNotificationChannel) map[string]string {
	if apiObject == nil {
		return nil
	}
	return map[string]string{
		"sns_topic_arn": aws.ToString(apiObject.SNSTopicArn),
	}
}

func flattenOutput(apiObject *types.StreamProcessorOutput) map[string]interface{} {
	if apiObject == nil {
		return nil
	}
	m := map[string]interface{}{}
	if v := apiObject.KinesisDataStream.Arn; v != nil {
		m["kinesis_video_stream"] = map[string]string{
			"arn": aws.ToString(v),
		}
	}
	if v := apiObject.S3Destination; v != nil {
		s := map[string]string{}
		if bucket := v.Bucket; bucket != nil {
			s["bucket"] = aws.ToString(bucket)
		}
		if prefix := v.KeyPrefix; prefix != nil {
			s["key_prefix"] = aws.ToString(prefix)
		}
		m["s3_destination"] = s
	}
	return m
}

func flattenRegionsOfInterest(apiObjects []types.RegionOfInterest) []map[string]interface{} { 
	if len(apiObjects) == 0 {
		return nil
	}
	var l []map[string]interface{}
	for _, apiObject := range apiObjects {
		m := map[string]interface{}{}
		if v := apiObject.BoundingBox; v != nil {
			m["bounding_box"] = flattenBoundingBox(v)
		}
		if v := apiObject.Polygon; v != nil {
			m["polygon"] = flattenPoints(v)
		}
		l = append(l, m)
	}
	return l
}

func flattenBoundingBox(apiObject *types.BoundingBox) map[string]float32 {
	b := map[string]float32{}
	if height := apiObject.Height; height != nil {
		b["height"] = aws.ToFloat32(height)
	}
	if left := apiObject.Height; left != nil {
		b["left"] = aws.ToFloat32(left)
	}
	if top := apiObject.Height; top != nil {
		b["top"] = aws.ToFloat32(top)
	}
	if width := apiObject.Width; width != nil {
		b["width"] = aws.ToFloat32(width)
	}
	return b
}

func flattenPoint(apiObject types.Point) map[string]float32 {
	p := map[string]float32{}
	if x := apiObject.X; x != nil {
		p["x"] = aws.ToFloat32(x)
	}
	if y := apiObject.Y; y != nil {
		p["y"] = aws.ToFloat32(y)
	}
	return p
}

func flattenPoints(apiObjects []types.Point) []map[string]float32 {
	if len(apiObjects) == 0 {
		return nil
	}
	var l []map[string]float32
	for _, apiObject := range apiObjects {
		l = append(l, flattenPoint(apiObject))
	}
	return l
}


func flattenSettings(apiObject *types.StreamProcessorSettings) map[string]interface{} {
	if apiObject == nil {
		return nil
	}
	m := map[string]interface{}{}
	if v := apiObject.ConnectedHome; v != nil {
		c := map[string]interface{}{}
		if labs := v.Labels; len(labs) != 0 {
			c["labels"] = flex.FlattenStringValueList(labs)
		}
		if min := v.MinConfidence; min != nil {
			c["min_confidence"] = aws.ToFloat32(min)
		}
		m["connected_home"] = c
	}
	if v := apiObject.FaceSearch; v != nil {
		f := map[string]interface{}{}
		if id := v.CollectionId; v != nil {
			f["collection_id"] = aws.ToString(id)
		}
		if thr := v.FaceMatchThreshold; thr != nil {
			f["face_match_threshold"] = aws.ToFloat32(thr)
		}
		m["face_search"] = f
	}
	return m
}

// TIP: Remember, as mentioned above, expanders take a Terraform data structure
// and return something that you can send to the AWS API. In other words,
// expanders translate from Terraform -> AWS.
//
// See more:
// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/

func extractInput(tfList []interface{}) *types.StreamProcessorInput {
	if len(tfList) == 0 {
		return nil
	}
	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}
	result := &types.StreamProcessorInput{}
	if v, ok := tfMap["kinesis_video_stream"]; ok {

		result.KinesisVideoStream = &types.KinesisVideoStream{
			Arn: extractArn(v.([]interface{})),
		}
	}
	return result
}

func extractArn(tfList []interface{}) *string {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	if v, ok := tfMap["arn"].(string); ok && v != "" {
		return aws.String(v)
	} else {
		return nil
	}
}

func extractDataSharingPreference(tfList []interface{}) *types.StreamProcessorDataSharingPreference {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.StreamProcessorDataSharingPreference{}

	if v, ok := tfMap["opt_in"].(bool); ok {
		result.OptIn = v
	}
	return result
}

func extractNotificationChannel(tfList []interface{}) *types.StreamProcessorNotificationChannel {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.StreamProcessorNotificationChannel{}

	if v, ok := tfMap["sns_topic_arn"].(string); ok && v != "" {
		result.SNSTopicArn = aws.String(v)
	}
	return result
}

func extractOutput(tfList []interface{}) *types.StreamProcessorOutput {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.StreamProcessorOutput{}

	if v, ok := tfMap["kinesis_data_stream"]; ok {
		result.KinesisDataStream = &types.KinesisDataStream{
			Arn: extractArn(v.([]interface{})),
		}
	}
	if v, ok := tfMap["s3_destination"]; ok {
		result.S3Destination = extractS3Destination(v.([]interface{}))
	}
	return result
}

func extractS3Destination(tfList []interface{}) *types.S3Destination {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.S3Destination{}

	if v, ok := tfMap["bucket"].(string); ok && v != "" {
		result.Bucket = aws.String(v)
	}
	if v, ok := tfMap["key_prefix"].(string); ok && v != "" {
		result.Bucket = aws.String(v)
	}
	return result
}

func extractRegionsOfInterest(tfList []interface{}) []types.RegionOfInterest {
	if len(tfList) == 0 {
		return nil
	}

	var results []types.RegionOfInterest

	for _, r := range tfList {
		m := r.(map[string]interface{})
		result := types.RegionOfInterest{}
		if v, ok := m["bounding_box"]; ok {
			result.BoundingBox = extractBoundingBox(v.([]interface{}))
		}
		if v, ok := m["polygon"]; ok {
			result.Polygon = extractPolygon(v.([]interface{}))
		}

		results = append(results, result)
	}

	return results
}

func extractBoundingBox(tfList []interface{}) *types.BoundingBox {
	if len(tfList) == 0 {
		return nil
	}
	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}
	result := &types.BoundingBox{}
	if v, ok := tfMap["height"].(float32); ok {
		result.Height = aws.Float32(v)
	}
	if v, ok := tfMap["left"].(float32); ok {
		result.Left = aws.Float32(v)
	}
	if v, ok := tfMap["top"].(float32); ok {
		result.Top = aws.Float32(v)
	}
	if v, ok := tfMap["width"].(float32); ok {
		result.Width = aws.Float32(v)
	}
	return result
}

func extractPolygon(tfList []interface{}) []types.Point {
	if len(tfList) == 0 {
		return nil
	}

	var results []types.Point

	for _, r := range tfList {
		m := r.(map[string]interface{})
		result := types.Point{}
		if v, ok := m["x"].(float32); ok {
			result.X = aws.Float32(v)
		}
		if v, ok := m["y"].(float32); ok {
			result.Y = aws.Float32(v)
		}
		results = append(results, result)
	}

	return results
}


func extractSettings(tfList []interface{}) *types.StreamProcessorSettings {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.StreamProcessorSettings{}

	if v, ok := tfMap["connected_home"]; ok {
		result.ConnectedHome = extractConnectedHome(v.([]interface{}))
	}
	if v, ok := tfMap["face_search"]; ok {
		result.FaceSearch = extractFaceSearch(v.([]interface{}))
	}
	return result
}

func extractSettingsForUpdate(tfList []interface{}) *types.StreamProcessorSettingsForUpdate {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.StreamProcessorSettingsForUpdate{}

	if v, ok := tfMap["connected_home"]; ok {
		result.ConnectedHomeForUpdate = extractConnectedHomeForUpdate(v.([]interface{}))
	}
	return result
}



func extractConnectedHome(tfList []interface{}) *types.ConnectedHomeSettings {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.ConnectedHomeSettings{}

	if v, ok := tfMap["labels"].([]interface{}); ok {
		result.Labels = flex.ExpandStringValueList(v)
	}
	if v, ok := tfMap["min_confidence"].(float32); ok {
		result.MinConfidence = aws.Float32(v)
	}
	return result
}

func extractConnectedHomeForUpdate(tfList []interface{}) *types.ConnectedHomeSettingsForUpdate {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.ConnectedHomeSettingsForUpdate{}

	if v, ok := tfMap["labels"].([]interface{}); ok {
		result.Labels = flex.ExpandStringValueList(v)
	}
	if v, ok := tfMap["min_confidence"].(float32); ok {
		result.MinConfidence = aws.Float32(v)
	}
	return result
}

func extractFaceSearch(tfList []interface{}) *types.FaceSearchSettings {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := &types.FaceSearchSettings{}

	if v, ok := tfMap["collection_id"].(string); ok && v != "" {
		result.CollectionId = aws.String(v)
	}
	if v, ok := tfMap["face_match_threshold"].(float32); ok {
		result.FaceMatchThreshold = aws.Float32(v)
	}
	return result
}
