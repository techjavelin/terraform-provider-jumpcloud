package api

import (
	"context"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

const API_ACCEPT_TYPE = "application/json"
const API_CONTENT_TYPE = "application/json"

type JumpCloudClientApiV2 struct {
	Client *jcapiv2.APIClient
	Auth   context.Context
}

type JumpCloudClientApiV1 struct {
	Client *jcapiv1.APIClient
	Auth   context.Context
}

func init() {
}

func (api JumpCloudClientApiV2) V2(client *jcapiv2.APIClient, auth context.Context) JumpCloudClientApiV2 {
	api.Client = client
	api.Auth = auth
	return api
}

func (api JumpCloudClientApiV1) V1(client *jcapiv1.APIClient, auth context.Context) JumpCloudClientApiV1 {
	api.Client = client
	api.Auth = auth
	return api
}
