package dummy

import (
	"fmt"
	"net"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type portWaiter struct {
	attempts int
	delay    time.Duration
}

type PortWaiter interface {
	Wait(serviceName string, host string, port int) error
}

func NewPortWaiter(attempts int, delay time.Duration) PortWaiter {
	return &portWaiter{
		attempts: attempts,
		delay:    delay,
	}
}

func (w *portWaiter) Wait(serviceName string, host string, port int) error {
	for i := 0; i < w.attempts; i++ {
		addr := fmt.Sprintf("%s:%d", host, port)

		_, err := net.Dial("tcp", addr)
		if err == nil {
			return nil
		}
		time.Sleep(w.delay)
	}

	return bosherr.Errorf("Timed out waiting for service %s to come up for %d", serviceName, w.attempts)
}
