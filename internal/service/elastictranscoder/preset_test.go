package elastictranscoder_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/elastictranscoder"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/PixarV/terraform-provider-ritt/internal/acctest"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
	tfet "github.com/PixarV/terraform-provider-ritt/internal/service/elastictranscoder"
)

func TestAccElasticTranscoderPreset_basic(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "elastictranscoder", regexp.MustCompile(`preset/.+`)),
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

func TestAccElasticTranscoderPreset_video_noCodec(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetVideoNoCodec(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
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

//https://github.com/terraform-providers/terraform-provider-aws/issues/14090
func TestAccElasticTranscoderPreset_audio_noBitRate(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetNoBitRateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
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

func TestAccElasticTranscoderPreset_disappears(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					acctest.CheckResourceDisappears(acctest.Provider, tfet.ResourcePreset(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Reference: https://github.com/PixarV/terraform-provider-ritt/issues/14087
func TestAccElasticTranscoderPreset_AudioCodecOptions_empty(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetAudioCodecOptionsEmptyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"audio_codec_options"}, // Due to incorrect schema (should be nested under audio)
			},
		},
	})
}

func TestAccElasticTranscoderPreset_description(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetDescriptionConfig(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
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

// Tests all configuration blocks
func TestAccElasticTranscoderPreset_full(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetFull1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					resource.TestCheckResourceAttr(resourceName, "audio.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio_codec_options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thumbnails.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_codec_options.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "video_watermarks.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPresetFull2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					resource.TestCheckResourceAttr(resourceName, "audio.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio_codec_options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "thumbnails.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_codec_options.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "video_watermarks.#", "1"),
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

// Reference: https://github.com/PixarV/terraform-provider-ritt/issues/695
func TestAccElasticTranscoderPreset_Video_frameRate(t *testing.T) {
	var preset elastictranscoder.Preset
	resourceName := "aws_elastictranscoder_preset.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, elastictranscoder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckElasticTranscoderPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPresetVideoFrameRateConfig(rName, "29.97"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticTranscoderPresetExists(resourceName, &preset),
					resource.TestCheckResourceAttr(resourceName, "video.0.frame_rate", "29.97"),
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

func testAccCheckElasticTranscoderPresetExists(name string, preset *elastictranscoder.Preset) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticTranscoderConn

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Preset ID is set")
		}

		out, err := conn.ReadPreset(&elastictranscoder.ReadPresetInput{
			Id: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		*preset = *out.Preset

		return nil
	}
}

func testAccCheckElasticTranscoderPresetDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ElasticTranscoderConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_elastictranscoder_preset" {
			continue
		}

		out, err := conn.ReadPreset(&elastictranscoder.ReadPresetInput{
			Id: aws.String(rs.Primary.ID),
		})

		if err == nil {
			if out.Preset != nil && aws.StringValue(out.Preset.Id) == rs.Primary.ID {
				return fmt.Errorf("Elastic Transcoder Preset still exists")
			}
		}

		if !tfawserr.ErrCodeEquals(err, elastictranscoder.ErrCodeResourceNotFoundException) {
			return fmt.Errorf("unexpected error: %s", err)
		}

	}
	return nil
}

func testAccPresetConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "mp4"
  name      = %[1]q

  audio {
    audio_packing_mode = "SingleTrack"
    bit_rate           = 320
    channels           = 2
    codec              = "mp3"
    sample_rate        = 44100
  }
}
`, rName)
}

func testAccPresetAudioCodecOptionsEmptyConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "mp4"
  name      = %[1]q

  audio {
    audio_packing_mode = "SingleTrack"
    bit_rate           = 320
    channels           = 2
    codec              = "mp3"
    sample_rate        = 44100
  }

  audio_codec_options {}
}
`, rName)
}

