package ids

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func CreateId(prefix string) string {
	t := time.Unix(10000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return prefix + ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
