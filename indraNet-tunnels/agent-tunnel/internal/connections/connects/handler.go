package connects

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	pb "agent-tunnel/proto"

	"github.com/gorilla/websocket"
	protobuf "google.golang.org/protobuf/proto"
)

func (c *TunnelClient) handleMessage(ctx context.Context, msg *pb.Envelope) error {

	switch m := msg.Message.(type) {
	case *pb.Envelope_Request:
		return c.handleConnect(ctx, m.Request)
	case *pb.Envelope_Data:
		return c.handleStream(ctx, m.Data)
	case *pb.Envelope_Close:
		return c.handleClose(ctx, m.Close)
	case *pb.Envelope_Control:
		log.Printf("[TUNNEL] control msg: %v", msg.Message)
	default:
		return fmt.Errorf("unknown message type: %T", m)
	}

	return nil
}

func (c *TunnelClient) handleConnect(_ context.Context, req *pb.TunnelRequest) error {

	conn, err := net.Dial("tcp", "localhost:"+c.Cfg.Portforward)
	if err != nil {
		return fmt.Errorf("dial [ERROR]: %w", err)
	}

	c.mu.Lock()
	c.Streams[req.Id] = conn
	defer c.mu.Unlock()

	// if len(req.Body) > 0 {
	// 	conn.Write(req.Body)
	// }

	if len(req.Body) > 0 {
		n, err := conn.Write(req.Body)
		log.Printf("[handleConnect] Sent %d bytes to Flask app for stream %s, err=%v", n, req.Id, err)
	} else {
		log.Printf("[handleConnect] Request body is empty for stream %s", req.Id)
	}

	// Create a new TCP connection for the stream
	go c.pipeToLocal(req.Id, conn)

	log.Printf("Stream %s connected [INFO]", req.Id)
	return nil
}

func (c *TunnelClient) handleStream(_ context.Context, chunk *pb.TunnelData) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, ok := c.Streams[chunk.Id]; ok {
		if _, err := conn.Write(chunk.Chunk); err != nil {
			return fmt.Errorf("write to stream %s: %w", chunk.Id, err)
		}
	} else {
		log.Printf("Received data for unknown stream ID: %s", chunk.Id)
	}
	return nil
}

func (c *TunnelClient) handleClose(_ context.Context, closeMsg *pb.TunnelClose) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, ok := c.Streams[closeMsg.Id]; ok {
		conn.Close()
		delete(c.Streams, closeMsg.Id)
	}
	return nil
}

func (c *TunnelClient) pipeToLocal(id string, conn net.Conn) error {
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		log.Printf("Failed to parse HTTP response: %v", err)
		return err
	}

	// Step 1: send TunnelResponse
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = strings.Join(v, ", ")
	}

	respMsg := &pb.TunnelResponse{
		Id:      id,
		Status:  int32(resp.StatusCode),
		Headers: headers,
	}

	env := &pb.Envelope{
		Message: &pb.Envelope_Response{Response: respMsg},
	}

	_ = c.send_envalope(env)

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Stream %s reached EOF", id)

				closeMsg := &pb.TunnelClose{Id: id}
				env := &pb.Envelope{
					Message: &pb.Envelope_Close{Close: closeMsg},
				}
				_ = c.send_envalope(env)

				return nil
			}
			log.Printf("Error reading from stream %s: %v", id, err)
			return err
		}

		dataMsg := &pb.TunnelData{
			Id:    id,
			Chunk: buf[:n],
		}

		env := &pb.Envelope{
			Message: &pb.Envelope_Data{Data: dataMsg},
		}

		_ = c.send_envalope(env)
	}
}

func (c *TunnelClient) send_envalope(env *pb.Envelope) error {
	b, err := protobuf.Marshal(env)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.BinaryMessage, b)
}
