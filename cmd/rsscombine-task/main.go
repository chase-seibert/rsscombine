package main

import (
  "github.com/rsscombine"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/spf13/viper"
  "log"
  "strings"
)

func main() {
  rsscombine.LoadConfig()
  bucket := viper.GetString("s3_bucket")
  filename := viper.GetString("s3_filename")
  combinedFeed := rsscombine.GetAtomFeed()
  atom, _ := combinedFeed.ToAtom()
  log.Printf("Rendered RSS with %v items", len(combinedFeed.Items))
  sess, err := session.NewSession(&aws.Config{
    //Region: aws.String("us-west-2")},
  })
  uploader := s3manager.NewUploader(sess)
  _, err = uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(bucket),
    Key: aws.String(filename),
    Body: strings.NewReader(atom),
  })
  if err != nil {
      log.Fatal("Unable to upload %q to %q, %v", filename, bucket, err)
  }
  log.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}
