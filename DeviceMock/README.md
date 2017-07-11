# Device Mock

This simple javascript file makes it easy to send mock request to the development cluster

By default it mocks a single device and sends every half second

## Install

`npm install`

## Output

This module will print any errors that are sent back from the server

## Environment Variables

There are Enviromnent variables that you can set to make this server behave diferently

### N: Default 1

setting N will make there be more than one device that updates

```
N=5 node devicemock.js
```

this will run 5 mock devices

### URL: Default http://localhost:8081

This will change the expected location of choral allegro

URL="http://choralcluster.csc.uvic.ca:8081" node devicemock.js

