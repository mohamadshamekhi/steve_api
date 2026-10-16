package redisarticleview

import (
	"context"
	"fmt"

	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type DB struct {
	adapter *redisadapter.Adapter
}

func New(adapter *redisadapter.Adapter) *DB {
	return &DB{adapter: adapter}
}

// AddViewer records that visitorID has viewed articleID, using a Redis set
// (article:views:<articleID>) as the dedup ledger. It returns true only the
// first time a given visitorID is seen for that article, meaning the caller
// should count this as a new unique view; subsequent calls for the same
// visitor/article pair return false and should not increment anything.
//
// Note: this dedup ledger lives only in Redis. If Redis data is lost (no
// persistence / flush), a returning visitor could be counted again — an
// accepted tradeoff for keeping the hot path cheap. The Mongo view_count
// field remains the durable source of truth for the count itself.
func (d DB) AddViewer(ctx context.Context, articleID, visitorID string) (bool, error) {
	const op = "redisarticleview.AddViewer"

	key := fmt.Sprintf("article:views:%s", articleID)
	added, err := d.adapter.Client().SAdd(ctx, key, visitorID).Result()
	if err != nil {
		return false, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithMeta(map[string]any{"article_id": articleID})
	}

	return added > 0, nil
}
