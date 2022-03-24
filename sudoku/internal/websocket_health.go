package sudoku

import (
	"context"
)

func init() {
	websocketPool.Add((*websocketHealthRequest)(nil), (*websocketHealthResponse)(nil))
}

type websocketHealthRequest string

func (websocketHealthRequest) Method() string {
	return "health"
}

func (r websocketHealthRequest) Validate(ctx context.Context) error {
	return nil
}

func (r websocketHealthRequest) Execute(ctx context.Context) (websocketResponse, error) {
	//srv := ctx.Value("srv").(*Service)
	return websocketHealthResponse("OK"), nil
}

// TODO handle and test
type websocketHealthResponse string

func (websocketHealthResponse) Method() string {
	return "health"
}

func (r websocketHealthResponse) Validate(ctx context.Context) error {
	return nil
}

func (r websocketHealthResponse) Execute(ctx context.Context) error {
	return nil
}
