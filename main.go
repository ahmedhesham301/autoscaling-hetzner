package main

import (
	"autoscaling-hezner/controller"
	"autoscaling-hezner/hetzner"
	"autoscaling-hezner/storage"
	"context"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	storage.SetupClient(ctx)
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))
	// data, err := os.ReadFile("go.sum")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = storage.Client.PutObject(ctx, &s3.PutObjectInput{
	// 	Bucket: aws.String("test"),
	// 	Body:   bytes.NewReader(data),
	// 	Key:    aws.String("test/bla/go.sum"),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// result, _, err := hetzner.HClient.Server.Create(ctx, hcloud.ServerCreateOpts{
	// 	Name:       "test",
	// 	Image:      &hcloud.Image{Name: "ubuntu-24.04"},
	// 	ServerType: &hcloud.ServerType{Name: "cx23"},
	// 	Location:   &hcloud.Location{Name: "hel1"},
	// })
	// print(result.Server.Status)
	// h, err := hetzner.HClient.Location.
	server := gin.Default()

	server.GET("/locations", controller.GetLocations)
	server.GET("/images", controller.GetImages)
	server.GET("/types", controller.GetTypes)
	server.PUT("/templates", controller.CreateTemplate)
	server.PUT("/groups", controller.CreateGroup)

	server.Run()
}
