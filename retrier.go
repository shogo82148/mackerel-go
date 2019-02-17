package mackerel

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"
)

// basic implementaion of exponential back off.
type retrier struct {
	delay time.Duration
}

var (
	mu          sync.Mutex
	retrierRand *rand.Rand
)

func init() {
	var seed int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
		seed = time.Now().UnixNano() // fall back to timestamp
	}
	retrierRand = rand.New(rand.NewSource(seed))
}

func (r *retrier) Next(ctx context.Context) bool {
	if r.delay == 0 {
		r.delay = 100 * time.Millisecond
		return true
	}

	mu.Lock()
	jitter := time.Duration(retrierRand.Float64() * float64(time.Second))
	mu.Unlock()

	timer := time.NewTimer(r.delay + jitter)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
		return false
	}
	r.delay *= 2
	if r.delay >= 60*time.Second {
		r.delay = 60 * time.Second
	}

	return true
}
