package api

import (
	"gopkg.in/resty.v1"
	"studentbookef/config"
)

/**
*this is a class that set up the Api address to consume "http://localhost:9099/sts/"
* Port: 9099 domain: bookstore
*It also set the type of messaging protocol in our case we will be using JSON
**/
const BASE_URL string = "http://localhost:8089/sts/" //connection port

func Rest() *resty.Request {
	return resty.R().SetAuthToken("").
		SetHeader("Accept", "application/json").
		SetHeader("email", "email").
		SetHeader("site", "site").
		SetHeader("Content-Type", "application/json")
}

var JSON = config.ConfigWithCustomTimeFormat
