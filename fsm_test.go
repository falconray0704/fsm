// Copyright (c) 2019 - Ray Ruan <falconray@yahoo.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsm

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func TestFSM_Transition(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
			{IDEvent: EventPause, IDsSrc:[]StateID{StateOpened}, IDDst:StatePaused},
			{IDEvent: EventClose, IDsSrc:[]StateID{StateOpened, StatePaused}, IDDst:StateClosed},
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateOpened}, IDDst:StateOpened},
		},
		Callbacks{
			{IDCallbackType: CallbackBeforeEvent, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventPause}: func(e *Event) { fmt.Println("Before event pause.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventClose}: func(e *Event) { fmt.Println("Before event close.")},
			{IDCallbackType: CallbackLeaveState, ID: StateOpened}: func(e *Event) { fmt.Println("Leave state opened.")},
			{IDCallbackType: CallbackLeaveState, ID: StatePaused}: func(e *Event) { fmt.Println("Leave state paused.")},
			{IDCallbackType: CallbackLeaveState, ID: StateClosed}: func(e *Event) { fmt.Println("Leave state closed."); e.Async()},
			{IDCallbackType: CallbackEnterState, ID: StateOpened}: func(e *Event) { fmt.Println("Got into state opened.")},
			{IDCallbackType: CallbackEnterState, ID: StatePaused}: func(e *Event) { fmt.Println("Got into state paused.")},
			{IDCallbackType: CallbackEnterState, ID: StateClosed}: func(e *Event) { fmt.Println("Got into state closed.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventOpen}: func(e *Event) { fmt.Println("After event open.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventPause}: func(e *Event) { fmt.Println("After event pause.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventClose}: func(e *Event) { fmt.Println("After event close.")},
		})
	assert.NoError(t, err, "NewFSM() expect no error.")
	err = fsm.Event(EventOpen)
	assert.Equal(t, AsyncError{}, err, "Async() transition expect AsyncError")
	assert.Equal(t, "async started", err.Error(), "Async Transition expect nil error inside AsyncError")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Async() transition expect keep same state before Transition()")

	// InTransitionError during in Async Transition
	err = fsm.Event(EventOpen)
	assert.Equal(t, InTransitionError{Event: EventStrOpen}, err, "During Async Transition expect InTransitionError if another Event happen")

	// Async Transition()
	err = fsm.Transition()
	assert.NoError(t, err, "Transition() expect no error ")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Async() Transition() expect StateStrOpened")

}

func TestNewFSM(t *testing.T) {
	var (
		stbe *stubbedCallbackEvent
		stbs *stubbedCallbackState
		stbErr error
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


	fsm, err := NewFSM(
		StateStartID,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
		},
		Callbacks{
		})
	assert.Equal(t, StateStartReserveError{}, err, "Init state should not used reserve StateStartID")
	assert.Nil(t, fsm, "Init state used reserve StateStartID expect nil fsm")

	fsm, err = NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
			{IDEvent: EventPause, IDsSrc:[]StateID{StateOpened}, IDDst:StatePaused},
			{IDEvent: EventClose, IDsSrc:[]StateID{StateOpened, StatePaused}, IDDst:StateClosed},
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateOpened}, IDDst:StateOpened},
		},
		Callbacks{
			{IDCallbackType: CallbackBeforeEvent, ID: EventStartID}: func(e *Event) { fmt.Println("Before all event.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventPause}: func(e *Event) { fmt.Println("Before event pause.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventClose}: func(e *Event) { fmt.Println("Before event close.")},
			{IDCallbackType: CallbackLeaveState, ID: StateStartID}: func(e *Event) { fmt.Println("Leave all state.")},
			{IDCallbackType: CallbackLeaveState, ID: StateOpened}: func(e *Event) { fmt.Println("Leave state opened.")},
			{IDCallbackType: CallbackLeaveState, ID: StatePaused}: func(e *Event) { fmt.Println("Leave state paused.")},
			{IDCallbackType: CallbackLeaveState, ID: StateClosed}: func(e *Event) { fmt.Println("Leave state closed.")},
			{IDCallbackType: CallbackEnterState, ID: StateStartID}: func(e *Event) { fmt.Println("Got into all state.")},
			{IDCallbackType: CallbackEnterState, ID: StateOpened}: func(e *Event) { fmt.Println("Got into state opened.")},
			{IDCallbackType: CallbackEnterState, ID: StatePaused}: func(e *Event) { fmt.Println("Got into state paused.")},
			{IDCallbackType: CallbackEnterState, ID: StateClosed}: func(e *Event) { fmt.Println("Got into state closed.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventStartID}: func(e *Event) { fmt.Println("After all event.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventOpen}: func(e *Event) { fmt.Println("After event open.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventPause}: func(e *Event) { fmt.Println("After event pause.")},
			{IDCallbackType: CallbackAfterEvent, ID: EventClose}: func(e *Event) { fmt.Println("After event close.")},
		})
	assert.NoError(t, err, "NewFSM() expect no error.")

	// Current() Is() Can() Cannot()
	assert.Equal(t, StateStrClosed, fsm.Current(), "Init state expect closed")
	assert.Equal(t, true, fsm.Is(StateClosed), "Init state expect StateClosed")
	assert.Equal(t, true, fsm.Can(EventOpen), "Open event transition Can() from StateClosed expect true")
	assert.Equal(t, false, fsm.Cannot(EventOpen), "Open event transition Cannot() from StateClosed expect false")
	assert.Equal(t, false, fsm.Can(EventClose), "Close event transition Can() from StateClosed expect false")
	assert.Equal(t, true, fsm.Cannot(EventClose), "Close event transition Cannot() from StateClosed expect true")

	// SetState()
	err = fsm.SetState(StateNonExist)
	assert.Equal(t, err, StateOutOfRangeError{ID:StateNonExist}, "NonExist state expect error")
	err = fsm.SetState(StateOpened)
	assert.NoError(t, err, "Set valid state StateOpened expect no error")
	err = fsm.SetState(StateClosed)
	assert.NoError(t, err, "Set valid state StateClosed expect no error")
	err = fsm.SetState(StateStartID)
	assert.Equal(t,StateStartReserveError{}, err, "Set reserve state StateStartID expect StateStartReserveError")

	// Event() ---> EventOutOfRangeError
	err = fsm.Event(EventNonExist)
	assert.Equal(t, EventOutOfRangeError{ID:EventNonExist}, err, "Non exist event expect EventOutOfRangeError")

	// Event() ---> InTransitionError
	err = fsm.Event(EventClose)
	assert.Equal(t, InvalidEventError{Event:EventStrClose, State: StateStrClosed}, err,"Not registered transition expect InvalidEventError")


	// closed ---> opened, success
	stbe, stbErr = withStubbedCallbackEvent(fsm, fsm.Event, EventOpen, CallbackBeforeEvent, EventOpen)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbe.err
	// err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")
	// opened ---> opened, NoTransitionError
	err = fsm.Event(EventOpen)
	assert.Equal(t, NoTransitionError{}, err, "Open transition from opened expect NoTransitionError.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")

	// opened ---> closed, success
	stbe, stbErr = withStubbedCallbackEvent(fsm, fsm.Event, EventClose, CallbackAfterEvent, EventClose)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbe.err
	// err = fsm.Event(EventClose)
	assert.NoError(t, err, "Close transition from opened expect success.")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Close transition expect closed")

	// closed ---> opened , success
	stbs, stbErr = withStubbedCallbackState(fsm, fsm.Event, EventOpen, CallbackLeaveState, EventOpen)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbs.err
	// err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")
	// opened ---> paused, success
	stbs, stbErr = withStubbedCallbackState(fsm, fsm.Event, EventPause, CallbackEnterState, EventPause)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbs.err
	//err = fsm.Event(EventPause)
	assert.NoError(t, err, "Pause transition from opened expect success.")
	assert.Equal(t, StateStrPaused, fsm.Current(), "Pause transition expect paused")

	// AvailableTransitions()
	trans := fsm.AvailableTransitions()
	expectTrans := []string{EventStrClose, EventStrOpen}
	expectTrans2 := []string{EventStrOpen, EventStrClose}
	assert.Equal(t, len(expectTrans), len(trans), "Valid transitions count expect 2 in state paused")
	assert.Equal(t, true,
		reflect.DeepEqual(trans, expectTrans) ||
		reflect.DeepEqual(trans, expectTrans2),
		"Valid transitions in state paused expect {open, close}")

	// paused --->  opened success
	stbe, stbErr = withStubbedCallbackEvent(fsm, fsm.Event, EventOpen, CallbackBeforeEvent, EventStartID)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbe.err
	// err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from paused expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")

	// opened ---> paused success
	err = fsm.Event(EventPause)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrPaused, fsm.Current(), "Pause transition expect paused")

	// paused --->  closed success
	stbe, stbErr = withStubbedCallbackEvent(fsm, fsm.Event, EventClose, CallbackAfterEvent, EventStartID)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbe.err
	// err = fsm.Event(EventClose)
	assert.NoError(t, err, "Close transition from paused expect success.")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Close transition expect closed")


	// closed ---> opened, success
	stbs, stbErr = withStubbedCallbackState(fsm, fsm.Event, EventOpen, CallbackLeaveState, StateStartID)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbs.err
	// err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")

	// opened ---> closed, success
	stbs, stbErr = withStubbedCallbackState(fsm, fsm.Event, EventClose, CallbackEnterState, StateStartID)
	assert.NoError(t, stbErr, "Expect no err on stub")
	assert.True(t, stbe.isCalled, "Expect stubbed callback was called")
	err = stbs.err
	// err = fsm.Event(EventClose)
	assert.NoError(t, err, "Close transition from opened expect success.")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Close transition expect closed")
}

func TestNewFSM_buildUpTransitions_DuplicateTransitionError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t,DuplicateTransitionError{event:EventStrOpen, state:StateStrClosed}, err, "Duplicate transition NewFSM() expect DuplicateTransitionError.")
	assert.Nil(t, fsm, "Duplicate transition NewFSM() expect nil fsm.")
}

