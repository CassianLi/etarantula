package ossutil

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"strings"
)

type AliOss struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	client          *oss.Client
}

// NewAliOss 创建一个AliOss实例
func NewAliOss(endpoint, accessKeyId, accessKeySecret string) (ali *AliOss, err error) {
	ali = &AliOss{
		Endpoint:        endpoint,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	ali.client, err = oss.New(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return ali, nil
}

// UploadByte 上传byte[]到指定的bucket，并指定文件路径，如：exampledir/example.txt
func (ali *AliOss) UploadByte(bucketName, objectKey string, data []byte) (err error) {
	bucket, err := ali.client.Bucket(bucketName)
	if err != nil {
		log.Println("get bucket failed, err: ", err)
		return err
	}

	// 将Byte数组上传至 objectKey(example/object.txt)
	err = bucket.PutObject(objectKey, bytes.NewReader(data))
	if err != nil {
		log.Println("put object failed, err: ", err)
	}
	return err
}

// UploadString 上传字符串到指定bucket，并指定文件路径，如：exampledir/example.txt
func (ali *AliOss) UploadString(bucketName, objectKey, data string) (err error) {
	bucket, err := ali.client.Bucket(bucketName)
	if err != nil {
		log.Println("get bucket failed, err: ", err)
		return err
	}

	// 指定Object存储类型为低频访问。
	storageType := oss.ObjectStorageClass(oss.StorageIA)

	// 指定Object访问权限为私有。
	objectAcl := oss.ObjectACL(oss.ACLPrivate)

	// 将字符串上传至 objectKey(example/object.txt)
	err = bucket.PutObject(objectKey, strings.NewReader(data), storageType, objectAcl)
	if err != nil {
		log.Println("put object failed, err: ", err)
	}
	return err
}

// UploadFile 上传本地文件到指定bucket，并指定文件路径，如：exampledir/example.txt
func (ali *AliOss) UploadFile(bucketName, objectKey, localFile string) (err error) {
	bucket, err := ali.client.Bucket(bucketName)
	if err != nil {
		log.Println("get bucket failed, err: ", err)
		return err
	}

	err = bucket.PutObjectFromFile(objectKey, localFile)
	if err != nil {
		log.Println("put object failed, err: ", err)
	}

	return err
}
