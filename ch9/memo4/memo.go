package memo

import (
	"sync"
)

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		/*
		   이 key에 대한 최초 요청
		   이 고루틴이 값을 계산하고 ready 상태를 알려야 한다
		*/
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e // memo.cache[key]에 nil이 아닌 값을 넣어두고 Unlock
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // ready 상태 브로드캐스트
	} else {
		memo.mu.Unlock() // key에 대한 요청이 nil이 아니라면 우선 Unlock
		_ = <-e.ready    // ready 채널이 close 되어서 제로값이 반환될때까지 대기. key당 f() 호출은 1번만 일어나게 된다.
	}
	return e.res.value, e.res.err
}
