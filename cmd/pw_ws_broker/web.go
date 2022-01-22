package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"plane.watch/lib/export"
	"plane.watch/lib/tile_grid"
	"plane.watch/lib/ws_protocol"
	"sync"
	"time"
)

//go:embed test-web
var testWebDir embed.FS

type (
	PwWsBrokerWeb struct {
		Addr      string
		ServeTest bool

		serveMux   http.ServeMux
		httpServer http.Server

		clients   ClientList
		listening bool
	}
	WsClient struct {
		conn    *websocket.Conn
		outChan chan ws_protocol.WsResponse
		cmdChan chan WsCmd
	}
	WsCmd struct {
		action string
		what   string
	}
	ClientList struct {
		//clients     map[*WsClient]chan ws_protocol.WsResponse
		clients sync.Map
	}
)

func (bw *PwWsBrokerWeb) configureWeb() error {
	bw.clients = ClientList{}
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
	log.Debug().Str("New Connection", r.RemoteAddr).Msg("New /planes WS")
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{ws_protocol.WsProtocolPlanes},
		InsecureSkipVerify: false,
		OriginPatterns: []string{
			"*localhost*",
			"*plane.watch*",
		},
		//		OriginPatterns:       nil, // TODO maybe set this for plane.watch?
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if nil != err {
		log.Error().Err(err).Msg("Failed to setup websocket connection")
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Failed to setup websocket connection"))
		return
	}

	log.Debug().Str("protocol", conn.Subprotocol()).Msg("Speaking...")
	switch conn.Subprotocol() {
	case ws_protocol.WsProtocolPlanes:
		client := NewWsClient(conn)
		bw.clients.addClient(client)
		client.Handle(r.Context())
		bw.clients.removeClient(client)
	default:
		_ = conn.Close(websocket.StatusPolicyViolation, "Unknown Subprotocol")
		log.Debug().Str("proto", conn.Subprotocol()).Msg("Bad connection, could not speak protocol")
		return
	}
}

func (bw *PwWsBrokerWeb) HealthCheck() bool {
	log.Info().Bool("Web Listening", bw.listening).Msg("Health check")
	return bw.listening
}

func (bw *PwWsBrokerWeb) HealthCheckName() string {
	return "WS Broker Web"
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	client := WsClient{
		conn:    conn,
		cmdChan: make(chan WsCmd),
		outChan: make(chan ws_protocol.WsResponse),
	}
	return &client
}

func (c *WsClient) Handle(ctx context.Context) {
	err := c.planeProtocolHandler(ctx, c.conn)
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if nil != err {
		if -1 == websocket.CloseStatus(err) {
			log.Error().Err(err).Msg("Failure in protocol handler")
		}
		return
	}
}
func (c *WsClient) AddSub(tileName string) {
	log.Debug().Msg("Add Sub")
	c.cmdChan <- WsCmd{
		action: ws_protocol.RequestTypeSubscribe,
		what:   tileName,
	}
	log.Debug().Msg("Add Sub Done")
}
func (c *WsClient) UnSub(tileName string) {
	log.Debug().Msg("Unsub")
	c.cmdChan <- WsCmd{
		action: ws_protocol.RequestTypeUnsubscribe,
		what:   tileName,
	}
	log.Debug().Msg("Unsub done")
}
func (c *WsClient) SendTiles() {
	log.Debug().Msg("Unsub")
	c.cmdChan <- WsCmd{
		action: ws_protocol.RequestTypeSubscribeList,
		what:   "",
	}
	log.Debug().Msg("Unsub done")
}

