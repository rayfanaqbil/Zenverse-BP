package config

import "os"

var ApiWaButton string = os.Getenv("URLAPIWABUTTON")

var GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail string