func TestNewFSM_validateCallbackMap_StateOutOfRangeError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
			{IDCallbackType: CallbackBeforeEvent, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
			{IDCallbackType: CallbackLeaveState, ID: StateNonExist}: func(e *Event) { fmt.Println("Leave state opened.")},
		})
	assert.Equal(t, StateOutOfRangeError{ID: StateNonExist}, err, "Non exist state callback register NewFSM() expect StateOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist state callback register NewFSM() expect nil fsm.")

}

func TestNewFSM_validateCallbackMap_EventOutOfRangeError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
			{IDCallbackType: CallbackBeforeEvent, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventNonExist}: func(e *Event) { fmt.Println("Before event open.")},
		})
	assert.Equal(t, EventOutOfRangeError{ID: EventNonExist}, err, "Non exist event callback register NewFSM() expect EventOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist event callback register NewFSM() expect nil fsm.")

}

func TestNewFSM_validateCallbackMap_CallbackTypeOutOfRangeError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
			{IDCallbackType: CallbackTypeSum, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
		})
	assert.Equal(t,CallbackTypeOutOfRangeError{Type: CallbackTypeSum}, err, "Non exist callback type register NewFSM() expect CallbackTypeOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist callback type register NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_nonExistSrcStateError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StateNonExist}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, StateOutOfRangeError{ID: StateNonExist}, err, "Non exist src state NewFSM() expect StateOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist src state NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_nonExistDstStateError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateNonExist},
		},
		Callbacks{
		})
	assert.Equal(t, StateOutOfRangeError{ID: StateNonExist}, err, "Non exist dst state NewFSM() expect StateOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist dst state NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_nonExistEventError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventNonExist, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, EventOutOfRangeError{ID: EventNonExist}, err, "Non exist event NewFSM() expect EventOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist event NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_Src_StateStartReserveError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused, StateStartID}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, StateStartReserveError{}, err, "Reserve state NewFSM() expect StateStartReserveError.")
	assert.Nil(t, fsm, "Reserve state NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_StateStartReserveError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateStartID},
		},
		Callbacks{
		})
	assert.Equal(t, StateStartReserveError{}, err, "Reserve state NewFSM() expect StateStartReserveError.")
	assert.Nil(t, fsm, "Reserve state NewFSM() expect nil fsm.")

}

