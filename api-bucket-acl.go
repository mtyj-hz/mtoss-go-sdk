package minio

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7/pkg/s3utils"
)

type AccessControlPolicy struct {
	XMLName xml.Name `xml:"AccessControlPolicy"`
	Owner
	AccessControlList struct {
		Grant struct {
			Grantee struct {
				Type string `xml:"Type"`
			} `xml:"Grantee"`
			Permission string `xml:"Permission"`
		} `xml:"Grant"`
	} `xml:"AccessControlList"`
}

func (c *Client) GetBucketAcl(ctx context.Context, bucketName string) (string, error) {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "", err
	}
	urlValues := make(url.Values)
	urlValues.Set("acl", "")

	reqMetadata := requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	}
	resp, err := c.executeMethod(ctx, http.MethodGet, reqMetadata)
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

func (c *Client) PutBucketAcl(ctx context.Context, bucketName, acl string) error {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	urlValues := make(url.Values)
	urlValues.Set("acl", "")

	head := make(http.Header)
	if acl = strings.Trim(acl, " "); acl == "" {
		acl = Default
	}
	head.Set(amzACL, acl)

	reqMetadata := requestMetadata{
		bucketName:   bucketName,
		customHeader: head,
		queryValues:  urlValues,
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
