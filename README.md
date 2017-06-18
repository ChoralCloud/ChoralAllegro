# Allegro

This server is used to direct data sent by the device onto the kafka server. It also does some verification of the JSON to ensure that he data and format of the JSON object is valid.

## To run this server locally

1. Install golang
2. Run "go get github.com/julienschmidt/httprouter"
3. Run go build to generate executable
4. Run the server, it will be listening on port 3000 for post requests

## Testing

The test is only ensuring that the server is working correctly and responding to, and verifying post requests.

1. Make sure the allegro server is running
2. Run "go test"
