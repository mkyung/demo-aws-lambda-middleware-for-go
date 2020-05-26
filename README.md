# Writing re-usable middleware for AWS Lambda in Golang

### TL;DR

This example shows how we can write middleware for the AWS Lambda and use it like the Go library `Negroni`.
The middleware is run in sequence and the `APIGatewayProxyResponse` is shared among the middleware to allow overwrite and update.

```go
func main() {
	chain := NewChain()
	chain.Add(Middleware1)
	chain.Add(Middleware2)
	c := chain.Build()
	lambda.Start(c)
}
```

Some examples for `Middleware1` and `Middleware2`:

```go
func Middleware1 (req events.APIGatewayProxyRequest, res events.APIGatewayProxyResponse) (events.APIGatewayProxyRequest, events.APIGatewayProxyResponse, error) {
	res.Body = res.Body + " middle1 "
	res.StatusCode = 201
	return req, res, nil
}

func Middleware2 (req events.APIGatewayProxyRequest, res events.APIGatewayProxyResponse) (events.APIGatewayProxyRequest, events.APIGatewayProxyResponse, error) {
	res.Body = res.Body + " middle2 "
	res.StatusCode = 200
	return req, res, nil
}
```

Results:

```json
{
  "statusCode": 200,
  "headers": null,
  "multiValueHeaders": null,
  "body": " middle1  middle2 "
}
```