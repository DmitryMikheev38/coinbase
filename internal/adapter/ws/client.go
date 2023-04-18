package ws

import (
	"coinbase/internal/models"
	"coinbase/internal/usecase/price"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type request struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []any    `json:"channels"`
}

type Client struct {
	url     string
	priceUC *price.UseCase
}

func NewClient(url string, priceUC *price.UseCase) *Client {
	return &Client{url: url, priceUC: priceUC}
}

func (c *Client) SubscribeToTicketChannels(ctx context.Context, ProductIDs []string) chan error {
	errChan := make(chan error)
	for _, ProductID := range ProductIDs {
		go func(id string) {
			errChan <- c.SubscribeToTicketChannel(ctx, id)
		}(ProductID)
	}

	return errChan
}

func (c *Client) SubscribeToTicketChannel(ctx context.Context, ProductIDs ...string) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := request{
		Type:       "subscribe",
		ProductIDs: ProductIDs,
		Channels:   []any{"ticker"},
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	log.Println("subscribe to:", ProductIDs)
	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			m := &models.Tick{}
			err = json.Unmarshal(message, m)
			if err != nil {
				log.Println("unmarshal:", err)
				continue
			}
			c.priceUC.SaveTick(m)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
