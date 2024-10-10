// Code generated by smithy-go-codegen DO NOT EDIT.

package efs

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Updates protection on the file system.
//
// This operation requires permissions for the
// elasticfilesystem:UpdateFileSystemProtection action.
func (c *Client) UpdateFileSystemProtection(ctx context.Context, params *UpdateFileSystemProtectionInput, optFns ...func(*Options)) (*UpdateFileSystemProtectionOutput, error) {
	if params == nil {
		params = &UpdateFileSystemProtectionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateFileSystemProtection", params, optFns, c.addOperationUpdateFileSystemProtectionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateFileSystemProtectionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateFileSystemProtectionInput struct {

	// The ID of the file system to update.
	//
	// This member is required.
	FileSystemId *string

	// The status of the file system's replication overwrite protection.
	//
	//   - ENABLED – The file system cannot be used as the destination file system in a
	//   replication configuration. The file system is writeable. Replication overwrite
	//   protection is ENABLED by default.
	//
	//   - DISABLED – The file system can be used as the destination file system in a
	//   replication configuration. The file system is read-only and can only be modified
	//   by EFS replication.
	//
	//   - REPLICATING – The file system is being used as the destination file system
	//   in a replication configuration. The file system is read-only and is only
	//   modified only by EFS replication.
	//
	// If the replication configuration is deleted, the file system's replication
	// overwrite protection is re-enabled, the file system becomes writeable.
	ReplicationOverwriteProtection types.ReplicationOverwriteProtection

	noSmithyDocumentSerde
}

// Describes the protection on a file system.
type UpdateFileSystemProtectionOutput struct {

	// The status of the file system's replication overwrite protection.
	//
	//   - ENABLED – The file system cannot be used as the destination file system in a
	//   replication configuration. The file system is writeable. Replication overwrite
	//   protection is ENABLED by default.
	//
	//   - DISABLED – The file system can be used as the destination file system in a
	//   replication configuration. The file system is read-only and can only be modified
	//   by EFS replication.
	//
	//   - REPLICATING – The file system is being used as the destination file system
	//   in a replication configuration. The file system is read-only and is only
	//   modified only by EFS replication.
	//
	// If the replication configuration is deleted, the file system's replication
	// overwrite protection is re-enabled, the file system becomes writeable.
	ReplicationOverwriteProtection types.ReplicationOverwriteProtection

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateFileSystemProtectionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpUpdateFileSystemProtection{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpUpdateFileSystemProtection{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateFileSystemProtection"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateFileSystemProtectionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateFileSystemProtection(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateFileSystemProtection(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateFileSystemProtection",
	}
}