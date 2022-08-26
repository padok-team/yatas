package acm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

type ACMGetObjectAPI interface {
	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)
}

func GetCertificates(svc ACMGetObjectAPI) []types.CertificateDetail {
	input := &acm.ListCertificatesInput{}
	result, err := svc.ListCertificates(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var certificatesArn []*string
	var certificates []types.CertificateDetail
	for _, r := range result.CertificateSummaryList {
		certificatesArn = append(certificatesArn, r.CertificateArn)
	}
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.ListCertificates(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		for _, r := range result.CertificateSummaryList {
			certificatesArn = append(certificatesArn, r.CertificateArn)
		}
	}

	for _, c := range certificatesArn {
		input := &acm.DescribeCertificateInput{
			CertificateArn: c,
		}
		result, err := svc.DescribeCertificate(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		certificates = append(certificates, *result.Certificate)
	}
	return certificates

}
