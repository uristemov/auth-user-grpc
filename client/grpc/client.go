package grpc

import (
	"context"
	"fmt"
	"github.com/uristemov/auth-user-grpc/models"
	"github.com/uristemov/auth-user-grpc/protobuf"
	"google.golang.org/grpc"
)

type Client struct {
	address  string
	conn     *grpc.ClientConn
	client   protobuf.UserClient
	dialOpts []grpc.DialOption
}

func NewClient(opts ...Option) (*Client, error) {
	cli := &Client{
		dialOpts: make([]grpc.DialOption, 0)}

	for _, opt := range opts {
		opt(cli)
	}

	err := cli.Connect()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.address, c.dialOpts...)
	if err != nil {
		return fmt.Errorf("error establishing gRPC connection to nats-streaming-reader: %s", err.Error())
	}

	c.conn = conn

	c.setupClient()

	return nil
}

func (c *Client) setupClient() {
	c.client = protobuf.NewUserClient(c.conn)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	resp, err := c.client.GetUserByEmail(ctx, &protobuf.GetUserByEmailRequest{Email: email})
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:        resp.Id,
		FirstName: resp.Firstname,
		LastName:  resp.Lastname,
		Email:     resp.Email,
		Password:  resp.Password,
	}, nil
}

func (c *Client) CreateUser(ctx context.Context, req *models.RegisterUser) (string, error) {

	grpcRequest := &protobuf.CreateUserRequest{
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	}

	resp, err := c.client.CreateUser(ctx, grpcRequest)
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
