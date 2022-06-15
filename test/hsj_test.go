package test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Cli        *minio.Client
	Core       *minio.Core
	Background context.Context
)

func TestMain(m *testing.M) {
	endpoint := "192.168.1.150:9000"
	// endpoint := "192.168.3.204:9000"
	//accessKeyID := "ZSfKOkJOzRlduaZ70sXJ"
	//secretAccessKey := "4EKyHtoIynJf8LeDs64qkd8an71OU5"
	//accessKeyID := "EG74PRI3L76SR46S38O3"
	//secretAccessKey := "UhYh3Uhbjsi4BeckcXWYlJD2AeKbdY9sN+BjsEBA"
	accessKeyID := "testsvc"
	secretAccessKey := "12345678"
	useSSL := false
	Background = context.Background()
	C, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	Cli = C
	if err != nil {
		fmt.Println(err.Error())
	}
	c, err := minio.NewCore(
		"192.168.1.184:9000",
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Core = c
	m.Run()
}

func TestMakeBucket(t *testing.T) {
	err := Cli.MakeBucket(Background, "hsj-test", minio.MakeBucketOptions{
		Region:       "test",
		StorageClass: "StorageClass",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
}

func TestBucketList(t *testing.T) {
	buckets, err := Cli.ListBuckets(Background)
	if err != nil {
		t.Log(err.Error())
		return
	}
	for i := range buckets {
		fmt.Println(buckets[i].StorageClass)
	}
}

func TestPutBucketLogging(t *testing.T) {
	err := Cli.PutBucketLogging(Background, "hsj-1234", minio.BucketLoggingRet{
		Enabled: &minio.BucketLoggingEnabled{
			TargetBucket: "hsj-1234",
			TargetPrefix: "",
		},
	})
	if err != nil {
		t.Log(err.Error())
		return
	}

}

func TestDeleteBucketLogging(t *testing.T) {
	if err := Cli.DeleteBucketLogging(Background, "hsj-1234"); err != nil {
		t.Log(err.Error())
		return
	}
}

func TestGetBucketLogging(t *testing.T) {
	res, err := Cli.GetBucketLogging(Background, "hsj-1234")
	if err != nil {
		t.Log(err.Error())
		return
	}
	fmt.Println(res.Enabled.TargetBucket)
}

func TestPutBucketAcl(t *testing.T) {
	if err := Cli.PutBucketAcl(Background, "hsj-1234", minio.PublicReadWrite); err != nil {
		t.Log(err.Error())
	}
}

func TestGetBucketAcl(t *testing.T) {
	acl, err := Cli.GetBucketAcl(Background, "hsj-1234")
	if err != nil {
		t.Log(err.Error())
		return
	}
	fmt.Println(acl)
}

func TestGetObjectACL(t *testing.T) {
	acl, err := Cli.GetObjectMtACL(Background, "hsj-1234", "7.gz")
	if err != nil {
		t.Log(err.Error())
		return
	}
	fmt.Println(acl)
}

func TestPutObjectACL(t *testing.T) {
	err := Cli.PutObjectACL(Background, "hsj-1234", "7.gz", minio.Private)
	if err != nil {
		t.Log(err.Error())
		return
	}
	//Cli.ListBuckets()
}

func TestListObjects(t *testing.T) {
	v2, err := Core.ListObjectsV2("mybucket111", "", "", "", "/", 4)
	if err != nil {
		t.Log(err.Error())
		return
	}
	fmt.Println(v2)
}

type Progress struct {
	Total   int64
	Current int64
	Percent int
}

func (p *Progress) Read(b []byte) (int, error) {
	if len(b) != 0 {
		p.Current += int64(len(b))
	}
	s := fmt.Sprintf("%.2f", float64(p.Current)/float64(p.Total))
	fmt.Println(s)
	return p.Percent, nil
}

func TestCore(t *testing.T) {
	buckets, err := Core.ListBuckets(Background)
	if err != nil {
		t.Log(err.Error())
		return
	}
	for i := range buckets {
		fmt.Println(buckets[i].Name)
	}
}

type UploadPartReq struct {
	PartNum int              // Number of the part uploaded.
	Part    minio.ObjectPart // Size of the part uploaded.
}

type UploadedPartRes struct {
	Error   error // Any error encountered while uploading the part.
	PartNum int   // Number of the part uploaded.
	Size    int64 // Size of the part uploaded.
	Part    minio.ObjectPart
}

func TestMuUpload(t *testing.T) {
	bucket := "hsj-123"
	open, err := os.Open("/Users/huangshijie/Downloads/原版.pdf")
	if err != nil {
		t.Log(err.Error())
		return
	}

	stat, err := open.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	size := stat.Size()
	opts := minio.PutObjectOptions{}
	totalPartsCount, partSize, lastPartSize, err := minio.OptimalPartInfo(size, opts.PartSize)
	if err != nil {
		t.Log(err.Error())
	}
	uploadID, err := Core.NewMultipartUpload(Background, bucket, stat.Name(), opts)
	if err != nil {
		t.Log(err.Error())
		return
	}
	defer func() {
		if err != nil {
			Core.AbortMultipartUpload(Background, bucket, stat.Name(), uploadID)
		}
	}()

	// Total data read and written to server. should be equal to 'size' at the end of the call.
	var totalUploadedSize int64

	// Complete multipart upload.
	completePart := make([]minio.CompletePart, 0)

	// Declare a channel that sends the next part number to be uploaded.
	// Buffered to 10000 because thats the maximum number of parts allowed
	// by S3.
	uploadPartsCh := make(chan UploadPartReq, 10000)

	// Declare a channel that sends back the response of a part upload.
	// Buffered to 10000 because thats the maximum number of parts allowed
	// by S3.
	uploadedPartsCh := make(chan UploadedPartRes, 10000)

	// Used for readability, lastPartNumber is always totalPartsCount.
	lastPartNumber := totalPartsCount

	// Send each part number to the channel to be processed.
	for p := 1; p <= totalPartsCount; p++ {
		uploadPartsCh <- UploadPartReq{PartNum: p}
	}
	close(uploadPartsCh)

	var partsBuf = make([][]byte, 4)
	for i := range partsBuf {
		partsBuf[i] = make([]byte, 0, partSize)
	}

	// Receive each part number from the channel allowing three parallel uploads.
	for w := 1; w <= 4; w++ {
		go func(w int, partSize int64) {
			// Each worker will draw from the part channel and upload in parallel.
			for uploadReq := range uploadPartsCh {

				// If partNumber was not uploaded we calculate the missing
				// part offset and size. For all other part numbers we
				// calculate offset based on multiples of partSize.
				readOffset := int64(uploadReq.PartNum-1) * partSize

				// As a special case if partNumber is lastPartNumber, we
				// calculate the offset based on the last part size.
				if uploadReq.PartNum == lastPartNumber {
					readOffset = (size - lastPartSize)
					partSize = lastPartSize
				}

				n, rerr := readFull(io.NewSectionReader(open, readOffset, partSize), partsBuf[w-1][:partSize])
				if rerr != nil && rerr != io.ErrUnexpectedEOF && err != io.EOF {
					uploadedPartsCh <- UploadedPartRes{
						Error: rerr,
					}
					// Exit the goroutine.
					return
				}

				// Get a section reader on a particular offset.
				hookReader := bytes.NewReader(partsBuf[w-1][:n])

				// Proceed to upload the part.
				// fmt.Println(fmt.Sprintf("线程number:%v,PartNum:%v",w,uploadReq.PartNum))
				objPart, err := Core.PutObjectPart(Background, bucket, stat.Name(), uploadID, uploadReq.PartNum,
					hookReader, partSize, "", "",
					opts.ServerSideEncryption)
				if err != nil {
					uploadedPartsCh <- UploadedPartRes{
						Error: err,
					}
					// Exit the goroutine.
					return
				}
				// Save successfully uploaded part metadata.
				uploadReq.Part = objPart

				// Send successful part info through the channel.
				uploadedPartsCh <- UploadedPartRes{
					Size:    objPart.Size,
					PartNum: uploadReq.PartNum,
					Part:    uploadReq.Part,
				}
			}
		}(w, partSize)
	}

	// Gather the responses as they occur and update any
	// progress bar.
	for u := 1; u <= totalPartsCount; u++ {
		uploadRes := <-uploadedPartsCh
		if uploadRes.Error != nil {
			// return minio.UploadInfo{}, uploadRes.Error
			t.Log(err.Error())
			return
		}
		// Update the totalUploadedSize.
		totalUploadedSize += uploadRes.Size
		// Store the parts to be completed in order.
		completePart = append(completePart, minio.CompletePart{
			ETag:       uploadRes.Part.ETag,
			PartNumber: uploadRes.Part.PartNumber,
		})
	}

	// Verify if we uploaded all the data.
	if totalUploadedSize != size {
		return
		// return UploadInfo{}, errUnexpectedEOF(totalUploadedSize, size, bucketName, objectName)
	}

	// Sort all completed parts.
	sort.Sort(completedParts(completePart))
	etag, err := Core.CompleteMultipartUpload(Background, bucket, stat.Name(), uploadID, completePart, opts)
	if err != nil {
		// return UploadInfo{}, err
		t.Log(err.Error())
		return
	}
	fmt.Println(etag)

}

var readFull = func(r io.Reader, buf []byte) (n int, err error) {
	// ReadFull reads exactly len(buf) bytes from r into buf.
	// It returns the number of bytes copied and an error if
	// fewer bytes were read. The error is EOF only if no bytes
	// were read. If an EOF happens after reading some but not
	// all the bytes, ReadFull returns ErrUnexpectedEOF.
	// On return, n == len(buf) if and only if err == nil.
	// If r returns an error having read at least len(buf) bytes,
	// the error is dropped.
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		// Some spurious io.Reader's return
		// io.ErrUnexpectedEOF when nn == 0
		// this behavior is undocumented
		// so we are on purpose not using io.ReadFull
		// implementation because this can lead
		// to custom handling, to avoid that
		// we simply modify the original io.ReadFull
		// implementation to avoid this issue.
		// io.ErrUnexpectedEOF with nn == 0 really
		// means that io.EOF
		if err == io.ErrUnexpectedEOF && nn == 0 {
			err = io.EOF
		}
		n += nn
	}
	if n >= len(buf) {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

type completedParts []minio.CompletePart

func (a completedParts) Len() int           { return len(a) }
func (a completedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

func EggFlower(src, dst string, container int) ([]string, int64, error) {
	if container < 8<<20 {
		container = 8 << 20
	}
	address := make([]string, 0)
	fileStream, err := os.Open(src)
	if err != nil {
		return address, 0, err
	}
	if !strings.HasSuffix(dst, "/") {
		return address, 0, errors.New("dst 格式错误，应该以'/'结尾")
	}
	stat, err := fileStream.Stat()
	if err != nil {
		return nil, 0, err
	}
	i := 0
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return address, stat.Size(), err
	}
	for {
		b := make([]byte, container)
		newReader := bufio.NewReader(fileStream)
		read, err := newReader.Read(b)
		if err != nil && err == io.EOF {
			return address, stat.Size(), err
		}
		if read == 0 {
			break
		}
		i++
		strBuff := strings.Builder{}
		strBuff.WriteString(dst)
		strBuff.WriteString(strconv.Itoa(i))
		strBuff.WriteString(".db")
		fileName := strBuff.String()
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return address, stat.Size(), err
		}
		writer := bufio.NewWriter(file)
		_, err = writer.Write(b[:read])
		if err != nil {
			return address, stat.Size(), err
		}
		if err := writer.Flush(); err != nil {
			return address, stat.Size(), err
		}
		if err := file.Close(); err != nil {
			return address, stat.Size(), err
		}
		address = append(address, fileName)
	}
	return address, stat.Size(), nil
}

func TestHttp(t *testing.T) {
	address, tootle, err := EggFlower("", "./Downloads/", 15<<20)
	if err != nil && err != io.EOF {
		t.Log(err.Error())
		return
	}
	objectName := "abb.zip"
	bucket := "hsj-1234"
	opts := minio.PutObjectOptions{NumThreads: 1, PartSize: 15 << 20}
	uploadID, err := Core.NewMultipartUpload(Background, bucket, objectName, opts)
	if err != nil {
		t.Log(err.Error())
		return
	}
	completePart := make([]minio.CompletePart, 0)
	size := int64(0)
	for i := range address {
		fileStrm, err := os.Open(address[i])
		if err != nil {
			t.Log(err, address[i])
			return
		}
		stat, err := fileStrm.Stat()
		if err != nil {
			t.Log(err.Error())
			return
		}
		PartNumber := i + 1
		objPart, err := Core.PutObjectPart(Background, bucket, objectName, uploadID, PartNumber,
			fileStrm, stat.Size(), "", "",
			opts.ServerSideEncryption)
		if err != nil {
			t.Log(err.Error())
			return
		}
		size += stat.Size()
		completePart = append(completePart, minio.CompletePart{
			ETag:       objPart.ETag,
			PartNumber: PartNumber,
		})
		fileStrm.Close()
	}
	sort.Sort(completedParts(completePart))
	etag, err := Core.CompleteMultipartUpload(Background, bucket, objectName, uploadID, completePart, opts)
	if err != nil {
		// return UploadInfo{}, err
		t.Log(err.Error())
		return
	}

	defer func() {
		if err != nil {
			Core.AbortMultipartUpload(Background, bucket, objectName, uploadID)
		}
	}()
	if tootle != size {
		err = errors.New("上传前后大小不一致")
		return
	}
	fmt.Println(etag)
}
