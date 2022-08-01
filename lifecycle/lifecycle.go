package lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhangel/go-frame.git/utils"
	"github.com/ahmetb/go-linq/v3"
)

var (
	lc           = &lifeCycle{}
	hooksTimeout = 30 * time.Second
	finalizing   = uint32(0)
	warnTimeout  = 3 * time.Second
)

type hookFunc struct {
	opts *Options
	fn   interface{}
}

func (s *hookFunc) run(ctx context.Context) {
	invoker := func() chan struct{} {
		done := make(chan struct{}, 1)

		go func() {
			workerCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			defer func() { done <- struct{}{} }()
			slow := uint32(0)

			t := time.NewTimer(warnTimeout)
			defer t.Stop()

			go func() {
				select {
				case <-t.C:
					atomic.StoreUint32(&slow, 1)
					log.Printf("Lifecycle hook %q blocked longer than %v...", s, warnTimeout)
					return
				case <-workerCtx.Done():
					return
				}
			}()

			fnType := reflect.TypeOf(s.fn)
			if fnType.Kind() != reflect.Func {
				log.Printf("Invalid lifecycle hook function.")
				return
			}

			if fnType.NumOut() != 0 {
				log.Printf("Invalid lifecycle hook function, only func() or func(context.Context) accepted.")
				return
			}

			switch fnType.NumIn() {
			case 0:
				reflect.ValueOf(s.fn).Call(nil)
			case 1:
				if fnType.In(0).Kind() != reflect.Interface || !fnType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
					log.Printf("Invalid lifecycle hook function, only func() or func(context.Context) accepted.")
					return
				}
				reflect.ValueOf(s.fn).Call([]reflect.Value{reflect.ValueOf(workerCtx)})
			default:
				log.Printf("Invalid lifecycle hook function, only func() or func(context.Context) accepted.")
				return
			}

			if atomic.LoadUint32(&slow) == 1 && workerCtx.Err() == nil {
				log.Printf("Lifecycle hook %q done.", s)
			}
		}()
		return done
	}

	select {
	case <-invoker():
		return
	case <-ctx.Done():
		log.Println(utils.FullCallStack())
		log.Printf("Lifecycle hook %q timeout exceeded, force terminated.", s)
		os.Exit(1)
		return
	}
}

func (s *hookFunc) String() string {
	if s.opts.name != "" {
		return s.opts.name
	}

	pc := reflect.ValueOf(s.fn).Pointer()
	if funcPtr := runtime.FuncForPC(pc); funcPtr != nil {
		fileName, line := funcPtr.FileLine(pc)
		if idx := strings.Index(fileName, ".git/"); idx != -1 {
			if idx := strings.LastIndex(fileName[:idx], "/"); idx != -1 {
				return fmt.Sprintf("%s:%d", fileName[idx+1:], line)
			} else {
				return funcPtr.Name()
			}
		} else {
			return funcPtr.Name()
		}
	} else {
		return fmt.Sprintf("%p", s.fn)
	}
}

type lifeCycle struct {
	mu sync.RWMutex

	onceInitialize sync.Once
	onceFinalize   sync.Once

	initializeHooks []hookFunc
	finalizeHooks   []hookFunc

	initialized uint32
}

func LifeCycle() *lifeCycle {
	return lc
}

func (s *lifeCycle) Initialize(timeout time.Duration) {
	s.onceInitialize.Do(func() {
		s.mu.RLock()
		initializeHooks := make([]hookFunc, len(s.initializeHooks))
		copy(initializeHooks, s.initializeHooks)
		s.mu.RUnlock()

		if timeout != 0 {
			hooksTimeout = timeout
		}

		runHooks(hooksTimeout, initializeHooks)

		atomic.StoreUint32(&s.initialized, 1)
	})
}

func (s *lifeCycle) Finalize() {
	s.onceFinalize.Do(func() {
		atomic.StoreUint32(&finalizing, 1)
		defer atomic.StoreUint32(&finalizing, 0)

		s.mu.RLock()
		finalizeHooks := make([]hookFunc, len(s.finalizeHooks))
		copy(finalizeHooks, s.finalizeHooks)
		s.mu.RUnlock()

		runHooks(hooksTimeout, finalizeHooks)
	})
}

func (s *lifeCycle) HookInitialize(hook func(), opt ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.initializeHooks = append(s.initializeHooks, hookFunc{
		fn:   hook,
		opts: generateOptions(opt...),
	})
}

func (s *lifeCycle) HookFinalize(hook func(ctx context.Context), opt ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.finalizeHooks = append(s.finalizeHooks, hookFunc{
		fn:   hook,
		opts: generateOptions(opt...),
	})
}

func (s *lifeCycle) IsInitialized() bool {
	return atomic.LoadUint32(&s.initialized) == 1
}

func (s *lifeCycle) IsFinalizing() bool {
	return atomic.LoadUint32(&finalizing) == 1
}

func Exit(code int) {
	LifeCycle().Finalize()
	os.Exit(code)
}

func runHooks(timeout time.Duration, hooks []hookFunc) {
	rootCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	linq.From(hooks).
		OrderByDescendingT(func(hook hookFunc) int64 { return int64(hook.opts.timeout) / 1e6 }).
		GroupByT(
			func(hook hookFunc) int32 { return hook.opts.priority },
			func(hook hookFunc) hookFunc { return hook },
		).
		OrderByDescendingT(func(group linq.Group) int32 { return group.Key.(int32) }).
		ForEachT(func(group linq.Group) {
			groupCtx, cancel := context.WithTimeout(rootCtx, group.Group[0].(hookFunc).opts.timeout)
			defer cancel()

			wg := sync.WaitGroup{}
			wg.Add(len(group.Group))

			for _, h := range group.Group {
				go func(hook hookFunc) {
					defer wg.Done()

					ctx, cancel := context.WithTimeout(groupCtx, hook.opts.timeout)
					defer cancel()

					hook.run(ctx)
				}(h.(hookFunc))
			}
			wg.Wait()
		})
}
