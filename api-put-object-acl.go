package mtoss_go_sdk

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/mtyj-hz/mtoss-go-sdk/pkg/s3utils"
)

func (c *Client) PutObjectACL(ctx context.Context, bucketName, objectName, acl string, opts ObjectAclMtOpts) error {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return err
	}

	urlValues := make(url.Values)
	urlValues.Set("acl", "")
	if opts.VersionId != "" {
		urlValues.Set("versionId", opts.VersionId)
	}
	head := make(http.Header)
	if acl = strings.Trim(acl, " "); acl == "" {
		acl = Default
	}
	head.Set(amzACL, acl)

	reqMetadata := requestMetadata{
		bucketName:   bucketName,
		objectName:   objectName,
		queryValues:  urlValues,
		customHeader: head,
	}
	resp, err := c.executeMethod(ctx, http.MethodPut, reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return httpRespToErrorResponse(resp, bucketName, "")
		}
	}
	return nil
}
