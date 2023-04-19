package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/stretchr/testify/assert"
)

type MockIMDSClient struct {
	GetInstanceIdentityDocumentFunc func(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

func (m *MockIMDSClient) GetInstanceIdentityDocument(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
	return m.GetInstanceIdentityDocumentFunc(ctx, params, optFns...)
}

func TestCallAws(t *testing.T) {
	tests := []struct {
		name                    string
		mockGetInstanceIdentity func(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
		expectedInstanceID      string
		expectedError           error
	}{
		{
			name: "valid instance identity document",
			mockGetInstanceIdentity: func(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
				return &imds.GetInstanceIdentityDocumentOutput{
					InstanceIdentityDocument: imds.InstanceIdentityDocument{
						InstanceID: "i-1234567890abcdef0",
					},
				}, nil
			},
			expectedInstanceID: "i-1234567890abcdef0",
			expectedError:      nil,
		},
		{
			name: "error getting instance identity document",
			mockGetInstanceIdentity: func(ctx context.Context, params *imds.GetInstanceIdentityDocumentInput, optFns ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
				return nil, errors.New("request failed")
			},
			expectedInstanceID: "",
			expectedError:      errors.New("request failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			NewIMDSClient = func(cfg aws.Config) IMDSClient {
				return &MockIMDSClient{
					GetInstanceIdentityDocumentFunc: tc.mockGetInstanceIdentity,
				}
			}
			defer func() { NewIMDSClient = func(cfg aws.Config) IMDSClient { return imds.NewFromConfig(cfg) } }()

			client := http.Client{}
			instanceID, err := CallAws(client)

			assert.Equal(t, tc.expectedInstanceID, instanceID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
