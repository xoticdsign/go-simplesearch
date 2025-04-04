package ssv1

import "github.com/gofiber/fiber/v2"

type SSv1 struct {
	Server Server
	Client Client
}

type Server struct {
	ServerImplementation *fiber.App
	Service              Servicer
}

type Servicer interface{}

type unimplementedService struct{}

type Client struct {
	ClientImplementation fiber.Client
}
