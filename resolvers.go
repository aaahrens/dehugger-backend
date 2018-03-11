package main

import (
	"github.com/graph-gophers/graphql-go"
)

type Resolver struct {
}

//top level func
func (r *Resolver) Websites(args struct{ First int32 }) []*websiteResolver {

	var l []*websiteResolver
	list := HugsCacher.GetHugs(0, 1)
	for _, website := range list {
		//spew.Dump(website)
		l = append(l, &websiteResolver{website})

	}
	return l
}

func (r *Resolver) FetchFollows(args struct{ Token *string }) []*websiteResolver {
	var l []*websiteResolver
	list := HugsCacher.GetHugs(0, 1)
	for _, website := range list {
		//spew.Dump(website)
		l = append(l, &websiteResolver{website})

	}
	return l
}

func (r *Resolver) FetchUser(args struct {
	Token *string
}) *userResolver {
	return &userResolver{
		user: &User{
			"Penis",
			[]Website{},
		}}
}

//site reporting
func (r *Resolver) ReportSite(args struct{ Input *WebsiteInput }) *bool {
	isError := true
	writeChannel <- Website{
		URL: *args.Input.Url,
	}
	return &isError
}

//user creation
func (r *Resolver) CreateUser(args struct{ User *CreateUserInput }) *userResolver {

	return &userResolver{
		user: &User{
			"Penis",
			[]Website{},
		}}

}

func (r *Resolver) LoginUser(args struct {
	Name     *string
	Password *string
}) *userResolver {
	return &userResolver{
		user: &User{
			"Penis",
			[]Website{},
		}}
}

//user resolver

//for users
type userResolver struct {
	user *User
}

func (user userResolver) Token() *string {
	return &user.user.AccessToken
}

func (user userResolver) Hugs() []*websiteResolver {
	l := make([]*websiteResolver, len(user.user.Hugs))
	for index, item := range user.user.Hugs {
		l[index] = &websiteResolver{&item}
	}
	return l
}

//for websites
type websiteResolver struct {
	w *Website
}

func (u *websiteResolver) ID() (*graphql.ID, error) {
	return &u.w.ID, nil
}

func (u *websiteResolver) Url() (*string, error) {

	return &u.w.URL, nil
}

func (u *websiteResolver) Description() (*string, error) {
	return &u.w.Description, nil
}

func (u *websiteResolver) Image() (*string, error) {
	return &u.w.Image, nil
}

func (u *websiteResolver) MetaData() (*string, error) {
	return &u.w.Metadata, nil
}

func (u *websiteResolver) Counter() (*int32, error) {
	return &u.w.Counter, nil
}