func testAccPresetDescriptionConfig(rName string, description string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container   = "mp4"
  description = %[2]q
  name        = %[1]q

  audio {
    audio_packing_mode = "SingleTrack"
    bit_rate           = 320
    channels           = 2
    codec              = "mp3"
    sample_rate        = 44100
  }
}
`, rName, description)
}

func testAccPresetFull1Config(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "mp4"
  name      = %[1]q

  audio {
    audio_packing_mode = "SingleTrack"
    bit_rate           = 128
    channels           = 2
    codec              = "AAC"
    sample_rate        = 48000
  }

  audio_codec_options {
    profile = "auto"
  }

  video {
    bit_rate             = "auto"
    codec                = "H.264"
    display_aspect_ratio = "16:9"
    fixed_gop            = "true"
    frame_rate           = "auto"
    keyframes_max_dist   = 90
    max_height           = 1080
    max_width            = 1920
    padding_policy       = "Pad"
    sizing_policy        = "Fit"
  }

  video_codec_options = {
    Profile                  = "main"
    Level                    = "4.1"
    MaxReferenceFrames       = 4
    InterlacedMode           = "Auto"
    ColorSpaceConversionMode = "None"
  }

  thumbnails {
    format         = "jpg"
    interval       = 5
    max_width      = 960
    max_height     = 540
    padding_policy = "Pad"
    sizing_policy  = "Fit"
  }
}
`, rName)
}

func testAccPresetFull2Config(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "mp4"
  name      = %[1]q

  audio {
    audio_packing_mode = "SingleTrack"
    bit_rate           = 96
    channels           = 2
    codec              = "AAC"
    sample_rate        = 44100
  }

  audio_codec_options {
    profile = "AAC-LC"
  }

  video {
    bit_rate             = "1600"
    codec                = "H.264"
    display_aspect_ratio = "16:9"
    fixed_gop            = "false"
    frame_rate           = "auto"
    max_frame_rate       = "60"
    keyframes_max_dist   = 240
    max_height           = "auto"
    max_width            = "auto"
    padding_policy       = "Pad"
    sizing_policy        = "Fit"
  }

  video_codec_options = {
    Profile                  = "main"
    Level                    = "2.2"
    MaxReferenceFrames       = 3
    InterlacedMode           = "Progressive"
    ColorSpaceConversionMode = "None"
  }

  video_watermarks {
    id                = "Terraform Test"
    max_width         = "20%%"
    max_height        = "20%%"
    sizing_policy     = "ShrinkToFit"
    horizontal_align  = "Right"
    horizontal_offset = "10px"
    vertical_align    = "Bottom"
    vertical_offset   = "10px"
    opacity           = "55.5"
    target            = "Content"
  }

  thumbnails {
    format         = "png"
    interval       = 120
    max_width      = "auto"
    max_height     = "auto"
    padding_policy = "Pad"
    sizing_policy  = "Fit"
  }
}
`, rName)
}

func testAccPresetVideoFrameRateConfig(rName string, frameRate string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "mp4"
  name      = %[1]q

  thumbnails {
    format         = "png"
    interval       = 120
    max_width      = "auto"
    max_height     = "auto"
    padding_policy = "Pad"
    sizing_policy  = "Fit"
  }

  video {
    bit_rate             = "auto"
    codec                = "H.264"
    display_aspect_ratio = "16:9"
    fixed_gop            = "true"
    frame_rate           = %[2]q
    keyframes_max_dist   = 90
    max_height           = 1080
    max_width            = 1920
    padding_policy       = "Pad"
    sizing_policy        = "Fit"
  }

  video_codec_options = {
    Profile                  = "main"
    Level                    = "4.1"
    MaxReferenceFrames       = 4
    InterlacedMode           = "Auto"
    ColorSpaceConversionMode = "None"
  }
}
`, rName, frameRate)
}

func testAccPresetVideoNoCodec(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "webm"
  type      = "Custom"
  name      = %[1]q

  audio {
    codec              = "vorbis"
    sample_rate        = 44100
    bit_rate           = 128
    channels           = 2
    audio_packing_mode = "SingleTrack"
  }

  thumbnails {
    format         = "png"
    interval       = 300
    max_width      = 640
    max_height     = 360
    sizing_policy  = "ShrinkToFit"
    padding_policy = "NoPad"
  }

  video {
    codec                = "vp9"
    keyframes_max_dist   = 90
    fixed_gop            = false
    bit_rate             = 600
    frame_rate           = 30
    max_width            = 640
    max_height           = 360
    display_aspect_ratio = "auto"
    sizing_policy        = "Fit"
    padding_policy       = "NoPad"
  }
}
`, rName)
}

func testAccPresetNoBitRateConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_elastictranscoder_preset" "test" {
  container = "wav"
  name      = %[1]q
  audio {
    audio_packing_mode = "SingleTrack"
    channels           = 2
    codec              = "pcm"
    sample_rate        = 44100
  }
}
`, rName)
}
