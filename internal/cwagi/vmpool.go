package cwagi

import (
	"context"
	"expvar"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/pborman/uuid"
)

// VMPool is a group of WebAssembly virtual machines dynamically spun up and down.
type VMPool struct {
	module         []byte
	name, mainFunc string
	vms            []*managedVM
	lock           sync.RWMutex
	maxSize        int
	work           chan workData
	vmCtr          *expvar.Int
	busyCtr        *expvar.Int
	requestCtr     *expvar.Int
}

// NewPool creates a new pool of WebAssembly workers with the given cwagi-linked code.
func NewPool(module []byte, name, mainFunc string, initSize, maxSize int) *VMPool {
	pid := uuid.New()

	vp := &VMPool{
		module:     module,
		name:       name,
		mainFunc:   mainFunc,
		maxSize:    maxSize,
		work:       make(chan workData, 32),
		vmCtr:      expvar.NewInt(pid + "::vmCount"),
		busyCtr:    expvar.NewInt(pid + "::busyCount"),
		requestCtr: expvar.NewInt(pid + "::requestCount"),
	}

	for range make([]struct{}, initSize) {
		vp.createVM()
	}

	return vp
}

func (vp *VMPool) Close() error {
	vp.lock.Lock()
	defer vp.lock.Unlock()

	for _, mvm := range vp.vms {
		mvm.cancel()
	}

	return nil
}

func (vp *VMPool) createVM() (*managedVM, error) {
	vp.lock.Lock()
	defer vp.lock.Unlock()

	log.Println("creating a new VM")
	vs, err := NewVM(vp.module, []string{vp.name, "mode", "cwagi"}, vp.name, vp.mainFunc)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	mvm := &managedVM{
		vms:    vs,
		cancel: cancel,
	}
	go mvm.work(ctx, vp.work, vp.busyCtr)

	vp.vms = append(vp.vms, mvm)
	vp.vmCtr.Add(1)

	return mvm, nil
}

func (vp *VMPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := workData{
		w:    w,
		r:    r,
		done: make(chan struct{}),
	}

	vp.work <- d

	if rand.Int63()%16 == 0 {
		log.Printf("chosen to be a decider!")

		vp.lock.RLock()
		vpWorkLen := len(vp.work)
		vpVmsLen := len(vp.vms)
		vp.lock.RUnlock()

		log.Printf("work: %d, vms: %d", vpWorkLen, vpVmsLen)

		if vpWorkLen > vpVmsLen {
			vp.createVM()
		}
	}

	<-d.done
	vp.requestCtr.Add(1)
}

type managedVM struct {
	busy   bool
	vms    *VMServer
	cancel context.CancelFunc
}

func (mvm *managedVM) work(ctx context.Context, ch <-chan workData, busyCtr *expvar.Int) {
	for {
		select {
		case <-ctx.Done():
			return

		case data := <-ch:
			mvm.busy = true
			busyCtr.Add(1)
			mvm.vms.ServeHTTP(data.w, data.r)
			busyCtr.Add(-1)
			mvm.busy = false
			data.done <- struct{}{}
		}
	}
}

type workData struct {
	w    http.ResponseWriter
	r    *http.Request
	done chan struct{}
}
