package rekognition_test
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
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"

	// TIP: You will often need to import the package that this test file lives
    // in. Since it is in the "test" context, it must import the package to use
    // any normal context constants, variables, or functions.
	tfrekognition "github.com/hashicorp/terraform-provider-aws/internal/service/rekognition"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this resource's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations


// TIP: ==== UNIT TESTS ====
// This is an example of a unit test. Its name is not prefixed with
// "TestAcc" like an acceptance test.
//
// Unlike acceptance tests, unit tests do not access AWS and are focused on a
// function (or method). Because of this, they are quick and cheap to run.
//
// In designing a resource's implementation, isolate complex bits from AWS bits
// so that they can be tested through a unit test. We encourage more unit tests
// in the provider.
//
// Cut and dry functions using well-used patterns, like typical flatteners and
// expanders, don't need unit testing. However, if they are complex or
// intricate, they should be unit tested.
// func TestStreamProcessorExampleUnitTest(t *testing.T) {
// 	testCases := []struct {
// 		TestName string
// 		Input    string
// 		Expected string
// 		Error    bool
// 	}{
// 		{
// 			TestName: "empty",
// 			Input:    "",
// 			Expected: "",
// 			Error:    true,
// 		},
// 		{
// 			TestName: "descriptive name",
// 			Input:    "some input",
// 			Expected: "some output",
// 			Error:    false,
// 		},
// 		{
// 			TestName: "another descriptive name",
// 			Input:    "more input",
// 			Expected: "more output",
// 			Error:    false,
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		t.Run(testCase.TestName, func(t *testing.T) {
// 			got, err := tfrekognition.FunctionFromResource(testCase.Input)

// 			if err != nil && !testCase.Error {
// 				t.Errorf("got error (%s), expected no error", err)
// 			}

// 			if err == nil && testCase.Error {
// 				t.Errorf("got (%s) and no error, expected error", got)
// 			}

// 			if got != testCase.Expected {
// 				t.Errorf("got %s, expected %s", got, testCase.Expected)
// 			}
// 		})
// 	}
// }


// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the resource as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// resource name.
//
// Acceptance test access AWS and cost money to run.
func TestAccRekognitionStreamProcessor_basicConnectedHome(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, rekognition.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_basicConnectedHome(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "input.#", "1"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "input.0.kinesis_video_stream.0.arn", "kinesisvideo", regexp.MustCompile(`stream/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "notification_channel.#", "1"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "notification_channel.0.sns_topic_arn", "sns", regexp.MustCompile(`.+`)),
					resource.TestCheckResourceAttr(resourceName, "settings.0.connected_home.0.labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.connected_home.0.labels.0", "ALL"),
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

func TestAccRekognitionStreamProcessor_connectedHome_tags(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_connectedHome_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccStreamProcessorConfig_connectedHome_tags2(rName, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccStreamProcessorConfig_connectedHome_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccRekognitionStreamProcessor_connectedHome_extendMinConfidence(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendMinConfidence(rName, "80"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "settings.0.connected_home.0.min_confidence", "80"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStreamProcessorConfig_extendMinConfidence(rName, "90"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "settings.0.connected_home.0.min_confidence", "90"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
				),
			},
		},
	})
}
func TestAccRekognitionStreamProcessor_connectedHome_extendDataSharedPreference(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendDataSharedPreference(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "data_sharing_preference.0.opt_in", "true"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStreamProcessorConfig_extendDataSharedPreference(rName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "data_sharing_preference.0.opt_in", "false"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
				),
			},
		},
	})
}

func TestAccRekognitionStreamProcessor_connectedHome_extendKmsKey(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendKmsKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor/.+$`)),
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

func TestAccRekognitionStreamProcessor_connectedHome_extendRegionsOfInterestBoundingBox_create(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestBoundingBox(rName, "0.2930403", "0.3922065", "0.15567766", "0.284666"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.bounding_box.*", map[string]string{
						"height": "0.2930403",
						"left":   "0.3922065",
						"top":    "0.15567766",
						"width":  "0.284666",
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

func TestAccRekognitionStreamProcessor_connectedHome_extendRegionsOfInterestBoundingBox_update(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestBoundingBox(rName, "0.2930403", "0.3922065", "0.15567766", "0.284666"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestBoundingBox(rName, "0.3930404", "0.3922066", "0.15567767", "0.284667"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.0.bounding_box.0.height", "0.3930404"),
				),
			},
			// {
			// 	Config: testAccStreamProcessorConfig_basicConnectedHome(rName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
			// 		resource.TestCheckResourceAttr(resourceName, "regions_of_interest.#", "0"),
			// 	),
			// },
		},
	})
}

func TestAccRekognitionStreamProcessor_connectedHome_extendRegionsOfInterestPolygon(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestPolygon(rName, "0.2930403", "0.3922065", "0.5102923", "0.7810281", "0.1209321", "0.2903921"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.0.polygon.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.2930403",
						"y": "0.3922065",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.5102923",
						"y": "0.7810281",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.1209321",
						"y": "0.2903921",
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
func TestAccRekognitionStreamProcessor_connectedHome_extendRegionsOfInterestPolygon_update(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Rekognition)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Rekognition),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestPolygon(rName, "0.2930403", "0.3922065", "0.5102923", "0.7810281", "0.1209321", "0.2903921"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.0.polygon.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStreamProcessorConfig_extendRegionsOfInterestPolygon(rName, "0.2930404", "0.3922066", "0.5102924", "0.7810282", "0.1209322", "0.2903922"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.0.polygon.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.2930404",
						"y": "0.3922066",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.5102924",
						"y": "0.7810282",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "regions_of_interest.0.polygon.*", map[string]string{
						"x": "0.1209322",
						"y": "0.2903922",
					}),
				),
			},
			{
				Config: testAccStreamProcessorConfig_basicConnectedHome(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "regions_of_interest.#", "0"),
				),
			},
		},
	})
}


