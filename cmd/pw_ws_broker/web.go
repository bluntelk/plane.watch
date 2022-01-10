package main

import (
	"context"
	"embed"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"nhooyr.io/websocket"
	"plane.watch/lib/export"
	"plane.watch/lib/tile_grid"
	"time"
)

//go:embed test-web
var testWebDir embed.FS

const (
	WsProtocolPlanes         = "planes"
	RequestTypeSubscribe     = "sub"
	RequestTypeSubscribeList = "sub-list"
	RequestTypeUnsubscribe   = "unsub"

	ResponseTypeError         = "error"
	ResponseTypePlaneLocation = "plane-location"
)

type (
	PwWsBrokerWeb struct {
		Addr      string
		ServeTest bool

		serveMux   http.ServeMux
		httpServer http.Server
	}
	WsRequest struct {
		Type     string
		GridTile string
	}
	WsResponse struct {
		Type     string
		Message  string
		location *export.PlaneLocation
	}
)

func (bw *PwWsBrokerWeb) configureWeb() error {
	bw.serveMux.HandleFunc("/", bw.indexPage)
	bw.serveMux.HandleFunc("/grid", bw.jsonGrid)
	bw.serveMux.HandleFunc("/planes", bw.servePlanes)

	if bw.ServeTest {
		bw.serveMux.Handle(
			"/test-web/",
			bw.logRequest(
				http.FileServer(http.FS(testWebDir)),
			),
		)
	}

	bw.httpServer = http.Server{
		Addr:         bw.Addr,
		Handler:      &bw.serveMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	return nil
}

func (bw *PwWsBrokerWeb) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("Remote", r.RemoteAddr).Str("Request", r.RequestURI).Msg("Web RQ")
		handler.ServeHTTP(w, r)
	})
}

func (bw *PwWsBrokerWeb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw.serveMux.ServeHTTP(w, r)
}

func (bw *PwWsBrokerWeb) indexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("Plane.Watch Websocket Broker"))
}

func (bw *PwWsBrokerWeb) jsonGrid(w http.ResponseWriter, r *http.Request) {
	grid := tile_grid.GetGrid()
	buf, err := json.MarshalIndent(grid, "", "  ")
	if nil != err {
		w.WriteHeader(500)
	}
	_, _ = w.Write(buf)
}

func (bw *PwWsBrokerWeb) servePlanes(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{WsProtocolPlanes},
		InsecureSkipVerify: false,
		//		OriginPatterns:       nil, // maybe set this for plane.watch?
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if nil != err {
		log.Error().Err(err).Msg("Failed to setup websocket connection")
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Failed to setup websocket connection"))
		return
	}

	switch conn.Subprotocol() {
	case WsProtocolPlanes:
		for {
			err = bw.planeProtocolHandler(r.Context(), conn)
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			if nil != err {
				log.Error().Err(err).Msg("Failure in protocol handler")
				return
			}
		}
	default:
		_ = conn.Close(websocket.StatusPolicyViolation, "Unknown Subprotocol")
		log.Debug().Str("proto", conn.Subprotocol()).Msg("Bad connection, could not speak protocol")
		return
	}
}

func (bw *PwWsBrokerWeb) planeProtocolHandler(ctx context.Context, conn *websocket.Conn) error {
	outMessages := make(chan WsResponse)

	// read from the connection for commands
	go func() {
		for {
			mt, frame, err := conn.Read(ctx)
			if nil != err {
				log.Debug().Err(err).Msg("Error from reading")
				return
			}
			switch mt {
			case websocket.MessageText:
				log.Debug().Bytes("Client Msg", frame).Msg("From Client")
				rq := WsRequest{}
				if err = json.Unmarshal(frame, &rq); nil != err {
					log.Warn().Err(err).Msg("Failed to understand message from client")
				}
				switch rq.Type {
				case RequestTypeSubscribe:
				case RequestTypeSubscribeList:
				case RequestTypeUnsubscribe:
				default:
					_ = bw.sendError(ctx, conn, "Unknown request type")
				}

			case websocket.MessageBinary:
				_ = bw.sendError(ctx, conn, "Please speak text")
			}
		}
	}()

	// write a stream of location information

	for planeMsg := range outMessages {
		_ = bw.sendPlaneMessage(ctx, conn, &planeMsg)
	}

	return nil
}

func (bw *PwWsBrokerWeb) sendError(ctx context.Context, conn *websocket.Conn, msg string) error {
	rs := WsResponse{
		Type:    ResponseTypeError,
		Message: msg,
	}
	return bw.sendPlaneMessage(ctx, conn, &rs)
}

func (bw *PwWsBrokerWeb) sendPlaneMessage(ctx context.Context, conn *websocket.Conn, planeMsg *WsResponse) error {
	buf, err := json.MarshalIndent(planeMsg, "", "  ")
	if nil != err {
		log.Debug().Err(err).Str("type", planeMsg.Type).Msg("Failed to marshal plane msg to send to client")
		return err
	}
	if err = bw.writeTimeout(ctx, conn, 3*time.Second, buf); nil != err {
		log.Debug().Err(err).Str("type", planeMsg.Type).Msg("Failed to send message to client")
		return err
	}
	return nil
}

func (bw *PwWsBrokerWeb) writeTimeout(ctx context.Context, conn *websocket.Conn, timeout time.Duration, msg []byte) error {
	ctxW, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return conn.Write(ctxW, websocket.MessageText, msg)
}
