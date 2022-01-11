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
	"sync"
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
	ResponseTypeAckSub        = "ack-sub"
	ResponseTypeAckUnsub      = "ack-unsub"
	ResponseTypeSubTiles      = "sub-list"
	ResponseTypePlaneLocation = "plane-location"
)

type (
	PwWsBrokerWeb struct {
		Addr      string
		ServeTest bool

		serveMux   http.ServeMux
		httpServer http.Server

		clients ClientList
	}
	WsClient struct {
		conn *websocket.Conn

		subLock sync.RWMutex
		subs    map[string]bool

		outChan chan WsResponse
	}

	WsRequest struct {
		Type     string `json:"type"`
		GridTile string `json:"gridTile"`
	}
	WsResponse struct {
		Type     string                        `json:"type"`
		Message  string                        `json:"message,omitempty"`
		Tiles    []string                      `json:"tiles,omitempty"`
		Location *export.EnrichedPlaneLocation `json:"location,omitempty"`
	}

	ClientList struct {
		clients     map[*WsClient]chan WsResponse
		clientsLock sync.RWMutex
	}
)

func (bw *PwWsBrokerWeb) configureWeb() error {
	bw.clients = ClientList{
		clients: make(map[*WsClient]chan WsResponse),
	}
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
		client := NewWsClient(conn)
		bw.clients.addClient(client, client.outChan)
		client.Handle(r.Context())
		bw.clients.removeClient(client)
	default:
		_ = conn.Close(websocket.StatusPolicyViolation, "Unknown Subprotocol")
		log.Debug().Str("proto", conn.Subprotocol()).Msg("Bad connection, could not speak protocol")
		return
	}
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	client := WsClient{
		conn:    conn,
		subs:    map[string]bool{},
		outChan: make(chan WsResponse),
	}
	return &client
}

func (c *WsClient) Handle(ctx context.Context) {
	err := c.planeProtocolHandler(ctx, c.conn)
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if nil != err {
		log.Error().Err(err).Msg("Failure in protocol handler")
		return
	}
}
func (c *WsClient) AddSub(tileName string) {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	c.subs[tileName] = true
}
func (c *WsClient) UnSub(tileName string) {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	c.subs[tileName] = false
}
func (c *WsClient) SubTiles() []string {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	tiles := make([]string, 0, len(c.subs))
	for k, v := range c.subs {
		if v {
			tiles = append(tiles, k)
		}
	}
	return tiles
}
func (c *WsClient) HasSub(tileName string) bool {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	val, ok := c.subs[tileName]
	return val && ok
}

func (c *WsClient) planeProtocolHandler(ctx context.Context, conn *websocket.Conn) error {
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
					c.AddSub(rq.GridTile)
					_ = c.sendAck(ctx, ResponseTypeAckSub, rq.GridTile)
				case RequestTypeSubscribeList:
					_ = c.sendPlaneMessage(ctx, &WsResponse{
						Type:  ResponseTypeSubTiles,
						Tiles: c.SubTiles(),
					})
				case RequestTypeUnsubscribe:
					c.UnSub(rq.GridTile)
					_ = c.sendAck(ctx, ResponseTypeAckUnsub, rq.GridTile)
				default:
					_ = c.sendError(ctx, "Unknown request type")
				}

			case websocket.MessageBinary:
				_ = c.sendError(ctx, "Please speak text")
			}
		}
	}()

	// write a stream of location information

	for planeMsg := range c.outChan {
		_ = c.sendPlaneMessage(ctx, &planeMsg)
	}

	return nil
}

func (c *WsClient) sendAck(ctx context.Context, ackType, tile string) error {
	rs := WsResponse{
		Type:  ackType,
		Tiles: []string{tile},
	}
	return c.sendPlaneMessage(ctx, &rs)
}

func (c *WsClient) sendError(ctx context.Context, msg string) error {
	rs := WsResponse{
		Type:    ResponseTypeError,
		Message: msg,
	}
	return c.sendPlaneMessage(ctx, &rs)
}

func (c *WsClient) sendPlaneMessage(ctx context.Context, planeMsg *WsResponse) error {
	buf, err := json.MarshalIndent(planeMsg, "", "  ")
	if nil != err {
		log.Debug().Err(err).Str("type", planeMsg.Type).Msg("Failed to marshal plane msg to send to client")
		return err
	}
	if err = c.writeTimeout(ctx, 3*time.Second, buf); nil != err {
		log.Debug().Err(err).Str("type", planeMsg.Type).Msg("Failed to send message to client")
		return err
	}
	return nil
}

func (c *WsClient) writeTimeout(ctx context.Context, timeout time.Duration, msg []byte) error {
	ctxW, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.conn.Write(ctxW, websocket.MessageText, msg)
}

func (cl *ClientList) addClient(c *WsClient, out chan WsResponse) {
	cl.clientsLock.Lock()
	defer cl.clientsLock.Unlock()
	cl.clients[c] = out
}

func (cl *ClientList) removeClient(c *WsClient) {
	cl.clientsLock.Lock()
	defer cl.clientsLock.Unlock()
	delete(cl.clients, c)
}

func (cl *ClientList) SendLocationUpdate(highLow, tile string, loc *export.EnrichedPlaneLocation) {
	cl.clientsLock.RLock()
	defer cl.clientsLock.RUnlock()

	for client, outChan := range cl.clients {
		if client.HasSub(tile) || client.HasSub("all"+highLow) {
			outChan <- WsResponse{
				Type:     ResponseTypePlaneLocation,
				Location: loc,
			}
		}
	}
}
