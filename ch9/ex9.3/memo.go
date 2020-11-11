package memo

import "fmt"

// Func ...
type Func func(key string, done <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	key           string
	res           result
	cancellist    []<-chan struct{}
	appendRequest chan request
	cancelRequest chan struct{}
	cancel        chan struct{}
	ready         chan struct{}
}

type request struct {
	key      string
	done     <-chan struct{}
	response chan<- result
}

// Memo ...
type Memo struct{ requests chan request }

// New ...
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

// Get ...
func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, done, response}
	res := <-response
	return res.value, res.err
}

// Close ...
func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	cancelEntry := make(chan string)
	for {
		select {
		case req, ok := <-memo.requests:
			if ok {
				e := cache[req.key]
				if e == nil {
					cancellist := make([]<-chan struct{}, 0)
					cancellist = append(cancellist, req.done)
					e = &entry{
						key:           req.key,
						cancellist:    cancellist,
						appendRequest: make(chan request),
						cancelRequest: make(chan struct{}),
						cancel:        make(chan struct{}),
						ready:         make(chan struct{}),
					}
					cache[req.key] = e
					go e.watchCancel(req.key, cancelEntry)
					go e.call(f, req.key)
				} else {
					select {
					case <-e.ready:
					default:
						e.appendRequest <- req
					}
				}

				go e.deliver(req.done, req.response)
			}
		case key, ok := <-cancelEntry:
			if ok {
				delete(cache, key)
			}
		}
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key, e.cancel)
	close(e.ready)
	close(e.appendRequest)
	e.cancellist = nil
}

func (e *entry) deliver(done <-chan struct{}, response chan<- result) {
	exit := false
	for !exit {
		select {
		case <-e.ready:
			//fmt.Printf("ready %s\n", e.key)
			response <- e.res
			exit = true
		case <-done:
			//fmt.Printf("deliver done %s\n", e.key)
		E:
			for {
				select {
				case <-e.ready:
					break E
				case <-e.cancel:
					response <- result{value: nil, err: fmt.Errorf("cancel")}
					break E
				case e.cancelRequest <- struct{}{}:
					response <- result{value: nil, err: fmt.Errorf("cancel")}
					break E
				}
			}
			exit = true
		}
	}
}

func (e *entry) watchCancel(key string, cancelEntry chan<- string) {
	exit := false
	for !exit {
		select {
		case req, ok := <-e.appendRequest:
			if ok {
				//fmt.Printf("append\n")
				e.cancellist = append(e.cancellist, req.done)
			}
		case _, ok := <-e.cancelRequest:
			//fmt.Printf("exit len=%d\n", len(e.cancellist))
			if !ok {
				exit = true
				break
			}
			count := 0
			for _, c := range e.cancellist {
				select {
				case <-c:
					count++
				default:
				}
			}
			//fmt.Printf("count=%d, len=%d\n", count, len(e.cancellist))
			if count == len(e.cancellist) {
				close(e.cancel)
				cancelEntry <- key
				exit = true
			}
		}
	}
}
