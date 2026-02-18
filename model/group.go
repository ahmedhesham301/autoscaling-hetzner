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

type Group struct {
	Id          int64    `json:"id"`
	TemplateId  int64    `json:"templateId"  binding:"required"`
	Zone        string   `json:"zone"        binding:"required"`
	Locations   []string `json:"locations"   binding:"required"`
	ServerTypes []string `json:"serverTypes" binding:"required"`
	MinSize     int      `json:"minSize"     binding:"required"`
	DesiredSize int      `json:"desiredSize" binding:"required"`
	MaxSize     int      `json:"maxSize"     binding:"required"`
}

func (g *Group) Save() error {
	jsonData, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		return err
	}
	_, err = storage.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test"),
		Body:   bytes.NewReader(jsonData),
		Key:    aws.String(fmt.Sprintf("test/groups/%d.json", g.Id)),
	})
	return err
}
