package ali

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Client struct {
	endpoint  string
	accessKey string
	objectKey string
	client    *oss.Client
}

func NewClient(endpoint, accessKey, objectKey string) (*Client, error) {
	innerCli, err := oss.New(endpoint, accessKey, objectKey)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		endpoint:  endpoint,
		accessKey: accessKey,
		objectKey: objectKey,
		client:    innerCli,
	}
	return cli, nil
}

func (cli *Client) Download(_ context.Context, bucketName, objectKey, localFilePath string) error {

	// 获取Bucket
	bucket, err := cli.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	partSize := int64(1024 * 1024 * 1024)
	// 下载文件到本地
	err = bucket.DownloadFile(objectKey, localFilePath, partSize, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) ListAndDownloadDir(_ context.Context, bucketName, prefix, localBasePath string) error {
	// 获取Bucket
	bucket, err := cli.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	partSize := int64(1024 * 1024 * 1024)
	// 使用递归方式列举目录下的所有对象
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Prefix(prefix), oss.Marker(marker))
		if err != nil {
			return err
		}

		for _, obj := range lsRes.Objects {
			relativeKey := strings.TrimPrefix(obj.Key, prefix)
			if relativeKey == "" || relativeKey == "/" {
				if !strings.HasSuffix(obj.Key, "/") {
					// 针对prefix是一个对象路径的情况
					fileName := filepath.Base(obj.Key)
					localFilePath := filepath.Join(localBasePath, fileName)
					if err = bucket.DownloadFile(obj.Key, localFilePath, partSize, nil); err != nil {
						return err
					}
				}
				continue
			}
			// 构建本地文件路径
			localFilePath := filepath.Join(localBasePath, relativeKey)
			if strings.HasSuffix(obj.Key, "/") {
				// 确保本地目录存在
				if err = os.MkdirAll(localFilePath, 0755); err != nil {
					return err
				}
				continue
			}

			// 下载单个文件
			if err = bucket.DownloadFile(obj.Key, localFilePath, partSize, nil); err != nil {
				return err
			}
		}

		if !lsRes.IsTruncated {
			break
		}
		marker = lsRes.NextMarker
	}

	return nil
}
