package taskrunner

type Runner struct {
	Controller ControlChan
	Error      ControlChan
	Data       DataChan
	DataSize   int
	Longlived  bool
	Dispatcher Fn
	Executor   Fn
}

func NewRunner(size int, longlived bool, d Fn, e Fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		Longlived:  longlived,
		DataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) StartDispatcher() {
	defer func() {
		if !r.Longlived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}

		case e := <-r.Error:
			if e == CLOSE {
				return
			}
		default:

		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.StartDispatcher()
}
