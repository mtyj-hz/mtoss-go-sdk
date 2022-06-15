package mtoss_go_sdk

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/mtyj-hz/mtoss-go-sdk/pkg/s3utils"
)

func (c *Client) PutBucketLogging(ctx context.Context, bucketName string, config BucketLoggingRet) error {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	buf, err := xml.Marshal(config)
	if err != nil {
		return err
	}

	urlValues := make(url.Values)
	urlValues.Set("logging", "")

	reqMetadata := requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentBody:      bytes.NewReader(buf),
		contentLength:    int64(len(buf)),
		contentMD5Base64: sumMD5Base64(buf),
		contentSHA256Hex: sum256Hex(buf),
	}

	// Execute PUT to set a bucket versioning.
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

type BucketLoggingRet struct {
	XMLName xml.Name              `xml:"BucketLoggingStatus"`
	Enabled *BucketLoggingEnabled `xml:"LoggingEnabled"`
}
type BucketLoggingEnabled struct {
	TargetBucket string `xml:"TargetBucket"`
	TargetPrefix string `xml:"TargetPrefix"`
}

func (c *Client) GetBucketLogging(ctx context.Context, bucketName string) (*BucketLoggingRet, error) {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return nil, err
	}
	urlValues := make(url.Values)
	urlValues.Set("logging", "")

	reqMetadata := requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	}
	resp, err := c.executeMethod(ctx, http.MethodGet, reqMetadata)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp, bucketName, "")
	}

	res := &BucketLoggingRet{}

	if err := xmlDecoder(resp.Body, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteBucketLogging(ctx context.Context, bucketName string) error {
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	urlValues := make(url.Values)
	urlValues.Set("logging", "")

	reqMetadata := requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	}
	resp, err := c.executeMethod(ctx, http.MethodDelete, reqMetadata)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp, bucketName, "")
	}

	return nil
}
