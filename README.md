# Allegro

This server is used to direct data sent by the device onto the kafka server. It also does some verification of the JSON to ensure that he data and format of the JSON object is valid.

## Verification

We need to ensure that data retrieved by the server is in the correct format before we attempt to do batch computations on it. Since we want people using this platform to follow our protocol, we expect a certain format for the data. The following expectations must be met:

- Data must be in JSON format.
- Timestamp must be within the last hour (implemented)
- Data must be correctly formatted (what is the expected format?)
- ID must be correctly formatted (it should be unique, we should probably assign these to devices when users add new ones or something)

## To run this server locally

1. Install golang
2. Run "go get github.com/julienschmidt/httprouter"
3. Run "go get github.com/Shopify/sarama"
4. Run "go build" to generate executable
5. Start the server, it will be listening on port 3000 for post requests

We might want to consider using godep for dependency management, especially if we add any more.

## Docker

There are two ways to run this project using docker

### With Choral Storm

Run choral storm as described in the readme, then run allegro by running

`docker-compose up`

### Locally

There is also a docker folder so you can run this in isolation from ChoralStorm. It takes considerably less memory to do this, and so it is useful for development purposes. To do this:

1. Ensure you aren't already running the image from ChoralStorm
2. cd docker
3. run "docker-compose up"

## Testing

Currently the test is not actually testing anything, but it does send a mock JSON object over to the allegro server, and we can see kafka consuming it. To run this:

1. Make sure the allegro and kafka servers are running
2. Run "go test"

### TODO

- Use httptest to monitor the server requests. To do this, we might have to remove the httprouter library.
- Handle requests concurrently, we can do this with go func()
