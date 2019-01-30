# Release Notes

## 0.1.2

* fixded etag for small files (smalles then 5 MB)
  * etag calculation is different for small files

## 0.1.1

* simplified hash loop to improve performance
* fixed missing lf in client

## 0.1.0

* fixed mod issue with client
* implemented s3 etag fetch in package
* implemented etag comparission in client

## 0.0.1

* only generates the fs etag like s3 multipart upload