func TestNewFSM_validateEventTransitionsMap_EventStartReserveError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventStartID, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, EventStartReserveError{}, err, "Reserve event NewFSM() expect EventStartReserveError.")
	assert.Nil(t, fsm, "Reserve event NewFSM() expect nil fsm.")

}

func TestNewFSM_validateStateMap_StateStartReserveMissingError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, StateStartReserveMissingError{}, err, "Missing reserve state NewFSM() expect StateStartReserveMissingError.")
	assert.Nil(t, fsm, "Missing reserve state NewFSM() expect nil fsm.")

	fsm, err = NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: "hello",
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, StateStartReserveMissingError{}, err, "Reserve state without reserve name NewFSM() expect StateStartReserveMissingError.")
	assert.Nil(t, fsm, "Reserve state without  reserve name NewFSM() expect nil fsm.")
}

func TestNewFSM_validateEventMap_EventStartReserveMissingError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, EventStartReserveMissingError{}, err, "Missing reserve event NewFSM() expect EventStartReserveMissingError.")
	assert.Nil(t, fsm, "Missing reserve event NewFSM() expect nil fsm.")

	fsm, err = NewFSM(
		StateClosed,
		EventMap{
			EventStartID: "hello",
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
		},
		Callbacks{
		})
	assert.Equal(t, EventStartReserveMissingError{}, err, "Reserve event without reserve name NewFSM() expect EventStartReserveMissingError.")
	assert.Nil(t, fsm, "Reserve event without  reserve name NewFSM() expect nil fsm.")
}


