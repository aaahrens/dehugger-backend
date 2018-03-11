package main

import "github.com/graph-gophers/graphql-go"

type Website struct {
	ID graphql.ID
	URL string
	Description string
	//base64
	Image string
	Metadata string
	Counter int32
	IsDown bool
}

type WebsiteInput struct {
	Url *string
	UUID *string
}

type CreateUserInput struct {
	Name *string
	PassWord *string
}

type User struct {
	AccessToken string
	Hugs []Website
}

