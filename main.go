package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	fs := flag.NewFlagSet("flags", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	from := fs.Int("from", 0, "Начальный ID")
	to := fs.Int("to", 0, "Конечный ID")
	workers := fs.Int("workers", 10, "Количество воркеров")
	timeout := fs.Int("timeout", 5, "Таймаут HTTP запроса")

	// проверим что все значения у флагов - числа
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println("Ошибка чтения флагов")
		fmt.Println("Корректный ввод: --from=N --to=M --workers=X --timeout=T")
		fmt.Println("--workers=X и --timeout==T указываются опционально")
		os.Exit(1)
	}

	// проверим что from и to - допустимые
	if *from <= 0 || *to <= 0 || *from > *to {
		fmt.Println("Некорректный ID")
		fmt.Println("Формат ввода: --from=N --to=M")
		fmt.Println("Где (N, M > 0) && (N <= M)")
		os.Exit(1)
	}

	// проверим число воркеров
	if *workers <= 0 {
		fmt.Println("Число воркеров не может быть меньше 1")
		os.Exit(1)
	}

	// проверим таймаут
	if *timeout <= 0 {
		fmt.Println("На HTTP запрос нужно выделять как минимум 1 секунду")
		os.Exit(1)
	}

	// Переведем timeout в секунды
	duration := time.Duration(*timeout) * time.Second

	// создадим контекст с защитой от Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// создадим 2 канала
	// jobs - для фильмов находящихся в работе, ограниченный количеством воркеров
	// results - для фильмов которые прошли обработку
	jobs := make(chan int, *workers)
	results := make(chan Result)

	var wg sync.WaitGroup

	// создадим горутины в соответствии с размером пула воркеров
	for cnt := 0; cnt < *workers; cnt++ {
		wg.Add(1)
		// будем запускать функцию с отложенным Done()
		go func() {
			defer wg.Done()
			getMovies(ctx, jobs, results, duration)
		}()
	}

	// и после того как прочитаем все сохраненные результаты
	// а именно когда горутины все горутины будут выполнены
	// закроем results
	go func() {
		wg.Wait()
		close(results)
	}()

	// сложим в jobs наши id а по окончании цикла закроем канал
	// так как мы уверены что ничего нового в него класться не будет0
	go func() {
		defer close(jobs)
		for id := *from; id <= *to; id++ {
			select {
			case <-ctx.Done():
				return

			case jobs <- id:
			}
		}
	}()

	// переберем полученные результаты из results
	for res := range results {
		// проверим на наличие ошибок
		if res.Err != nil {
			fmt.Println(res.Err)
			continue
		}

		// в случае если все хорошо выведем фильмы в нужном формате
		fmt.Printf(
			"%d - %s - %d - %s\n",
			res.Movie.ID,
			res.Movie.Title,
			res.Movie.Year,
			res.Movie.Director,
		)
	}
}
