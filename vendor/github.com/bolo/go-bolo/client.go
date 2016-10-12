package bolo

import "github.com/pebbe/zmq4"

func Connect(url string) (chan PDU, chan error, error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, nil, err
	}
	socket, err := context.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, nil, err
	}
	err = socket.Connect(url)
	if err != nil {
		return nil, nil, err
	}
	err = socket.SetSubscribe("")
	if err != nil {
		return nil, nil, err
	}

	pduChan := make(chan PDU, 5)
	errChan := make(chan error, 5)

	go func() {
		for {
			// FIXME: reconnect on failures? bail out at some point?
			msg, err := socket.RecvMessage(0)
			if err != nil {
				errChan <- err
			}

			// parse message into PDU interface
			pdu, err := ParsePDU(msg)
			if err != nil {
				errChan <- err
			} else {
				pduChan <- pdu
			}
		}
		close(pduChan)
		close(errChan)
		socket.Close()
	}()

	return pduChan, errChan, nil
}
