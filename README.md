## Initialize mtOss Client
mtoss client requires the following four parameters specified to connect to an Amazon S3 compatible object storage.

| Parameter  | Description|
| :---         |     :---     |
| endpoint   | URL to object storage service.   |
| _minio.Options_ | All the options such as credentials, custom transport etc. |

```go
package main

import (
	"log"

	mtoss "github.com/mtyj-hz/mtoss-go-sdk"
	"github.com/mtyj-hz/mtoss-go-sdk/pkg/credentials"
)

func main() {
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	// 初使化 minio client对象。
	minioClient, err := mtoss.New(endpoint, &mtoss.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient初使化成功
}
```
### getCid
```go
 package main

import (
	"log"
	"fmt"
	mtoss "github.com/mtyj-hz/mtoss-go-sdk"
	"github.com/mtyj-hz/mtoss-go-sdk/pkg/credentials"
	"context"
)

func main() {
	endpoint := "oss.mty.wang"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	// 初使化 minio client对象。
	mtossClient, err := mtoss.New(endpoint, &mtoss.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", mtossClient) // minioClient初使化成功

	ctx := context.Background()
	opts := mtoss.GetObjectOptions{}
	bucketName := "y1211"
	objectName := "app.zip"
	objectInfo, err := mtossClient.StatObject(ctx, bucketName,objectName , opts)
	if err != nil {
		return
	}
	fmt.Println(objectInfo.Cid)
}
```