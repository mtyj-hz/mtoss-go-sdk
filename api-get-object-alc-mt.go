package minio

import (
	"context"
	"net/http"
	"net/url"

	"github.com/minio/minio-go/v7/pkg/s3utils"
)

type ObjectAclMtOpts struct {
	VersionId string
}

func (c *Client) GetObjectMtACL(ctx context.Context, bucketName, objectName string, opts ObjectAclMtOpts) (string, error) {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "", err
	}
	if err := s3utils.CheckValidObjectName(bucketName); err != nil {
		return "", err
	}
	urlValues := make(url.Values)
	urlValues.Set("acl", "")
	if opts.VersionId != "" {
		urlValues.Set("versionId", opts.VersionId)
	}
	resp, err := c.executeMethod(ctx, http.MethodGet, requestMetadata{
		bucketName:  bucketName,
		objectName:  objectName,
		queryValues: urlValues,
	})
	if err != nil {
		return "", err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return "", httpRespToErrorResponse(resp, bucketName, "")
	}

	res := &AccessControlPolicy{}

	if err := xmlDecoder(resp.Body, res); err != nil {
		return "", err
	}

	return res.AccessControlList.Grant.Permission, nil
}
