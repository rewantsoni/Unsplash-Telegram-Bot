package main

type data struct {
	Query string
	Random bool
}

type response struct {
	Urls url
	User user
}

type url struct {
	Full string
}

type user struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Links link
}

type link struct {
	Html string
}