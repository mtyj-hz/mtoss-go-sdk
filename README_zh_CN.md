## 初始化麦田OSS Client
mtoss client需要以下4个参数来连接与Amazon S3兼容的对象存储。

| 参数  | 描述|
| :---         |     :---     |
| endpoint   | 对象存储服务的URL   |
| accessKeyID | Access key是唯一标识你的账户的用户ID。 |
| secretAccessKey | Secret key是你账户的密码。 |
| secure | true代表使用HTTPS |


```go
package main

import (
	"log"

	mtoss "github.com/mtyj-hz/mtoss-go-sdk"
	"github.com/mtyj-hz/mtoss-go-sdk/pkg/credentials"
)

func main() {
	endpoint := "oss.mty.wang"
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

### 获取文件cid
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