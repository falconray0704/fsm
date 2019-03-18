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

package fsm_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
	. "github.com/falconray0704/fsm"
)





func TestNewFSM(t *testing.T) {
	const (
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
			StateOpened: StateStrOpened,
			StatePaused: StateStrPaused,
			StateClosed: StateStrClosed },
		Events{
			{IDEvent: EventOpen, IDsSrc:[]StateID{StateClosed, StatePaused}, IDDst:StateOpened},
			{IDEvent: EventPause, IDsSrc:[]StateID{StateOpened}, IDDst:StatePaused},
			{IDEvent: EventClose, IDsSrc:[]StateID{StateOpened, StatePaused}, IDDst:StateClosed},
		},
		Callbacks{
			{IDCallbackType: CallbackBeforeEvent, ID: EventOpen}: func(e *Event) { fmt.Println("Before event open.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventPause}: func(e *Event) { fmt.Println("Before event pause.")},
			{IDCallbackType: CallbackBeforeEvent, ID: EventClose}: func(e *Event) { fmt.Println("Before event close.")},
			{IDCallbackType: CallbackLeaveState, ID: StateOpened}: func(e *Event) { fmt.Println("Leave state opened.")},
			{IDCallbackType: CallbackLeaveState, ID: StatePaused}: func(e *Event) { fmt.Println("Leave state paused.")},
			{IDCallbackType: CallbackLeaveState, ID: StateClosed}: func(e *Event) { fmt.Println("Leave state closed.")},
			{IDCallbackType: CallbackEnterState, ID: StateOpened}: func(e *Event) { fmt.Println("Got into state opened.")},
			{IDCallbackType: CallbackEnterState, ID: StatePaused}: func(e *Event) { fmt.Println("Got into state paused.")},
			{IDCallbackType: CallbackEnterState, ID: StateClosed}: func(e *Event) { fmt.Println("Got into state closed.")},
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


	// closed ---> opened, success
	err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")
	// opened ---> closed, success
	err = fsm.Event(EventClose)
	assert.NoError(t, err, "Close transition from opened expect success.")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Close transition expect closed")

	// opened ---> paused, success
	err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")
	err = fsm.Event(EventPause)
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
	err = fsm.Event(EventOpen)
	assert.NoError(t, err, "Open transition from paused expect success.")
	assert.Equal(t, StateStrOpened, fsm.Current(), "Open transition expect opened")

	// paused --->  closed success
	err = fsm.Event(EventPause)
	assert.NoError(t, err, "Open transition from closed expect success.")
	assert.Equal(t, StateStrPaused, fsm.Current(), "Pause transition expect paused")
	err = fsm.Event(EventClose)
	assert.NoError(t, err, "Close transition from paused expect success.")
	assert.Equal(t, StateStrClosed, fsm.Current(), "Close transition expect closed")


}

func ExampleFSM_Transition() {
	const (
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
			EventOpen: EventStrOpen,
			EventPause: EventStrPause,
			EventClose: EventStrClose },
		StateMap{
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

