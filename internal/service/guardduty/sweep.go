//go:build sweep
// +build sweep

package guardduty

import (
	"fmt"
	"log"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/guardduty"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
	"github.com/PixarV/terraform-provider-ritt/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_guardduty_detector", &resource.Sweeper{
		Name:         "aws_guardduty_detector",
		F:            sweepDetectors,
		Dependencies: []string{"aws_guardduty_publishing_destination"},
	})

	resource.AddTestSweepers("aws_guardduty_publishing_destination", &resource.Sweeper{
		Name: "aws_guardduty_publishing_destination",
		F:    sweepPublishingDestinations,
	})
}

func sweepDetectors(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).GuardDutyConn
	input := &guardduty.ListDetectorsInput{}
	var sweeperErrs *multierror.Error

	err = conn.ListDetectorsPages(input, func(page *guardduty.ListDetectorsOutput, lastPage bool) bool {
		for _, detectorID := range page.DetectorIds {
			id := aws.StringValue(detectorID)
			input := &guardduty.DeleteDetectorInput{
				DetectorId: detectorID,
			}

			log.Printf("[INFO] Deleting GuardDuty Detector: %s", id)
			_, err := conn.DeleteDetector(input)
			if tfawserr.ErrCodeContains(err, "AccessDenied") {
				log.Printf("[WARN] Skipping GuardDuty Detector (%s): %s", id, err)
				continue
			}
			if err != nil {
				sweeperErr := fmt.Errorf("error deleting GuardDuty Detector (%s): %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping GuardDuty Detector sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error retrieving GuardDuty Detectors: %w", err)
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepPublishingDestinations(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).GuardDutyConn
	var sweeperErrs *multierror.Error

	detect_input := &guardduty.ListDetectorsInput{}

	err = conn.ListDetectorsPages(detect_input, func(page *guardduty.ListDetectorsOutput, lastPage bool) bool {
		for _, detectorID := range page.DetectorIds {
			list_input := &guardduty.ListPublishingDestinationsInput{
				DetectorId: detectorID,
			}

			err = conn.ListPublishingDestinationsPages(list_input, func(page *guardduty.ListPublishingDestinationsOutput, lastPage bool) bool {
				for _, destination_element := range page.Destinations {
					input := &guardduty.DeletePublishingDestinationInput{
						DestinationId: destination_element.DestinationId,
						DetectorId:    detectorID,
					}

					log.Printf("[INFO] Deleting GuardDuty Publishing Destination: %s", *destination_element.DestinationId)
					_, err := conn.DeletePublishingDestination(input)

					if err != nil {
						sweeperErr := fmt.Errorf("error deleting GuardDuty Publishing Destination (%s): %w", *destination_element.DestinationId, err)
						log.Printf("[ERROR] %s", sweeperErr)
						sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
					}
				}
				return !lastPage
			})
		}
		return !lastPage
	})

	if err != nil {
		sweeperErr := fmt.Errorf("Error receiving Guardduty detectors for publishing sweep : %w", err)
		log.Printf("[ERROR] %s", sweeperErr)
		sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
	}

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping GuardDuty Publishing Destination sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error retrieving GuardDuty Publishing Destinations: %s", err)
	}

	return sweeperErrs.ErrorOrNil()
}
