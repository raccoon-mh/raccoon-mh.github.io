package strapi

import (
	"github.com/gobuffalo/envy"
)

type StrapiResponse struct {
	Data  *Data  `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

var (
	STRAPI_URL         string
	STRAPI_ADMIN_TOKEN string
)

func init() {
	STRAPI_URL = envy.Get("STRAPI_URL", "http://localhost:1337")
	STRAPI_ADMIN_TOKEN = envy.Get("STRAPI_ADMIN_TOKEN", "")
}

type Data struct {
	ID         int         `json:"id"`
	Attributes interface{} `json:"attributes"`
}

type Error struct {
	Status  int          `json:"status"`
	Name    string       `json:"name"`
	Message string       `json:"message"`
	Details ErrorDetails `json:"details"`
}

type ErrorDetails struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Path    []string `json:"path"`
	Message string   `json:"message"`
	Name    string   `json:"name"`
}