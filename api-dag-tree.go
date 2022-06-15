package mtoss_go_sdk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mtyj-hz/mtoss-go-sdk/pkg/s3utils"
)

type DagTreeOpts struct {
	VersionId string
	Cid       string
}

func (c *Client) GetDagTree(ctx context.Context, bucketName, objectName string, opts DagTreeOpts) ([]byte, error) {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return nil, err
	}
	if err := s3utils.CheckValidObjectName(bucketName); err != nil {
		return nil, err
	}
	urlValues := make(url.Values)
	urlValues.Set("dag", "")
	if opts.VersionId != "" {
		urlValues.Set("versionId", opts.VersionId)
	}
	if opts.Cid != "" {
		urlValues.Set("cid", opts.Cid)
	}

	reqMetadata := requestMetadata{
		bucketName:  bucketName,
		objectName:  objectName,
		queryValues: urlValues,
	}
	resp, err := c.executeMethod(ctx, http.MethodGet, reqMetadata)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return nil, httpRespToErrorResponse(resp, bucketName, "")
		}
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}
