[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![LICENSE](https://img.shields.io/hexpm/l/pulsar.svg)](https://github.com/andvikram/goreal/blob/master/LICENSE)

# GoReal

An application to simplify PubSub using [Websocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API) and [Apache Pulsar](http://pulsar.apache.org)

## Getting Started

Fire up Pulsar with Docker:

```
docker run -it \
  -p 6650:6650 \
  -p 8080:8080 \
  --mount source=pulsardata,target=/pulsar/data \
  --mount source=pulsarconf,target=/pulsar/conf \
  apachepulsar/pulsar:2.5.1 \
  bin/pulsar standalone
```

Build the application:

```
go build
```

Start the application:

```
./goreal run -c <full_path_to_config_directory>
```

### Config file

Please make sure to create a config YAML file with the name `goreal-config.yml`

The default directory is $HOME.

Sample config file content:

```yml
## OIRIGIN should hold the address of the server
##  rendering the client that will connect with GoReal
ORIGIN: http://localhost:4000
MESSAGE_BUS: pulsar
MESSAGE_BUS_URL: pulsar://localhost:6650

APPLICATION:
  NAME: goreal
  SCHEME: http
  HOST: localhost
  PORT: :4100
  LOG_LEVEL: debug
  LOG_FILE_PATH: /tmp/goreal.log

# When using SSL/TLS
SECURE:
  CERT_PATH: /home/user/cert.pem
  KEY_PATH: /home/user/key.pem
```

### Client side

On the frontend, connect with GoReal over Websocket when you want to subscribe to a topic.

Sample subscribe function:

```js
const ws = new WebSocket("ws://<GoReal_server_address>/subscribe");
// ws is provided to the function 
// so that it can be used to close on exit events
function subscribe(ws, topicID, peer) {
  // This message format is required by GoReal
  // Subscription is given a name using peer
  const msg = { topicID:  topicID, peer: peer };
  ws.onopen = function () {
    console.log("WebSocket connection established with server");
    ws.send(JSON.stringify(msg));
  }
  ws.onclose = function () {
    console.log("WebSocket connection closed by server");
  }
  ws.onmessage = function (event) {
    let mData = JSON.parse(event.data);
    console.log("Received message:", mData);
  }
}
```

### Server side

Use the `PublishToTopicRoute` endpoint to publish your message to the provided topic.

### API

```js
// SubscribeTopicRoute defines Websocket path to subscribe and listen to a topic
SubscribeTopicRoute = "/subscribe"

// PublishToTopicRoute defines HTTP path for publishing a message in a topic
PublishToTopicRoute = "/topics/:topic_id/messages/new"
  
  Required POST data:
    { message: JSON.stringify(message) }
```
