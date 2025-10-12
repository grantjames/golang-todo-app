package stores

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"grantjames.github.io/todo-app/types"
)

func newTestActor(tb testing.TB) (*TodoStoreActor, context.CancelFunc) {
	tb.Helper()
	a := NewTodoStoreActor(NewInMemoryTodoStore())
	ctx, cancel := context.WithCancel(context.Background())
	go a.Run(ctx)
	tb.Cleanup(cancel)
	return a, cancel
}

func TestConcurrentCreateAndList(t *testing.T) {
	t.Parallel()

	a, _ := newTestActor(t)

	const N = 1000
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			resp := make(chan types.AddTodoResponse)
			id := strconv.Itoa(i)
			a.Send(types.AddTodoRequest{Todo: types.NewTodo("Test todo "+id, nil), Resp: resp})
			<-resp
		}()
	}

	wg.Wait()

	// List and verify count
	responseChan := make(chan types.GetAllTodosResponse)
	a.cmds <- types.GetAllTodosRequest{Resp: responseChan}

	res := <-responseChan
	if len(res.Todos) != N {
		t.Fatalf("expected %d todos, got %d", N, len(res.Todos))
	}
}