func TestNewFSM_initStateError(t *testing.T) {
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

	fsm, err := NewFSM(
		StateNonExist,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
		},
		Callbacks{
		})
	assert.Equal(t, StateOutOfRangeError{ID: StateNonExist}, err, "Non exist init state NewFSM() expect StateOutOfRangeError.")
	assert.Nil(t, fsm, "Non exist init state NewFSM() expect nil fsm.")

}

func ExampleFSM_Transition() {
	const (
		StateStartID = iota
		StateOpened = iota
		StatePaused = iota
		StateClosed = iota
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
	)
	const (
		EventStrOpen = "open"
		EventStrPause = "paused"
		EventStrClose = "close"
	)

	fsm, err := NewFSM(
		StateClosed,
		EventMap{
			EventStartID: EventStartStr,
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
			StateStartID: StateStartStr,
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
			Events{
				{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
				{IDEvent: EventPause, IDsSrc:[]StateID{StateOpened}, IDDst:StatePaused},
				{IDEvent: EventClose, IDsSrc:[]StateID{StateOpened, StatePaused}, IDDst:StateClosed},
			},
			Callbacks{
				{IDCallbackType: CallbackEnterState, ID: StateOpened}: func(e *Event) { fmt.Println("Got into state opened.")},
				{IDCallbackType: CallbackEnterState, ID: StatePaused}: func(e *Event) { fmt.Println("Got into state paused.")},
				{IDCallbackType: CallbackEnterState, ID: StateClosed}: func(e *Event) { fmt.Println("Got into state closed.")},
			})

	if  err != nil {
		log.Fatalln("NewFSM() expect no error.")
	}

	// closed ---> opened, success
	err = fsm.Event(EventOpen)
	if  err != nil || fsm.Current() != StateStrOpened {
		log.Fatalln("Open transition from closed expect success.")
	}
	// opened ---> closed, success
	err = fsm.Event(EventClose)
	if  err != nil || fsm.Current() != StateStrClosed {
		log.Fatalln("Close transition from opened expect success.")
	}

	// opened ---> paused, success
	err = fsm.Event(EventOpen)
	if  err != nil || fsm.Current() != StateStrOpened {
		log.Fatalln("Open transition from closed expect success.")
	}
	err = fsm.Event(EventPause)
	if  err != nil || fsm.Current() != StateStrClosed {
		log.Fatalln("Pause transition from opened expect success.")
	}

	// paused --->  opened success
	err = fsm.Event(EventOpen)
	if  err != nil || fsm.Current() != StateStrOpened {
		log.Fatalln("Open transition from paused expect success.")
	}
	// paused --->  closed success
	err = fsm.Event(EventPause)
	if  err != nil || fsm.Current() != StateStrClosed {
		log.Fatalln("Pause transition from opened expect success.")
	}
	err = fsm.Event(EventClose)
	if  err != nil || fsm.Current() != StateStrClosed {
		log.Fatalln("Close transition from paused expect success.")
	}
}

// stub state callback
type stubbedCallbackState struct {
	isCalled bool
	countCalled int
	err error

	callbackType CallbackType
	callbackID CKey
	preCB Callback
	preCBM CallbackMap
}

func (cb *stubbedCallbackState) EventCall(event * Event)  {
	cb.isCalled = true
	cb.countCalled++
}

func withStubbedCallbackState(fsm *FSM, eventFunc func(eventID int, args ...interface{}) error,eventID EventID, callbackType CallbackType, callbackID CKey) (*stubbedCallbackState, error) {
	var (
		err error
		stb = &stubbedCallbackState{}
	)
	if err = stb.stub(fsm, callbackType, callbackID); err != nil {
		return nil, err
	}
	defer stb.unstub()

	stb.err = eventFunc(eventID)

	return stb, nil
}

func (cb *stubbedCallbackState) stub(fsm *FSM, callbackType CallbackType, callbackID CKey) error {
	var (
		ok bool
	)
	cb.callbackType = callbackType
	cb.callbackID = callbackID

	switch callbackType {
	case CallbackLeaveState:
		cb.preCBM = fsm.callbacksLeaveState
	case CallbackEnterState:
		cb.preCBM = fsm.callbacksEnterState
	}

	if cb.preCB, ok = cb.preCBM[cb.callbackID]; !ok {
		return errors.New("stub callback fail")
	}

	cb.preCBM[cb.callbackID] = cb
	return nil
}

func (cb *stubbedCallbackState) unstub() {
	cb.preCBM[cb.callbackID] = cb.preCB
}

// stub event callback
type stubbedCallbackEvent struct {
	isCalled bool
	countCalled int
	err error

	callbackType CallbackType
	callbackID CKey
	preCB Callback
	preCBM CallbackMap
}

func (cb *stubbedCallbackEvent) EventCall(event * Event)  {
	cb.isCalled = true
	cb.countCalled++
}

func withStubbedCallbackEvent(fsm *FSM, eventFunc func(eventID int, args ...interface{}) error,eventID EventID, callbackType CallbackType, callbackID CKey) (*stubbedCallbackEvent, error) {
	var (
		err error
		stb = &stubbedCallbackEvent{}
	)
	if err = stb.stub(fsm, callbackType, callbackID); err != nil {
		return nil, err
	}
	defer stb.unstub()

	stb.err = eventFunc(eventID)

	return stb, nil
}

func (cb *stubbedCallbackEvent) stub(fsm *FSM, callbackType CallbackType, callbackID CKey) error {
	var (
		ok bool
	)
	cb.callbackType = callbackType
	cb.callbackID = callbackID

	switch callbackType {
	case CallbackBeforeEvent:
		cb.preCBM = fsm.callbacksBeforeEvent
	case CallbackAfterEvent:
		cb.preCBM = fsm.callbacksAfterEvent
	}

	if cb.preCB, ok = cb.preCBM[cb.callbackID]; !ok {
		return errors.New("stub callback fail")
	}

	cb.preCBM[cb.callbackID] = cb
	return nil
}

func (cb *stubbedCallbackEvent) unstub() {
	cb.preCBM[cb.callbackID] = cb.preCB
}

