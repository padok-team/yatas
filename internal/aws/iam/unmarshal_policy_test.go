// Package policy provides a custom function to unmarshal AWS policies.
package iam

import "testing"

func TestPolicy_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Version    string
		ID         string
		Statements []Statement
	}
	type args struct {
		policy []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UnmarshalJSON",
			fields: fields{
				Version:    "",
				ID:         "",
				Statements: nil,
			},
			args: args{
				policy: []byte(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Sid": "",
							"Effect": "Allow",
							"Action": "*",
							"Resource": "*",
							"NotAction": "*",
							"Condition": {
								"IpAddress": {
									"aws:SourceIp": ""
								}
							},
							"Principal": {
								"AWS": "*"
							}
						}
					]
				}`),
			},
			wantErr: false,
		},
		{
			name: "UnmarshalJSON",
			fields: fields{
				Version:    "",
				ID:         "",
				Statements: nil,
			},
			args: args{
				policy: []byte(`{
					"Version": "2012-10-17",
					"Statement": 
						{	
							"StatementID": "",
							"Sid": "",
							"Effect": "Allow",
							"Action": "*",
							"Resource": "*",
							"NotPrincipal": {
								"AWS": "*"
							}
						}
					
				}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policyJSON := &Policy{
				Version:    tt.fields.Version,
				ID:         tt.fields.ID,
				Statements: tt.fields.Statements,
			}
			if err := policyJSON.UnmarshalJSON(tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("Policy.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicy_UnmarshalJSONFail(t *testing.T) {
	type fields struct {
		Version    string
		ID         string
		Statements []Statement
	}
	type args struct {
		policy []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UnmarshalJSON",
			fields: fields{
				Version:    "",
				ID:         "",
				Statements: nil,
			},
			args: args{
				policy: []byte(`{
					"Version": "2012-10-17",
					"ID": "",
					"Statements": [
						{
							"Sid": "",
							"Effect": "Allow",
							"Action": "*",
							"Resource": "*"
						}
					]
				}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policyJSON := &Policy{
				Version:    tt.fields.Version,
				ID:         tt.fields.ID,
				Statements: tt.fields.Statements,
			}
			if err := policyJSON.UnmarshalJSON(tt.args.policy); (err == nil) != tt.wantErr {
				t.Errorf("Policy.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
