package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	cfg := parseConfig()
	finalTimeout := time.Duration(cfg.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), finalTimeout)
	defer cancel()
	resultChan := make(chan Result, len(cfg.URLs))
	for _, url := range cfg.URLs {
		go fetchHedged(ctx, url, resultChan)

	}
	errorsCount := 0

	for i := 0; i < len(cfg.URLs); i++ {
		select {
		case res := <-resultChan:
			if res.Err != nil {
				errorsCount++
				continue
			}

			cancel()

			fmt.Printf("%s %s (получено от %s)\n", res.Proto, res.Status, res.URL)
			fmt.Println("Заголовки ответа:")
			for key, values := range res.Headers {
				for _, value := range values {
					fmt.Printf("  %s: %s\n", key, value)
				}
			}
			fmt.Println("\nТело ответа:")
			fmt.Println(string(res.Body))
			os.Exit(0)

		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "Ошибка: Превышен общий таймаут ожидания (%v)\n", 15)
			os.Exit(228)
		}
	}

	if errorsCount == len(cfg.URLs) {
		fmt.Fprintf(os.Stderr, "Ошибка: Все указанные URL (%d шт.) завершились ошибкой запроса.\n", errorsCount)
		os.Exit(1)
	}
}