func TestAccRekognitionStreamProcessor_basic(t *testing.T) {
	ctx := acctest.Context(t)
    // TIP: This is a long-running test guard for tests that run longer than
    // 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, rekognition.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, rekognition.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "maintenance_window_start_time.0.day_of_week"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "user.*", map[string]string{
						"console_access": "false",
						"groups.#":       "0",
						"username":       "Test",
						"password":       "TestTest1234",
					}),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rekognition", regexp.MustCompile(`streamprocessor:+.`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
			},
		},
	})
}

func TestAccRekognitionStreamProcessor_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var streamprocessor rekognition.DescribeStreamProcessorResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rekognition_stream_processor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, rekognition.EndpointsID)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, rekognition.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamProcessorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamProcessorConfig_basic(rName, testAccStreamProcessorVersionNewer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamProcessorExists(ctx, resourceName, &streamprocessor),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfrekognition.ResourceStreamProcessor(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckStreamProcessorDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_rekognition_stream_processor" {
				continue
			}

			input := &rekognition.DescribeStreamProcessorInput{
				StreamProcessorId: aws.String(rs.Primary.ID),
			}
			_, err := conn.DescribeStreamProcessorWithContext(ctx, &rekognition.DescribeStreamProcessorInput{
				StreamProcessorId: aws.String(rs.Primary.ID),
			})
			if err != nil {
				if tfawserr.ErrCodeEquals(err, rekognition.ErrCodeNotFoundException) {
					return nil
				}
				return err
			}

			return create.Error(names.Rekognition, create.ErrActionCheckingDestroyed, tfrekognition.ResNameStreamProcessor, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckStreamProcessorExists(ctx context.Context, name string, streamprocessor *rekognition.DescribeStreamProcessorResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.Rekognition, create.ErrActionCheckingExistence, tfrekognition.ResNameStreamProcessor, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.Rekognition, create.ErrActionCheckingExistence, tfrekognition.ResNameStreamProcessor, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionConn()
		resp, err := conn.DescribeStreamProcessorWithContext(ctx, &rekognition.DescribeStreamProcessorInput{
			StreamProcessorId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.Rekognition, create.ErrActionCheckingExistence, tfrekognition.ResNameStreamProcessor, rs.Primary.ID, err)
		}

		*streamprocessor = *resp

		return nil
	}
}

func testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionConn()

	input := &rekognition.ListStreamProcessorsInput{}
	_, err := conn.ListStreamProcessorsWithContext(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckStreamProcessorNotRecreated(before, after *rekognition.DescribeStreamProcessorResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.StreamProcessorId), aws.StringValue(after.StreamProcessorId); before != after {
			return create.Error(names.Rekognition, create.ErrActionCheckingNotRecreated, tfrekognition.ResNameStreamProcessor, aws.StringValue(before.StreamProcessorId), errors.New("recreated"))
		}

		return nil
	}
}

