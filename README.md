# chan-TCP-nels

Implementing bidirectional Go channels over TCP
* Send/Receive pair of channels that can exchange data between themselves over TCP
* This allows the channels to wait and send information over the internet no matter where they are instantiated
* Current implementation exhange messages between local and remote server and close connection, this can be expanded

## Usage
(This order can be reversed, it does not matter)
* In the root directory, run go run main.go -> channel starts trying to dial :4000
* In the /remote directory, run go run main.go -> previous server (at :3000) connects here and sends a message on a receive channel here
* Both channels send and receive messages to and from the other one

## TODO
* Expand functionality to more than just send and receive a message and close
* Better interfacing for usage
