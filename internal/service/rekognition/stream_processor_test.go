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
	//
	// Also, AWS Go SDK v2 may handle nested structures differently than v1,
	// using the services/rekognition/types package. If so, you'll
	// need to import types and reference the nested types, e.g., as
	// types.<Type Name>.
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
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
	"github.com/hashicorp/terraform-provider-aws/names"
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
func TestStreamProcessorExampleUnitTest(t *testing.T) {
	testCases := []struct {
		TestName string
		Input    string
		Expected string
		Error    bool
	}{
		{
			TestName: "empty",
			Input:    "",
			Expected: "",
			Error:    true,
		},
		{
			TestName: "descriptive name",
			Input:    "some input",
			Expected: "some output",
			Error:    false,
		},
		{
			TestName: "another descriptive name",
			Input:    "more input",
			Expected: "more output",
			Error:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			got, err := tfrekognition.FunctionFromResource(testCase.Input)

			if err != nil && !testCase.Error {
				t.Errorf("got error (%s), expected no error", err)
			}

			if err == nil && testCase.Error {
				t.Errorf("got (%s) and no error, expected error", got)
			}

			if got != testCase.Expected {
				t.Errorf("got %s, expected %s", got, testCase.Expected)
			}
		})
	}
}

// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the resource as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// resource name.
//
// Acceptance test access AWS and cost money to run.
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
			acctest.PreCheckPartitionHasService(t, names.RekognitionEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RekognitionEndpointID),
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

func TestAccRekognitionStreamProcessor_faceSearch(t *testing.T) {
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
			acctest.PreCheckPartitionHasService(t, names.RekognitionEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RekognitionEndpointID),
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
			acctest.PreCheckPartitionHasService(t, names.RekognitionEndpointID)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RekognitionEndpointID),
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
		conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_rekognition_stream_processor" {
				continue
			}

			input := &rekognition.DescribeStreamProcessorInput{
				StreamProcessorId: aws.String(rs.Primary.ID),
			}
			_, err := conn.DescribeStreamProcessor(ctx, &rekognition.DescribeStreamProcessorInput{
				StreamProcessorId: aws.String(rs.Primary.ID),
			})
			if err != nil {
				var nfe *types.ResourceNotFoundException
				if errors.As(err, &nfe) {
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

		conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionClient()
		resp, err := conn.DescribeStreamProcessor(ctx, &rekognition.DescribeStreamProcessorInput{
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
	conn := acctest.Provider.Meta().(*conns.AWSClient).RekognitionClient()

	input := &rekognition.ListStreamProcessorsInput{}
	_, err := conn.ListStreamProcessors(ctx, input)

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

func testAccStreamProcessorConfig_baseIAMRole() string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
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
  role       = aws_iam_role.test[count.index].name
  policy_arn = "arn:aws:iam:policy/service-role/AmazonRekognitionServiceRole"
}
`)
}

func testAccStreamProcessorConfig_faceSearchCreate(rName string) string {
	return acctest.ConfigCompose(testAccStreamProcessorConfig_baseIAMRole(), fmt.Sprintf(`
resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
	kinesis_video_stream {
		arn = aws_kinesis_video_stream.test.arn
	}
  }	
  output {
	kinesis_data_stream {
		arn = aws_kinesis_stream.test.arn
	}
  }
  role_arn 				  = aws_iam_role.test.arn
  settings {
	connected_home
  }
}
`, rName))
}

func testAccStreamProcessorConfig_basic(rName, version string) string {
	return fmt.Sprintf(`
resource "aws_kinesis_video_stream" "test" {
  name = "test"
}

resource "aws_kinesis_stream" "test" {
  name             = "test"
  stream_mode_details {
    stream_mode = "ON_DEMAND"
  }
}

resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "rekognition.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_rekognition_stream_processor" "test" {
  name             = %[1]q
  input {
	kinesis_video_stream {
		arn = aws_kinesis_video_stream.test.arn
	}
  }	
  output {
	kinesis_data_stream {
		arn = aws_kinesis_stream.test.arn
	}
  }
  role_arn 				  = 
  settings {

  }
}
`, rName, version)
}
