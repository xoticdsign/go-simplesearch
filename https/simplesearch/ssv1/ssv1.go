package ssv1

import (
	"github.com/gofiber/fiber/v2"
)

type SSv1 struct {
	Server Server
	Client Client
}

type Server struct {
	ServerImplementation *fiber.App
	Handlers             Handlerer
}

type Handlerer interface{}

type UnimplementedHandlers struct{}

type MakeSearchRequest struct {
	SearchFor string                   `json:"search_for"`
	Filters   MakeSearchRequestFilters `json:"filters"`
}

type MakeSearchRequestFilters struct {
	PriceBottom float64 `json:"price_bottom"`
	PriceTop    float64 `json:"price_top"`
}

type MakeSearchResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func (u *UnimplementedHandlers) MakeSearch(c *fiber.Ctx) error {
	return c.JSON(MakeSearchResponse{
		Result: "method unimplemented",
	})
}

type Client struct {
	ClientImplementation fiber.Client
}
