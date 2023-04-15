package ws

import (
	"coinbase/internal/models"
	"coinbase/internal/usecase/price"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"log"
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

func (c *Client) SubscribeToTicketChannels(ctx context.Context, ProductIDs []string) error {
	g, errctx := errgroup.WithContext(ctx)
	for _, ProductID := range ProductIDs {
		func(id string) {
			g.Go(func() error { return c.SubscribeToTicketChannel(errctx, id) })
		}(ProductID)
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
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

	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			m := &models.Tick{}
			err = json.Unmarshal(message, m)
			if err != nil {
				log.Println("json parse:", err)
				return err
			}
			if m.ProductID == "" || m.BestAsk == "" || m.BestBid == "" {
				continue
			}
			c.priceUC.SaveTick(m)
		}
	}
}