func testAccStreamProcessorConfig_basic(rName, version string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "test" {
  name = %[1]q
}

resource "aws_rekognition_stream_processor" "test" {
  stream_processor_name             = %[1]q
  engine_type             = "ActiveRekognition"
  engine_version          = %[2]q
  host_instance_type      = "rekognition.t2.micro"
  security_groups         = [aws_security_group.test.id]
  authentication_strategy = "simple"
  storage_type            = "efs"

  logs {
    general = true
  }

  user {
    username = "Test"
    password = "TestTest1234"
  }
}
`, rName, version)
}


func testAccStreamProcessorConfig_baseConnectedHome(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
	name = %[1]q
  path = "/service-role/"
  assume_role_policy = jsonencode({
    "Version" = "2012-10-17",
    "Statement" = [
        {
            "Effect" = "Allow",
            "Principal" = {
                "Service" = [
                    "rekognition.amazonaws.com",
                ]
            },
            "Action" = "sts:AssumeRole"
        }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRekognitionServiceRole"
}

resource "aws_iam_role_policy" "test" {
	name = %[1]q
	role = aws_iam_role.test.id
	policy = jsonencode({
			"Version" = "2012-10-17",
			"Statement" = [
					{
							"Effect" = "Allow",
							"Action" = [
								"s3:PutObject"
							],
							"Resource" = [
								"${aws_s3_bucket.test.arn}/*"
							]
					}
			]
	})
}

resource "aws_kinesis_video_stream" "test" {
  name = %[1]q
}

resource "aws_s3_bucket" "test" {
	bucket = %[1]q
}	

resource "aws_sns_topic" "test" {
	name = format("%%s-%%s", "AmazonRekognition", %[1]q)
}

`, rName)
}

func testAccStreamProcessorConfig_basicConnectedHome(rName string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
}
`, rName))
}

func testAccStreamProcessorConfig_extendMinConfidence(rName string, minConf string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
	data_sharing_preference {
		opt_in = true
	}
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
			min_confidence = %[2]q
		}
	}
}
`, rName, minConf))
}

func testAccStreamProcessorConfig_extendDataSharedPreference(rName string, optIn string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
	data_sharing_preference {
		opt_in = %[2]q
	}
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
}
`, rName, optIn))
}

func testAccStreamProcessorConfig_extendKmsKey(rName string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_kms_key" "test" {
	description = %[1]q
}

resource "aws_rekognition_stream_processor" "test" {
  name = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			arn = aws_s3_bucket.test.arn
		}
  }
	kms_key_id = aws_kms_key.test.key_id
  role_arn = aws_iam_role.test.arn
  settings {
	connected_home
		labels = ["ALL"]
  }
}
`, rName))
}

func testAccStreamProcessorConfig_extendRegionsOfInterestBoundingBox(rName string, height string, left string, top string, width string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
	regions_of_interest {
		bounding_box {
			height = %[2]q
			left = %[3]q
			top = %[4]q
			width = %[5]q
		}
	}
}
`, rName, height, left, top, width))
}

func testAccStreamProcessorConfig_extendRegionsOfInterestPolygon(rName string, x1 float64, y1 string, x2 string, y2 string, x3 string, y3 string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
	regions_of_interest {
		polygon {
			x = %[2]q
			y = %[3]q
		}
		polygon {
			x = %[4]q
			y = %[5]q
		}
		polygon {
			x = %[6]q
			y = %[7]q
		}
	}
}
`, rName, x1, y1, x2, y2, x3, y3))
}

func testAccStreamProcessorConfig_connectedHome_tags1(rName string, tagKey1 string, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
	tags = {
		%[2]q = %[3]q
	}
}
`, rName, tagKey1, tagValue1))
}

func testAccStreamProcessorConfig_connectedHome_tags2(rName string, tagKey1 string, tagValue1 string, tagKey2 string, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccStreamProcessorConfig_baseConnectedHome(rName),
		fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
		kinesis_video_stream {
			arn = aws_kinesis_video_stream.test.arn
		}
  }	
  output {
		s3_destination {
			bucket = aws_s3_bucket.test.id
		}
  }
	notification_channel {
		sns_topic_arn = aws_sns_topic.test.arn
	}
  role_arn 				  = aws_iam_role.test.arn
  settings {
		connected_home {
			labels = ["ALL"]
		}
	}
	tags = {
		%[2]q = %[3]q
		%[4]q = %[5]q
	}
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}