package id

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	pool *sync.Pool
)

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &generator{r: ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)}
		},
	}
}

type generator struct {
	r io.Reader
}

func (g *generator) New() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), g.r)
}

func MakeULID() ulid.ULID {
	g := pool.Get().(*generator)
	id := g.New()
	pool.Put(g)
	return id
}
