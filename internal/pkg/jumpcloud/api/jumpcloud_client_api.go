package api

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

const API_ACCEPT_TYPE = "application/json"
const API_CONTENT_TYPE = "application/json"

type JumpCloudClientApi struct {
	Client *jcapiv2.APIClient
	Auth   context.Context
}

func init() {
	fmt.Println("api package initialized")
}

func (api JumpCloudClientApi) New(client *jcapiv2.APIClient, auth context.Context) JumpCloudClientApi {
	api.Client = client
	api.Auth = auth
	return api
}
