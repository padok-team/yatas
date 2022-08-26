// Package policy provides a custom function to unmarshal AWS policies.
package iam

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// UnmarshalJSON decodifies input JSON info to awsPolicy type
func (policyJSON *Policy) UnmarshalJSON(policy []byte) error {

	var raw interface{}
	var err error
	var statementList []Statement

	err = json.Unmarshal(policy, &raw)
	if err != nil {
		return err
	}
	// Parsing content of JSON element as empty interface
	switch object := raw.(type) {
	// All elelements
	case map[string]interface{}:
		for key, value := range object {
			switch key {
			case "Version":
				policyJSON.Version = value.(string)
			case "ID":
				policyJSON.ID = value.(string)
			case "Statement":
				statementList = make([]Statement, 0)
				// Statement level - slice -> []interface{} , single element -> map[string]interface
				switch statement := value.(type) {
				// Statement slice -> iterate over elements, parse and store into slice
				case []interface{}:
					// statement slice
					// iterate over statements
					for _, statementValue := range statement {
						statement := Statement{}
						// Type assertion to format info
						statementMap := statementValue.(map[string]interface{})
						// Parse statement
						statement.Parse(statementMap)
						// Append statement to slice
						statementList = append(statementList, statement)
					}
				// Single statement -> parse and store it into slice
				case map[string]interface{}:
					statementMap := Statement{}
					// Parse statement
					statementMap.Parse(statement)
					statementList = append(statementList, statementMap)
				}
				// Assign statements slice to Policy
				policyJSON.Statements = statementList
			}
		}
	}
	return err
}

// Parse decodifies input JSON info into Statement type
func (statementJSON *Statement) Parse(statement map[string]interface{}) {

	// Definitions
	var principal, notPrincipal, action, notAction, resource, notResource, condition []string
	var err error

	/* Iterate over map elements, each key element (statementKey) is the statement element
	identifer and each value element (statementValue) the statement element value */
	for statementKey, statementValue := range statement {
		// Switch case over key type (identifying Statement elements)
		switch statementKey {
		case "StatementID":
			// Type assertion to assign
			statementJSON.StatementID = statementValue.(string)
		case "Effect":
			// Type assertion to assign
			statementJSON.Effect = statementValue.(string)
		case "Principal":
			// principal(statementValue) can be map[string][]string/string -> needs processing
			// Initialize map
			statementJSON.Principal = make(map[string][]string)
			// procesing map
			mapStatement := statementValue.(map[string]interface{})
			// iterate over key principal (keyPrincipal) and value principal (valuePrincipal)
			for keyPrincipal, valuePrincipal := range mapStatement {
				// valuePrincipal can be string or []string
				switch valuePrincipal := valuePrincipal.(type) {
				case string:
					// As map each element is identified with a key and has a value
					principal = make([]string, 0)
					statementJSON.Principal[keyPrincipal] = append(principal, valuePrincipal)
				case []interface{}:
					/* If value is an interface we know we have an []string -> knowing final type
					we can use mapstructure (which uses reflect) to store as final type */
					err = mapstructure.Decode(statementValue, &statementJSON.Principal)
					if err != nil {
						log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement principal element").Err(err).Msg("")
					}
				}
			}
		case "NotPrincipal":
			// Same case as principal
			// notprincipal has to be statementValue = map[string][]string/string -> needs processing
			// Same procedure as Principal
			// Intialize map
			statementJSON.NotPrincipal = make(map[string][]string)
			// procesing map (statementValue)
			mapStatement := statementValue.(map[string]interface{})
			for keyNotPrincipal, valueNotPrincipal := range mapStatement {
				// valueNotPrincipal can be string or []string
				switch vnp := valueNotPrincipal.(type) {
				case string:
					notPrincipal = make([]string, 0)
					statementJSON.NotPrincipal[keyNotPrincipal] = append(notPrincipal, vnp)
				case []interface{}:
					err = mapstructure.Decode(statementValue, &statementJSON.NotPrincipal)
					if err != nil {
						log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement not principal element").Err(err).Msg("")
					}
				}
			}
		case "Action":
			// We only have now string or []string, process with type assertion and mapstructure
			// Action can be string or []string
			switch statementValue := statementValue.(type) {
			case string:
				action = make([]string, 0)
				statementJSON.Action = append(action, statementValue)
			case []interface{}:
				err = mapstructure.Decode(statementValue, &statementJSON.Action)
				if err != nil {
					log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement action element").Err(err).Msg("")
				}
			}
		case "NotAction":
			// Same as Action
			// NotAction can be string or []string
			switch statementValue := statementValue.(type) {
			case string:
				notAction = make([]string, 0)
				statementJSON.NotAction = append(notAction, statementValue)
			case []interface{}:
				err = mapstructure.Decode(statementValue, &statementJSON.NotAction)
				if err != nil {
					log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement not action element").Err(err).Msg("")
				}
			}
		case "Resource":
			// Same as Action
			// Resource can be string or []string
			switch statementValue := statementValue.(type) {
			case string:
				resource = make([]string, 0)
				statementJSON.Resource = append(resource, statementValue)
			case []interface{}:
				err = mapstructure.Decode(statementValue, &statementJSON.Resource)
				if err != nil {
					log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement resource element").Err(err).Msg("")
				}
			}
		case "NotResource":
			// Same as Action
			// NotResource can be string or []string
			switch statementValue := statementValue.(type) {
			case string:
				notResource = make([]string, 0)
				statementJSON.NotResource = append(notResource, statementValue)
			case []interface{}:
				err = mapstructure.Decode(statementValue, &statementJSON.NotResource)
				if err != nil {
					log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement not resource element").Err(err).Msg("")
				}
			}
		case "Condition":
			// Condition can be string, []string or map(lot of options)
			switch statementValue := statementValue.(type) {
			case string:
				condition = make([]string, 0)
				statementJSON.Condition = append(condition, statementValue)
			case []interface{}:
				err = mapstructure.Decode(statementValue, &statementJSON.Condition)
				if err != nil {
					log.Error().Str("Error parsing policies", "Error using mapstructure parsing Policy statement condition element").Err(err).Msg("")
				}
			// If map format as raw text and store it as string
			case map[string]interface{}:
				condition = make([]string, 0)
				statementJSON.Condition = append(condition, fmt.Sprintf("%v", statementValue))
			}
		}
	}
}
