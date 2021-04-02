package producer

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"plane.watch/lib/tracker"
	"plane.watch/lib/tracker/beast"
	"plane.watch/lib/tracker/mode_s"
	"time"
)


func NewBeastFetcher(host, port string) tracker.Producer {
	p := NewProducer("Beast Fetcher for: " + net.JoinHostPort(host, port))

	source := newSource(p.String(), "beast://" + net.JoinHostPort(host, port))
	p.fetcher(host, port, func(conn net.Conn) error {
		scan := bufio.NewScanner(conn)
		scan.Split(ScanBeast)

		for scan.Scan() {
			msg := scan.Bytes()
			p.addFrame(beast.NewFrame(msg), source)
		}

		return scan.Err()
	})

	return p
}

func NewBeastListener(host, port string) tracker.Producer {
	p := NewProducer("Beast Listener for: " + net.JoinHostPort(host, port))

	go func() {
		// todo: Listen for incoming connection
	}()

	return p
}

func NewBeastFile(filePaths[] string, delay bool) tracker.Producer {
	p := NewProducer("Beast File")

	lastTimeStamp := time.Duration(0)
	p.readFiles(filePaths, func(reader io.Reader, fileName string) error {
		var count uint64
		scanner := bufio.NewScanner(reader)
		scanner.Split(ScanBeast)
		for scanner.Scan() {
			count++
			frame := beast.NewFrame(scanner.Bytes())
			if nil == frame {
				continue
			}
			if delay {
				currentTs := frame.(*mode_s.Frame).BeastTicksNs()
				if lastTimeStamp > 0 && lastTimeStamp < currentTs {
					time.Sleep(currentTs - lastTimeStamp)
				}
				lastTimeStamp = currentTs
			}

			p.AddEvent(tracker.NewFrameEvent(frame, newSource(p.String(), "file://" + fileName)))
		}
		p.addInfo("We processed %d lines", count)
		return scanner.Err()
	})

	return p
}

func ScanBeast(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// skip until we get our first 0x1A (message start)
	i := bytes.IndexByte(data, 0x1A)
	// byte 2 is our message type, so it tells us how long this message is
	msgLen := 0
	switch data[i+1] {
	case 0x31:
		// mode-ac 11 bytes (2+8)
		// 1(esc), 1(type), 6(mlat), 1(signal), 2(mode-ac)
		msgLen = 11
	case 0x32:
		// mode-s short 16 bytes
		// 1(esc), 1(type), 6(mlat), 1(signal), 7(mode-s short)
		msgLen = 16
	case 0x33:
		// mode-s long 23 bytes
		// 1(esc), 1(type), 6(mlat), 1(signal), 14(mode-s extended squitter)
		msgLen = 23
	case 0x34:
		// Config Settings and Stats
		// 1(esc), 1(type), 6(mlat), 1(unused), (1)DIP Config, (1)timestamp error ticks
		msgLen = 11
	default:
		// unknown? assume we got an out of sequence and skip
		return i+1, nil, nil
	}
	bufLen := len(data)
	//println("type", data[i+1], "input len", bufLen, "msg len",msgLen)
	if bufLen >= i + msgLen {
		// we have enough in our buffer
		// account for double escapes
		advance = msgLen
		max := msgLen * 2
		token = make([]byte, msgLen)
		dataIndex := i // start at the <esc>/0x1a
		tokenIndex := 0
		token[tokenIndex] = data[dataIndex]
		for dataIndex < i+advance-1 && dataIndex < i+max {
			dataIndex++ // first inc skips past the first 0x1a
			tokenIndex++
			// can we get to the next byte?
			if dataIndex+1 > bufLen {
				// run out of buffer, want more
				return 0, nil, nil
			}

			token[tokenIndex] = data[dataIndex]
			if data[dataIndex] != 0x1A {
				continue
			}
			if data[dataIndex+1] == 0x1A { // skip over the second <esc>
				advance++
				dataIndex++
			}
		}
		return advance, token, nil
	}
	// we want more data!
	return 0, nil, nil
}
