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


### Devices

This is a command line argument which gives a list of space separated device ids to mock

`node devicemock.js --devices b_pwSqQU0im_4dxRq5LDwByYaahQ6SKP8noz-8V0QZh5KYOG7KGs4Gp3N56VzEan-jsEdmc2RWrD6GY1d-HJv6TG3ClGykJmaHvFCFVCz5TZwkyqjpr7UGy2FlKU_7sJ BKTaSl3uyvfjq-vKG6EIMxDRk5a4Imw_ngimc0nNHDPswYwBo5E1CVIMxkOCRSovcp2cJnH-VMRV7SWOqz12MkL-DjS9ubDeesOwOWspEtENPoeivy1Zfujtv7BS8gl9`
