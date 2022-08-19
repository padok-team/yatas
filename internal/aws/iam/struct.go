package iam

// Policy represents an AWS IAM policy document
type Policy struct {
	Version    string      `json:"Version"`
	ID         string      `json:"ID,omitempty"`
	Statements []Statement `json:"Statement"`
}

// Statement represents the body of an AWS IAM policy document
type Statement struct {
	StatementID  string              `json:"StatementID,omitempty"`  // Statement ID, service specific
	Effect       string              `json:"Effect"`                 // Allow or Deny
	Principal    map[string][]string `json:"Principal,omitempty"`    // principal that is allowed or denied
	NotPrincipal map[string][]string `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       []string            `json:"Action"`                 // allowed or denied action
	NotAction    []string            `json:"NotAction,omitempty"`    // matches everything except
	Resource     []string            `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  []string            `json:"NotResource,omitempty"`  // matches everything except
	Condition    []string            `json:"Condition,omitempty"`    // conditions for when a policy is in effect
}

type UserPolicies struct {
	UserName string
	Policies []Policy
}
