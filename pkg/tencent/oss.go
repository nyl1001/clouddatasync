package tencent

import (
	"context"
	"net/http"
	"net/url"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type client struct {
	secretID  string
	secretKey string
}

func NewClient(secretID, secretKey string) (*client, error) {
	cli := &client{
		secretID:  secretID,
		secretKey: secretKey,
	}
	return cli, nil
}

func (cli *client) Download(bucketURL, objectKey, localFilePath string) error {
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cli.secretID,
			SecretKey: cli.secretKey,
			Transport: &http.Transport{},
		},
	})

	_, err := c.Object.GetToFile(context.Background(), objectKey, localFilePath, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) ListAndDownloadDir(bucketURL, prefix, localBasePath string) error {
	// 初始化COS客户端
	// 注意：腾讯云COS列举操作可能需要手动处理分页，此处仅作示例简化处理
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cli.secretID,
			SecretKey: cli.secretKey,
			Transport: &http.Transport{},
		},
	})

	opt := &cos.BucketGetOptions{
		Prefix:  prefix,
		MaxKeys: 1000, // 设置每次请求返回的最大对象数量，根据实际情况调整
	}

	// 此处需自行实现循环调用以处理分页，此处省略循环逻辑
	_, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		return err
	}
	// 实际上还需要处理返回的对象列表，逐个下载

	return nil
}
