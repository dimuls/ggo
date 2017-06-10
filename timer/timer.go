package timer

import (
	"errors"
	"time"
)

type mode int

const (
	base mode = iota
	period
)

type Parameters struct {
	Base    int `json:"base"`
	ByoYomi int `json:"byoYomi"`
	Periods int `json:"periods"`
	Moves   int `json:"moves"`
}

type Callbacks struct {
	OnBaseOver   func()
	OnPeriodOver func()
	OnOver       func()
}

type Timer struct {
	parameters Parameters
	callbacks  Callbacks

	base    time.Duration
	byoYomi time.Duration
	periods int
	moves   int

	mode mode

	over bool

	startedAt time.Time

	timer *time.Timer
}

func NewTimer(parameters Parameters, callbacks Callbacks) (*Timer, error) {

	if parameters.Base < 0 {
		return nil, errors.New("base should be greater or equal to zero")
	}

	if parameters.ByoYomi < 0 {
		return nil, errors.New("byo-yomi should be greater or equal to zero")
	}

	if parameters.Base == 0 && parameters.ByoYomi == 0 {
		return nil, errors.New("both base and byo-yomi duration can't be zero")
	}

	if parameters.ByoYomi == 0 {
		if parameters.Periods != 0 {
			return nil, errors.New("periods should be zero")
		}
		if parameters.Moves != 0 {
			return nil, errors.New("moves should be zero")
		}
	} else {
		if parameters.Periods < 1 {
			return nil, errors.New("periods should be greater than zero")
		}
		if parameters.Moves < 1 {
			return nil, errors.New("moves should be greater than one")
		}
		if parameters.Periods > 1 && parameters.Moves > 1 {
			return nil, errors.New("both periods and moves can't be greater than one")
		}
	}

	t := &Timer{
		parameters: parameters,
		callbacks:  callbacks,
		base:       time.Duration(parameters.Base) * time.Second,
		byoYomi:    time.Duration(parameters.ByoYomi) * time.Second,
		periods:    parameters.Periods,
		moves:      parameters.Moves,
		over:       false,
		timer:      nil,
	}

	if parameters.Base > 0 {
		t.mode = base
	} else {
		t.mode = period
	}

	return t, nil
}

func (t *Timer) Switch() {
	if t.over {
		return
	}
	if t.timer == nil {
		t.switchOn()
	} else {
		t.switchOff()
	}
}

func (t *Timer) switchOn() {
	if t.mode == base {
		t.startBaseTimer()
	} else {
		t.startPeriodTimer()
	}
	t.startedAt = time.Now()
}

func (t *Timer) startBaseTimer() {
	t.timer = time.AfterFunc(t.base, t.onBaseOver)
}

func (t *Timer) startPeriodTimer() {
	t.timer = time.AfterFunc(t.byoYomi, t.onPeriodOver)
}

func (t *Timer) switchOff() {
	if t.mode == base {
		t.stopBaseTimer()
	} else {
		t.stopPeriodTimer()
		t.moves--
		if t.moves == 0 {
			t.moves = t.parameters.Moves
			t.byoYomi = time.Duration(t.parameters.ByoYomi) * time.Second
		}
	}
	t.timer = nil
}

func (t *Timer) stopBaseTimer() {
	t.timer.Stop()
	t.base = t.base - time.Now().Sub(t.startedAt)
}

func (t *Timer) stopPeriodTimer() {
	t.timer.Stop()
	t.byoYomi = t.byoYomi - time.Now().Sub(t.startedAt)
}

func (t *Timer) onBaseOver() {
	if t.parameters.ByoYomi == 0 {
		t.over = true
		if t.callbacks.OnOver != nil {
			t.callbacks.OnOver()
		}
	} else {
		t.mode = period
		t.startPeriodTimer()
		if t.callbacks.OnBaseOver != nil {
			t.callbacks.OnBaseOver()
		}
	}
}

func (t *Timer) onPeriodOver() {
	t.periods--
	if t.periods == 0 {
		t.over = true
		if t.callbacks.OnOver != nil {
			t.callbacks.OnOver()
		}
	} else {
		t.byoYomi = time.Duration(t.parameters.ByoYomi) * time.Second
		t.startPeriodTimer()
		if t.callbacks.OnPeriodOver != nil {
			t.callbacks.OnPeriodOver()
		}
	}
}
