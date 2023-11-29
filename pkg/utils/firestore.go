package utils

import (
	"context"
	"log"

	firestore "cloud.google.com/go/firestore"
	"golang.ngrok.com/ngrok"
	"google.golang.org/api/option"
)

func InitializeFirestore() (firestore.Client, error) {
	// firestore
	opt := option.WithCredentialsFile("internal/config/urlsaver-f3a58-firebase-adminsdk-zitx7-6e40480f8d.json")
	client, err := firestore.NewClient(context.Background(), "urlsaver-f3a58", opt)
	return *client, err
}

func UpdateUrl(client firestore.Client, tun ngrok.Tunnel) (err error) {
	urlsCol := client.Collection("urls")
	socialMediaCommentDoc := urlsCol.Doc("mqsASg8cOVfa52ESQDLx")

	_, err = socialMediaCommentDoc.Update(context.TODO(), []firestore.Update{{Path: "url", Value: tun.URL()}})
	if err != nil {
		log.Fatal(err)
	}

	return
}
