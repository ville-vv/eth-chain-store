package schedule

import (
	"github.com/robfig/cron/v3"
	"time"
)

type ScheduleFunc func()

type IWork interface {
	Scheme() string
	Spec() string
	WorkFunc() ScheduleFunc
}

type NewIWork func() IWork

type Schedule struct {
	works map[string]IWork
	c     *cron.Cron
}

func NewSchedule(wks ...IWork) *Schedule {
	sc := &Schedule{works: make(map[string]IWork), c: cron.New(cron.WithLocation(time.Local), cron.WithSeconds())}
	sc.Add(wks...)
	return sc
}

func (sel *Schedule) Add(wks ...IWork) *Schedule {
	for _, wk := range wks {
		sel.add(wk)
	}
	return sel
}

func (sel *Schedule) add(wk IWork) {
	if sel.works == nil {
		sel.works = map[string]IWork{}
	}
	sel.works[wk.Scheme()] = wk
	return
}

func (sel *Schedule) AddNew(fcs ...NewIWork) *Schedule {
	for _, f := range fcs {
		sel.add(f())
	}
	return sel
}

func (sel *Schedule) Start() error {
	c := sel.c
	for _, w := range sel.works {
		_, err := c.AddFunc(w.Spec(), w.WorkFunc())
		if err != nil {
			return err
		}
	}
	c.Run()
	return nil
}
