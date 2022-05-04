//go:build generate
// +build generate

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

const filename = `create_tags_gen.go`

var (
	retryCreateOnNotFound = flag.Bool("RetryCreateOnNotFound", true, "retry create if resource not found")
	tagOp                 = flag.String("TagOp", "CreateTags", "tag function")
	tagOpBatchSize        = flag.String("TagOpBatchSize", "", "tag function batch size")
	tagInCustomVal        = flag.String("TagInCustomVal", "", "tag input custom value")
	tagInIDElem           = flag.String("TagInIDElem", "Resources", "tag input identifier field")
	tagInIDNeedSlice      = flag.Bool("TagInIDNeedSlice", true, "tag input identifier requires slice")
	tagInTagsElem         = flag.String("TagInTagsElem", "Tags", "tag input tags field")
	tagResTypeElem        = flag.String("TagResTypeElem", "", "tag resource type field")
	tagTypeIDElem         = flag.String("TagTypeIDElem", "", "tag type identifier field")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tmain.go [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

type TemplateData struct {
	AWSService     string
	ClientType     string
	ServicePackage string

	ParentNotFoundError   string
	RetryCreateOnNotFound bool
	TagInCustomVal        string
	TagInIDElem           string
	TagInIDNeedSlice      bool
	TagInTagsElem         string
	TagOp                 string
	TagOpBatchSize        string
	TagResTypeElem        string
	TagTypeIDElem         string
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	templateData := TemplateData{
		AWSService:            "ec2",
		ServicePackage:        "ec2",
		ClientType:            "*ec2.EC2",
		RetryCreateOnNotFound: *retryCreateOnNotFound,
		TagInCustomVal:        *tagInCustomVal,
		TagInIDElem:           *tagInIDElem,
		TagInIDNeedSlice:      *tagInIDNeedSlice,
		TagInTagsElem:         *tagInTagsElem,
		TagOp:                 *tagOp,
		TagOpBatchSize:        *tagOpBatchSize,
		TagResTypeElem:        *tagResTypeElem,
		TagTypeIDElem:         *tagTypeIDElem,
	}

	if templateData.ServicePackage == "ec2" {
		templateData.ParentNotFoundError = `
if tfawserr.ErrCodeContains(err, ".NotFound") {
	err = &resource.NotFoundError{
		LastError:   err,
		LastRequest: input,
	}
}
`
	}

	tmpl, err := template.New("createtags").Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by internal/service/ec2/generate/createtags/main.go; DO NOT EDIT.

package {{ .ServicePackage }}

import (
	"fmt"
	"time"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/{{ .AWSService }}"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
	tftags "github.com/PixarV/terraform-provider-ritt/internal/tags"
)

const EventualConsistencyTimeout = 5 * time.Minute

// CreateTags creates {{ .ServicePackage }} service tags for new resources.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func CreateTags(conn {{ .ClientType }}, identifier string{{ if .TagResTypeElem }}, resourceType string{{ end }}, tagsMap interface{}) error {
	tags := tftags.New(tagsMap)

	{{- if .TagOpBatchSize }}
	for _, tags := range tags.Chunks({{ .TagOpBatchSize }}) {
	{{- end }}
	input := &{{ .AWSService }}.{{ .TagOp }}Input{
		{{- if not ( .TagTypeIDElem ) }}
		{{- if .TagInIDNeedSlice }}
		{{ .TagInIDElem }}: aws.StringSlice([]string{identifier}),
		{{- else }}
		{{ .TagInIDElem }}: aws.String(identifier),
		{{- end }}
		{{- if .TagResTypeElem }}
		{{ .TagResTypeElem }}:    aws.String(resourceType),
		{{- end }}
		{{- end }}
		{{- if .TagInCustomVal }}
		{{ .TagInTagsElem }}:       {{ .TagInCustomVal }},
		{{- else }}
		{{ .TagInTagsElem }}:       Tags(tags.IgnoreAWS()),
		{{- end }}
	}

	{{- if .RetryCreateOnNotFound }}

	_, err := tfresource.RetryWhenNotFound(EventualConsistencyTimeout, func() (interface{}, error) {
		output, err := conn.{{ .TagOp }}(input)

		{{ .ParentNotFoundError }}

		return output, err
	})
	{{- else }}
	_, err := conn.{{ .TagOp }}(input)
	{{- end }}

	if err != nil {
		return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
	}
	{{- if .TagOpBatchSize }}
	}
	{{- end }}

	return nil
}
`
