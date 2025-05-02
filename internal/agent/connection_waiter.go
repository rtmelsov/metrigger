package agent

import (
	"context"
	"fmt"
	"net"
	"time"
)

func WaitForServer(ctx context.Context, address string) error {
	var timeouts = []int{1, 3, 5}
	for _, el := range timeouts {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation canceled: %w", ctx.Err())
		default:
			conn, err := net.Dial("tcp", address)
			if err == nil {
				err := conn.Close()
				if err != nil {
					return err
				}
				return nil // Сервер доступен
			}
			time.Sleep(time.Duration(el) * time.Second)
		}
	}
	return fmt.Errorf("services not available at %s after %v", address, timeouts[len(timeouts)-1])
}
