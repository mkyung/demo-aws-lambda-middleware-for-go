package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaHandler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type MiddlewareHandler func(events.APIGatewayProxyRequest, events.APIGatewayProxyResponse) (events.APIGatewayProxyRequest, events.APIGatewayProxyResponse, error)

type Chain struct {
	head       *Middleware
	middleware *Middleware
}

type Middleware struct {
	current MiddlewareHandler
	next    *Middleware
}

func (m *Middleware) run(request events.APIGatewayProxyRequest, response events.APIGatewayProxyResponse) (events.APIGatewayProxyResponse, error) {
	req, res, err := m.current(request, response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if m.next == nil {
		return res, nil
	}
	return m.next.run(req, res)
}

func (c *Chain) Add(m MiddlewareHandler) {
	midd := &Middleware{}
	midd.current = m
	c.middleware.next = midd
	c.middleware = midd
}

func (c *Chain) Build() LambdaHandler {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res := events.APIGatewayProxyResponse{}
		res2, err := c.head.next.run(req, res)
		return res2, err
	}
}

func NewChain() Chain {
	midd := &Middleware{}
	return Chain{
		head:       midd,
		middleware: midd,
	}
}

func main() {
	chain := NewChain()
	chain.Add(Middleware1)
	chain.Add(Middleware2)
	c := chain.Build()
	lambda.Start(c)
}

func Middleware1(req events.APIGatewayProxyRequest, res events.APIGatewayProxyResponse) (events.APIGatewayProxyRequest, events.APIGatewayProxyResponse, error) {
	res.Body = res.Body + " middle1 "
	res.StatusCode = 201
	return req, res, nil
}

func Middleware2(req events.APIGatewayProxyRequest, res events.APIGatewayProxyResponse) (events.APIGatewayProxyRequest, events.APIGatewayProxyResponse, error) {
	res.Body = res.Body + " middle2 "
	res.StatusCode = 200
	return req, res, nil
}
