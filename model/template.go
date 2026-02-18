package model

import (
	"autoscaling-hezner/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Template struct {
	Id          int64  `json:"id"`
	Image       string `json:"image"      binding:"required"`
	ServerType  string `json:"serverType" binding:"required"`
	CloudConfig string `json:"cloudConfig"`
}

func (t *Template) Save() error {
	jsonData, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}
	_, err = storage.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test"),
		Body:   bytes.NewReader(jsonData),
		Key:    aws.String(fmt.Sprintf("test/templates/%d.json", t.Id)),
	})
	return err
}

func GetTemplate(id int64) error {
	_, err := storage.Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String(fmt.Sprintf("test/templates/%d.json", id)),
	})
	fmt.Printf("/test/templates/%d.json", id)
	return err
}
