package pgxpool

import (
	"context"
	"time"
)

// Acquire acquires a connection from the pool. If the connection has exceeded
// MaxConnLifetime, it is destroyed and a new one is attempted.
func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {
	for {
		res, err := p.pool.Acquire(ctx)
		if err != nil {
			return nil, err
		}

		conn := res.Value().(*connResource)
		if p.config.MaxConnLifetime > 0 && time.Since(conn.createdAt) > p.config.MaxConnLifetime {
			res.Destroy()
			continue
		}

		return &Conn{res: res, pool: p}, nil
	}
}