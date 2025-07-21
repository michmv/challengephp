package cmd

import (
	"challengephp/lib"
	"challengephp/src/config"
	"challengephp/src/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func SeedDatabase(configPathDefault string) lib.Command {
	configPath := configPathDefault
	return lib.Command{
		Use:   "seed",
		Short: "Заполняем базу данными",
		Run:   runSeed(configPath),
	}
}

func runSeed(configPath string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig(configPath)
		if err != nil {
			lib.Exit(err.Tap())
		}

		pool, err := db.CreateDB(conf.DB)
		if err != nil {
			lib.Exit(err.Tap())
		}
		defer pool.Close()

		fmt.Println("Старт")

		err = seeding(pool, context.Background())
		if err != nil {
			lib.Exit(err.Tap())
		}

		fmt.Println("Конец")
	}
}

func seeding(pool *pgxpool.Pool, ctx context.Context) lib.Error {
	err := prepareDb(pool, ctx)
	if err != nil {
		return err.Tap()
	}

	resUsers := make(chan lib.Error)
	resTypes := make(chan lib.Error)
	resEvents := make(chan lib.Error)

	go seedUsers(pool, ctx, resUsers)
	go seedTypes(pool, ctx, resTypes)
	go seedEvents(pool, ctx, resEvents)
	err = <-resUsers
	if err != nil {
		return err.Tap()
	}
	err = <-resTypes
	if err != nil {
		return err.Tap()
	}
	err = <-resEvents
	if err != nil {
		return err.Tap()
	}

	err = recoveryDbAfterPrepare(pool, ctx)
	if err != nil {
		return err.Tap()
	}

	return nil
}

func seedUsers(pool *pgxpool.Pool, ctx context.Context, res chan lib.Error) {
	data := make([]string, 0, 100)
	now := time.Now().Format("2006-01-02 15:04:05")
	for i := 1; i <= 1_000; i++ {
		data = append(data, fmt.Sprintf("(%d, 'user%d', '%s')", i, i, now))
	}
	query := "INSERT INTO " + db.TableUsers + " (id, name, created_at) VALUES " + strings.Join(data, ",") + ";"
	_, err := pool.Exec(ctx, query)
	if err != nil {
		fmt.Println(query)
		res <- lib.Err(err)
	}

	res <- nil
}

func seedTypes(pool *pgxpool.Pool, ctx context.Context, res chan lib.Error) {
	data := make([]string, 0, 100)
	for i := 1; i <= 100; i++ {
		data = append(data, fmt.Sprintf("(%d, 'type%d')", i, i))
	}
	query := "INSERT INTO " + db.TableEventTypes + " (id, name) VALUES " + strings.Join(data, ",") + ";"
	_, err := pool.Exec(ctx, query)
	if err != nil {
		fmt.Println(query)
		res <- lib.Err(err)
	}

	res <- nil
}

func seedEvents(pool *pgxpool.Pool, ctx context.Context, res chan lib.Error) {
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	workers := 2000

	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(newCtx, pool, i, &wg, errCh)
	}

	select {
	case err := <-errCh:
		cancel()
		wg.Wait()
		close(errCh)
		res <- lib.Err(err)
	case <-waitGroupDone(&wg):
		close(errCh)
		res <- nil
	}
}

func waitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func worker(ctx context.Context, pool *pgxpool.Pool, offset int, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	total := 5000

	rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now().Format("2006-01-02 15:04:05")
	data := make([]string, 0, total)
	for i := 1; i <= total; i++ {
		articleId := rand.Intn(100) + 1
		userId := rand.Intn(1000) + 1
		typeId := rand.Intn(100) + 1
		data = append(data, fmt.Sprintf("(%d, '%s', '{\"page\": \"/article%d\"}', %d, %d)", i+(offset*total), now, articleId, userId, typeId))
	}
	query := "INSERT INTO " + db.TableEvents + " (id, timestamp, metadata, user_id, type_id) VALUES " + strings.Join(data, ",") + ";"

	select {
	case <-ctx.Done():
		return
	default:
	}

	// вставляем sql
	_, err := pool.Exec(ctx, query)
	if err != nil {
		errCh <- err
		return
	}
}

func prepareDb(pool *pgxpool.Pool, ctx context.Context) lib.Error {
	for _, q := range sliceStringMerge(db.PrepareSql, db.DropIndex) {
		_, err := pool.Exec(ctx, q)
		if err != nil {
			return lib.Err(err)
		}
	}
	return nil
}

func recoveryDbAfterPrepare(pool *pgxpool.Pool, ctx context.Context) lib.Error {
	for _, q := range sliceStringMerge(db.CreateIndex, db.RecoveryDbAfterPrepare) {
		_, err := pool.Exec(ctx, q)
		if err != nil {
			return lib.Err(err)
		}
	}
	return nil
}

func sliceStringMerge(a []string, b []string) []string {
	for _, i := range b {
		a = append(a, i)
	}
	return a
}
