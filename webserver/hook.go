package webserver

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Hook func(ctx context.Context) error

func BuildShutdownServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		wg.Add(len(servers))
		doneChan := make(chan interface{})
		for _, server := range servers {
			go func(s Server) {
				err := s.Shutdown(ctx)
				if err != nil {
					fmt.Printf("shutdown error: %v", err)
				}
				wg.Done()
			}(server)
		}
		go func() {
			wg.Wait()
			doneChan <- "done"
		}()

		select {
		case <-ctx.Done():
			fmt.Printf("close timeout\n")
			return errors.New("close timeout")
		case <-doneChan:
			fmt.Printf("close success\n")
			return nil
		}
	}
}
