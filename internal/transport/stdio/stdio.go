package stdio

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zmcp/odata-mcp/internal/transport"
)

// StdioTransport implements the Transport interface for stdio communication
type StdioTransport struct {
	reader  *bufio.Reader
	writer  io.Writer
	handler transport.Handler
}

// New creates a new stdio transport
func New(handler transport.Handler) *StdioTransport {
	return &StdioTransport{
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
		handler: handler,
	}
}

// Start begins processing messages from stdio
func (t *StdioTransport) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := t.ReadMessage()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				// Log error but continue processing
				fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
				continue
			}

			// Process request if it has a method
			if msg.Method != "" && t.handler != nil {
				response, err := t.handler(ctx, msg)
				if err != nil {
					// Send error response
					errorResponse := &transport.Message{
						JSONRPC: "2.0",
						ID:      msg.ID,
						Error: &transport.Error{
							Code:    -32603,
							Message: err.Error(),
						},
					}
					if err := t.WriteMessage(errorResponse); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing error response: %v\n", err)
					}
				} else if response != nil {
					if err := t.WriteMessage(response); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
					}
				}
			}
		}
	}
}

// ReadMessage reads a line-delimited JSON message from stdin
func (t *StdioTransport) ReadMessage() (*transport.Message, error) {
	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var msg transport.Message
	if err := json.Unmarshal(line, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// WriteMessage writes a JSON message to stdout
func (t *StdioTransport) WriteMessage(msg *transport.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if _, err := t.writer.Write(data); err != nil {
		return err
	}

	if _, err := t.writer.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

// Close closes the transport (no-op for stdio)
func (t *StdioTransport) Close() error {
	return nil
}