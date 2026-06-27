package sync

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/cache"

	"github.com/redis/go-redis/v9"
)

// cursorTTL keeps each day's cursor/done keys around long enough to survive a
// restart, then lets them self-clean.
const cursorTTL = 48 * time.Hour

func pageKey(entity, date string) string { return fmt.Sprintf("sync:%s:page:%s", entity, date) }
func doneKey(entity, date string) string { return fmt.Sprintf("sync:%s:done:%s", entity, date) }

// NextPage returns the next page to fetch for entity/date, defaulting to 1.
func NextPage(ctx context.Context, entity, date string) (int, error) {
	v, err := cache.RDB.Get(ctx, pageKey(entity, date)).Int()
	if errors.Is(err, redis.Nil) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	if v < 1 {
		return 1, nil
	}
	return v, nil
}

// SetNextPage stores the next page to fetch (called only after a page is
// committed to MySQL, so the cursor always points at the first un-imported page).
func SetNextPage(ctx context.Context, entity, date string, page int) error {
	return cache.RDB.Set(ctx, pageKey(entity, date), page, cursorTTL).Err()
}

// IsDone reports whether entity finished its run for date.
func IsDone(ctx context.Context, entity, date string) (bool, error) {
	n, err := cache.RDB.Exists(ctx, doneKey(entity, date)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// MarkDone flags entity as finished for date.
func MarkDone(ctx context.Context, entity, date string) error {
	return cache.RDB.Set(ctx, doneKey(entity, date), "1", cursorTTL).Err()
}

// ClearCursor removes the page + done keys so the next run starts from page 1.
func ClearCursor(ctx context.Context, entity, date string) error {
	return cache.RDB.Del(ctx, pageKey(entity, date), doneKey(entity, date)).Err()
}
