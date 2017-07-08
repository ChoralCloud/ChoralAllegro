var request = require('request')

const interval = 500
const N = (0|process.env.N) || 1;

for(var i = 0; i < N; ++i)
  setInterval(sendRequest("" + i), interval)

function sendRequest(dev_id) {
  // will be ignored
  return function(){
    var now = new Date()

    var data = {
      "device_id": dev_id,
      "user_secret":  "zyx321",
      "device_data": {
        "sensor_data": Math.random() * interval,
      },
      "device_timestamp": now.getTime()
    }

    var options = {
      method: 'POST',
      uri: process.env.URL || "http://localhost:8081",
      json: true,
      body: data
    }

    request(options, function (err, res, body) {
      if (err) {
        console.error('error posting json: ', err, body)
        return
      }

      if(body.error){
        console.log(body.error)
      }
    });
  }
}
