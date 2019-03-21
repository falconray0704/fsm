package fsmBenchmark_test

import (
	"errors"
	falconFSM "github.com/falconray0704/fsm"
	loopLabFSM "github.com/looplab/fsm"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
	"testing"
)

func initSys() (clean func()){
	var (
		errOnce error
	)

	sysLoggerCfg := slog.NewSysLogCfg()
	errOnce = sysCfg.LoadFileCfgs("./sysDatas/cfgs/appCfgs.yaml", "sysLogger", &sysLoggerCfg)
	if errOnce != nil {
		log.Fatalf("Loading system configs fail:%s\n", errOnce.Error())
	}
	_, _, err := slog.Init(sysLoggerCfg)
	if err != nil {
		log.Fatalf("Init system logger fail: %s.\n", err.Error())
	}
	slog.Debug("System init finished.")

	return func() {
		slog.Sync()
		slog.Close()
	}

}

var (
	sysClean func()
)

func init() {
	sysClean = initSys()
}

func callbackSetBool_falconFSM(event *falconFSM.Event) {
	var (
		v *bool
		ok bool
	)
	if v, ok = event.Args[0].(*bool); ok {
		*v = false
	} else {
		event.Err = errors.New("callback args error")
	}
}

func Benchmark_falconFSM(b *testing.B) {
	var (
		err error
		fsm *falconFSM.FSM
	)
	const (
		StateStartID = iota
		StateOpened = iota
		StatePaused = iota
		StateClosed = iota
		StateNonExist
	)
	const (
		StateStrOpened = "opened"
		StateStrPaused = "paused"
		StateStrClosed = "closed"
	)
	const (
		EventStartID = iota
		EventOpen = iota
		EventPause = iota
		EventClose = iota
		EventNonExist
	)
	const (
		EventStrOpen = "open"
		EventStrPause = "paused"
		EventStrClose = "close"
	)

	fsm, err = falconFSM.NewFSM(
		StateClosed,
		falconFSM.EventMap{
			EventStartID: falconFSM.EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		falconFSM.StateMap{
			StateStartID: falconFSM.StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		falconFSM.Events{
			{IDEvent: EventOpen, IDsSrc:[]falconFSM.StateID{StateClosed, StatePaused}, IDDst:StateOpened},
			{IDEvent: EventPause, IDsSrc:[]falconFSM.StateID{StateOpened}, IDDst:StatePaused},
			{IDEvent: EventClose, IDsSrc:[]falconFSM.StateID{StateOpened, StatePaused}, IDDst:StateClosed},
		},
		falconFSM.Callbacks{
			{IDCallbackType: falconFSM.CallbackBeforeEvent, ID: EventOpen}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackBeforeEvent, ID: EventPause}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackBeforeEvent, ID: EventClose}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackLeaveState, ID: StateOpened}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackLeaveState, ID: StatePaused}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackLeaveState, ID: StateClosed}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackEnterState, ID: StateOpened}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackEnterState, ID: StatePaused}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackEnterState, ID: StateClosed}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackAfterEvent, ID: EventOpen}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackAfterEvent, ID: EventPause}: callbackSetBool_falconFSM,
			{IDCallbackType: falconFSM.CallbackAfterEvent, ID: EventClose}: callbackSetBool_falconFSM,
		})

	///////////////
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var (
			a, b, c bool
		)
		if err = fsm.Event(EventOpen, &a); err != nil || a != false {break}
		if err = fsm.Event(EventPause, &b); err != nil || b != false {break}
		if err = fsm.Event(EventClose, &c); err != nil || c != false {break}
	}
}

func callbackSetBool_loopLabFSM(event *loopLabFSM.Event) {
	var (
		v *bool
		ok bool
	)
	if v, ok = event.Args[0].(*bool); ok {
		*v = false
	} else {
		event.Err = errors.New("callback args error")
	}
}

func Benchmark_loopLabFSM(b *testing.B) {
	var (
		err error
		fsm *loopLabFSM.FSM
	)
	const (
		StateStartID = iota
		StateOpened = iota
		StatePaused = iota
		StateClosed = iota
		StateNonExist
	)
	const (
		StateStrOpened = "opened"
		StateStrPaused = "paused"
		StateStrClosed = "closed"
	)
	const (
		EventStartID = iota
		EventOpen = iota
		EventPause = iota
		EventClose = iota
		EventNonExist
	)
	const (
		EventStrOpen = "open"
		EventStrPause = "paused"
		EventStrClose = "close"
	)

	fsm = loopLabFSM.NewFSM(
		StateStrClosed,
		loopLabFSM.Events{
			{Name: EventStrOpen, Src: []string{StateStrClosed, StateStrPaused}, Dst: StateStrOpened},
			{Name: EventStrPause, Src: []string{StateStrOpened}, Dst: StateStrPaused},
			{Name: EventStrClose, Src: []string{StateStrOpened, StateStrPaused}, Dst: StateStrClosed},
		},
		loopLabFSM.Callbacks{
			"before_" + EventStrOpen: callbackSetBool_loopLabFSM,
			"before_" + EventStrPause: callbackSetBool_loopLabFSM,
			"before_" + EventStrClose: callbackSetBool_loopLabFSM,
			"leave_" + StateStrOpened: callbackSetBool_loopLabFSM,
			"leave_" + StateStrPaused: callbackSetBool_loopLabFSM,
			"leave_" + StateStrClosed: callbackSetBool_loopLabFSM,
			"enter_" + StateStrOpened: callbackSetBool_loopLabFSM,
			"enter_" + StateStrPaused: callbackSetBool_loopLabFSM,
			"enter_" + StateStrClosed: callbackSetBool_loopLabFSM,
			"after_" + EventStrOpen: callbackSetBool_loopLabFSM,
			"after_" + EventStrPause: callbackSetBool_loopLabFSM,
			"after_" + EventStrClose: callbackSetBool_loopLabFSM,
		},
	)

	///////////////
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var (
			a, b, c bool
		)
		if err = fsm.Event(EventStrOpen, &a); err != nil || a != false {break}
		if err = fsm.Event(EventStrPause, &b); err != nil || b != false {break}
		if err = fsm.Event(EventStrClose, &c); err != nil || c != false {break}
	}
}



