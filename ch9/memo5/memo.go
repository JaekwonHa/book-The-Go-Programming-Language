package memo

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f) // New로 Memo 객체를 생성할때마다 관리 고루틴인 memo.server가 생긴다
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry) //cache 변수의 scope을 해당 관리 고루틴 내부로 제한한다
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// 관리 고루틴은 memo 객체당 1개이기 때문에 memo.reqeusts 채널로 온 것을 한번에 1건씩 처리한다
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // e.call 에서는 f 함수를 수행하고 e.ready 채널을 브로드캐스트한다
		}
		go e.deliver(req.response) // e.deliver 은 e.ready 가 브로드캐스트 되었다면 바로 캐시된 값을 바로 반환한다
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
