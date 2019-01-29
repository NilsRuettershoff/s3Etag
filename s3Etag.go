package s3Etag

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// CalculateLocalETag calculates the local etag
func CalculateLocalETag(path string, chunksize int) (etag string, err error) {
	_, ok := os.Stat(path)
	if os.IsNotExist(ok) {
		return "", errors.Wrap(err, fmt.Sprintf("inputfile %s does not exists", path))
	}
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "unable to open file")
	}
	chunkbytes := int64(chunksize * 1024 * 1024)
	// simple sequential way
	buffer := make([]byte, chunkbytes, chunkbytes)
	bufr := bytes.NewReader(buffer)
	reader := bufio.NewReader(f)
	gh := md5.New()
	counter := 0
	var n int
	run := true
	// generate hash of chunk and write result to global hash
	for run {
		h := md5.New()
		bufr.Reset(buffer)
		n, err = io.ReadFull(reader, buffer)
		if err == io.ErrUnexpectedEOF {
			run = false
		}
		if err == io.EOF {
			run = false
			continue
		}
		_, err = io.CopyN(h, bufr, int64(n))
		gh.Write(h.Sum(nil))
		counter++
	}
	etag = hex.EncodeToString(gh.Sum(nil))
	etag += fmt.Sprintf("-%d", counter)
	return etag, err
}

// FetchS3Etag fetches S3 etag of an object via s3client from bucket
func FetchS3Etag(s3client *s3.S3, bucket string, key string) (etag string, err error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key)}
	result, err := s3client.HeadObject(input)
	if err != nil {
		return "", err
	}
	etag = *result.ETag
	// api adds double quotes at the beginning and end
	etag = strings.Replace(etag, "\"", "", -1)
	return etag, err
}
