package docdb

import (
	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/docdb"
)

// Takes the result of flatmap.Expand for an array of parameters and
// returns Parameter API compatible objects
func expandParameters(configured []interface{}) []*docdb.Parameter {
	parameters := make([]*docdb.Parameter, 0, len(configured))

	// Loop over our configured parameters and create
	// an array of aws-sdk-go compatible objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		p := &docdb.Parameter{
			ApplyMethod:    aws.String(data["apply_method"].(string)),
			ParameterName:  aws.String(data["name"].(string)),
			ParameterValue: aws.String(data["value"].(string)),
		}

		parameters = append(parameters, p)
	}

	return parameters
}

// Flattens an array of Parameters into a []map[string]interface{}
func flattenParameters(list []*docdb.Parameter, parameterList []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		if i.ParameterValue != nil {
			name := aws.StringValue(i.ParameterName)

			// Check if any non-user parameters are specified in the configuration.
			parameterFound := false
			for _, configParameter := range parameterList {
				if configParameter.(map[string]interface{})["name"] == name {
					parameterFound = true
				}
			}

			// Skip parameters that are not user defined or specified in the configuration.
			if aws.StringValue(i.Source) != "user" && !parameterFound {
				continue
			}

			result = append(result, map[string]interface{}{
				"apply_method": aws.StringValue(i.ApplyMethod),
				"name":         aws.StringValue(i.ParameterName),
				"value":        aws.StringValue(i.ParameterValue),
			})
		}
	}
	return result
}
