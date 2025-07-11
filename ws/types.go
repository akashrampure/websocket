package ws

import "errors"

type MessageHandler func(clientID string, msg []byte)

var ErrClientNotConnected = errors.New("client is not connected")