func (c *WsClient) planeProtocolHandler(ctx context.Context, conn *websocket.Conn) error {
	// read from the connection for commands
	go func() {
		for {
			mt, frame, err := conn.Read(ctx)
			if nil != err {
				if errors.Is(err, io.EOF) || websocket.CloseStatus(err) >= 0 {
					return
				}
				log.Debug().Err(err).Int("Close Status", int(websocket.CloseStatus(err))).Msg("Error from reading")
				return
			}
			switch mt {
			case websocket.MessageText:
				log.Debug().Bytes("Client Msg", frame).Msg("From Client")
				rq := ws_protocol.WsRequest{}
				if err = json.Unmarshal(frame, &rq); nil != err {
					log.Warn().Err(err).Msg("Failed to understand message from client")
				}
				switch rq.Type {
				case ws_protocol.RequestTypeSubscribe:
					c.AddSub(rq.GridTile)
				case ws_protocol.RequestTypeSubscribeList:
					c.SendTiles()
					if nil != err {
						return
					}
				case ws_protocol.RequestTypeUnsubscribe:
					c.UnSub(rq.GridTile)
				default:
					_ = c.sendError(ctx, "Unknown request type")
				}

			case websocket.MessageBinary:
				_ = c.sendError(ctx, "Please speak text")
			}
		}
	}()

	// write a stream of location information
	subs := make(map[string]bool)

	for {
		var err error
		select {
		case cmdMsg := <-c.cmdChan:
			switch cmdMsg.action {
			case "exit":
				return nil
			case ws_protocol.RequestTypeSubscribe:
				subs[cmdMsg.what] = true
				err = c.sendAck(ctx, ws_protocol.ResponseTypeAckSub, cmdMsg.what)
			case ws_protocol.RequestTypeUnsubscribe:
				delete(subs, cmdMsg.what)
				err = c.sendAck(ctx, ws_protocol.ResponseTypeAckUnsub, cmdMsg.what)
			case ws_protocol.RequestTypeSubscribeList:
				tiles := make([]string, 0, len(subs))
				for k, v := range subs {
					if v {
						tiles = append(tiles, k)
					}
				}
				err = c.sendPlaneMessage(ctx, &ws_protocol.WsResponse{
					Type:  ws_protocol.ResponseTypeSubTiles,
					Tiles: tiles,
				})

			}
		case planeMsg := <-c.outChan:
			err = c.sendPlaneMessage(ctx, &planeMsg)
		}

		if nil != err {
			return err
		}
	}
}

func (c *WsClient) sendAck(ctx context.Context, ackType, tile string) error {
	rs := ws_protocol.WsResponse{
		Type:  ackType,
		Tiles: []string{tile},
	}
	return c.sendPlaneMessage(ctx, &rs)
}

func (c *WsClient) sendError(ctx context.Context, msg string) error {
	rs := ws_protocol.WsResponse{
		Type:    ws_protocol.ResponseTypeError,
		Message: msg,
	}
	return c.sendPlaneMessage(ctx, &rs)
}

func (c *WsClient) sendPlaneMessage(ctx context.Context, planeMsg *ws_protocol.WsResponse) error {
	buf, err := json.MarshalIndent(planeMsg, "", "  ")
	if nil != err {
		log.Debug().Err(err).Str("type", planeMsg.Type).Msg("Failed to marshal plane msg to send to client")
		return err
	}
	if err = c.writeTimeout(ctx, 3*time.Second, buf); nil != err {
		log.Debug().
			Err(err).
			Str("type", planeMsg.Type).
			Msgf("Failed to send message to client. %+v", err)
		return err
	}
	return nil
}

func (c *WsClient) writeTimeout(ctx context.Context, timeout time.Duration, msg []byte) error {
	ctxW, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.conn.Write(ctxW, websocket.MessageText, msg)
}

func (cl *ClientList) addClient(c *WsClient) {
	log.Debug().Msg("Add Client")
	cl.clients.Store(c, true)
	prometheusNumClients.Inc()
	log.Debug().Msg("Add Client Done")
}

func (cl *ClientList) removeClient(c *WsClient) {
	log.Debug().Msg("Remove Client")
	close(c.outChan)
	cl.clients.Delete(c)
	prometheusNumClients.Dec()
	log.Debug().Msg("Remove Client Done")
}

// SendLocationUpdate sends an update to each listening client
// todo: make this threaded?
func (cl *ClientList) SendLocationUpdate(highLow, tile string, loc *export.EnrichedPlaneLocation) {
	cl.clients.Range(func(key, value interface{}) bool {
		defer func() {
			if r := recover(); nil != r {
				log.Error().Msgf("Panic: %v", r)
			}
		}()
		client := key.(*WsClient)
		client.outChan <- ws_protocol.WsResponse{
			Type:     ws_protocol.ResponseTypePlaneLocation,
			Location: loc,
		}
		return true
	})
}
