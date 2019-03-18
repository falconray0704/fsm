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

// Package fsm implements a finite state machine.
//
// It is heavily based on two FSM implementations:
//
// Javascript Finite State Machine
// https://github.com/jakesgordon/javascript-state-machine
//
// Fysom for Python
// https://github.com/oxplot/fysom (forked at https://github.com/mriehl/fysom)
//
package fsm

import (
	"sync"
)

// transitioner is an interface for the FSM's transition function.
type transitioner interface {
	//transition(*FSM) error
	transition() error
}

type StateID = int
type EventID = int

type EventMap map[EventID]string
type StateMap map[StateID]string

type CallbackMap map[cKey]Callback

// FSM is the state machine that holds the current state.
//
// It has to be created with NewFSM to function properly.
type FSM struct {
	// current is the state that the FSM is currently in.
	//current string
	current StateID

	// transitions maps events and source states to destination states.
	transitions map[eKey]StateID

	// callbacks maps events and targers to callback functions.
	//callbacks map[cKey]Callback

	// transition is the internal transition functions used either directly
	// or when Transition is called in an asynchronous state transition.
	transition func()
	// transitionerObj calls the FSM's transition() function.
	// transitionRunner transitioner

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex


	// map of events callbacks
	eventMap	EventMap
	stateMap	StateMap
	callbacksBeforeEvent	CallbackMap
	callbacksLeaveState		CallbackMap
	callbacksEnterState		CallbackMap
	callbacksAfterEvent		CallbackMap
}

type CallBackDesc struct {
	IDCallbackType	CallbackType
	ID			int // ID of event or state
}

// EventDesc represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type EventDesc struct {
	// Name is the event name used when calling for a transition.
	IDEvent EventID

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	IDsSrc []StateID

	// Dst is the destination state that the FSM will be in if the transition
	// succeds.
	IDDst StateID
}

// Callback is a function type that callbacks should use. Event is the current
// event info as the callback happens.
type Callback func(*Event)

// Events is a shorthand for defining the transition map in NewFSM.
type Events []EventDesc

// Callbacks is a shorthand for defining the callbacks in NewFSM.a
type Callbacks map[CallBackDesc]Callback

// NewFSM constructs a FSM from events and callbacks.
//
// The events and transitions are specified as a slice of Event structs
// specified as Events. Each Event is mapped to one or more internal
// transitions from Event.Src to Event.Dst.
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
//
// There are also two short form versions for the most commonly used callbacks.
// They are simply the name of the event or state:
//
// 1. <NEW_STATE> - called after entering <NEW_STATE>
//
// 2. <EVENT> - called after event named <EVENT>
//
// If both a shorthand version and a full version is specified it is undefined
// which version of the callback will end up in the internal map. This is due
// to the psuedo random nature of Go maps. No checking for multiple keys is
// currently performed.
func NewFSM(initState int, eventMap EventMap, stateMap StateMap,
			eventTransitions []EventDesc,
			callbackDescMap map[CallBackDesc]Callback) (*FSM, error) {
	var (
		ok bool
		err error
	)

	if  _, ok = stateMap[initState]; !ok {
		return nil, StateOutOfRangeError{ID: initState}
	}

	if err = validateEventTransitionsMap(eventTransitions, eventMap, stateMap); err != nil {
		return nil, err
	}

	if err = validateCallbackMap(callbackDescMap, eventMap, stateMap); err != nil {
		return nil, err
	}


	f := &FSM{
		//transitionRunner: &transitionerStruct{},
		eventMap:		eventMap,
		stateMap:		stateMap,
		current:		initState,
		transitions:	make(map[eKey]StateID),
		callbacksBeforeEvent:		make(CallbackMap),
		callbacksLeaveState:		make(CallbackMap),
		callbacksEnterState:		make(CallbackMap),
		callbacksAfterEvent:		make(CallbackMap),
		//transitions:     make(map[eKey]string),
	}

	if err = f.buildUpTransitions(eventTransitions); err != nil {
		return nil, err
	}

	f.buildUpCallbackMap(callbackDescMap)

	return f, nil
}

func (f *FSM) buildUpCallbackMap(callbackDescMap map[CallBackDesc]Callback) {
	var (
		desc CallBackDesc
		cb Callback
	)

	for desc, cb = range callbackDescMap {
		switch desc.IDCallbackType {
		case CallbackBeforeEvent:
			f.callbacksBeforeEvent[desc.ID] = cb
		case CallbackLeaveState:
			f.callbacksLeaveState[desc.ID] = cb
		case CallbackEnterState:
			f.callbacksEnterState[desc.ID] = cb
		case CallbackAfterEvent:
			f.callbacksAfterEvent[desc.ID] = cb
		}
	}
}

/*
func registerCallback(callbackMap CallbackMap, id int, callback Callback) error {
	var (
		ok bool
	)

	if _, ok = callbackMap[id]; ok {
		callbackMap[id] = callback
		return DuplicateCallbackError{}
	}

	callbackMap[id] = callback
	return nil
}
*/

func (f *FSM) buildUpTransitions(eventTransitions []EventDesc) error {
	var (
		ok bool
		ed EventDesc
	)
	for _, ed = range eventTransitions {
		for _, src := range ed.IDsSrc {
			ek := eKey{event: ed.IDEvent, src: src }
			if _, ok = f.transitions[ek]; ok {
				f.transitions[ek] = ed.IDDst
				return DuplicateTransitionError{event: f.eventMap[ek.event], state: f.stateMap[ek.src]}
			}
			f.transitions[ek] = ed.IDDst
		}
	}
	return nil
}

func validateEventTransitionsMap(eventTransitions []EventDesc, eventMap EventMap, stateMap StateMap) error {
	var (
		ok bool
		ed EventDesc
	)
	for _, ed = range eventTransitions {
		if _, ok = eventMap[ed.IDEvent]; !ok {
			return EventOutOfRangeError{ID: ed.IDEvent}
		}
		if _, ok = stateMap[ed.IDDst]; !ok {
			return StateOutOfRangeError{ID: ed.IDDst}
		}
		for _, srcState := range ed.IDsSrc {
			if _, ok = stateMap[srcState]; !ok {
				return StateOutOfRangeError{ID: srcState}
			}
		}
	}

	return nil
}

func validateCallbackMap(callbackDescMap map[CallBackDesc]Callback, eventMap EventMap, stateMap StateMap) error {
	var (
		ok bool
	)
	for desc := range callbackDescMap {
		switch desc.IDCallbackType {
		case CallbackBeforeEvent, CallbackAfterEvent:
			if _, ok = eventMap[desc.ID]; !ok && desc.ID != CKeyAll {
				return EventOutOfRangeError{ID: desc.ID}
			}
		case CallbackLeaveState, CallbackEnterState:
			if _, ok = stateMap[desc.ID]; !ok && desc.ID != CKeyAll {
				return StateOutOfRangeError{ID: desc.ID}
			}
		default:
			return CallbackTypeOutOfRangeError{Type: desc.IDCallbackType}
		}
	}
	return nil
}

// Current returns the current state of the FSM.
func (f *FSM) Current() string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.stateMap[f.current]
}

// Is returns true if state is the current state.
func (f *FSM) Is(stateID int) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return stateID == f.current
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *FSM) SetState(stateID int) error {
	if _, ok :=  f.stateMap[stateID]; !ok {
		return StateOutOfRangeError{ID: stateID}
	}

	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = stateID
	return nil
}

// Can returns true if event can occur in the current state.
func (f *FSM) Can(eventID int) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	_, ok := f.transitions[eKey{eventID, f.current}]
	return ok && (f.transition == nil)
}

// AvailableTransitions returns a list of transitions avilable in the
// current state.
func (f *FSM) AvailableTransitions() []string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	var transitions []string
	for key := range f.transitions {
		if key.src == f.current {
			transitions = append(transitions, f.eventMap[key.event])
		}
	}
	return transitions
}

// Cannot returns true if event can not occure in the current state.
// It is a convenience method to help code read nicely.
func (f *FSM) Cannot(eventID int) bool {
	return !f.Can(eventID)
}

// Event initiates a state transition with the named event.
//
// The call takes a variable number of arguments that will be passed to the
// callback, if defined.
//
// It will return nil if the state change is ok or one of these errors:
//
// - event X inappropriate because previous transition did not complete
//
// - event X inappropriate in current state Y
//
// - event X does not exist
//
// - internal error on state transition
//
// The last error should never occur in this situation and is a sign of an
// internal bug.
func (f *FSM) Event(eventID int, args ...interface{}) error {
	var (
		ok bool
		eventName string
	)

	if eventName, ok = f.eventMap[eventID]; !ok {
		return EventOutOfRangeError{}
	}

	f.eventMu.Lock()
	defer f.eventMu.Unlock()

	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	if f.transition != nil {
		return InTransitionError{eventName}
	}

	dst, ok := f.transitions[eKey{eventID, f.current}]
	if !ok {
		return InvalidEventError{eventName, f.stateMap[f.current]}
	}

	e := &Event{eventID, f.current, dst, nil, args, false, false}

	err := f.beforeEventCallbacks(e)
	if err != nil {
		return err
	}

	if f.current == dst {
		f.afterEventCallbacks(e)
		return NoTransitionError{e.Err}
	}

	// Setup the transition, call it later.
	f.transition = func() {
		f.stateMu.Lock()
		f.current = dst
		f.stateMu.Unlock()

		f.enterStateCallbacks(e)
		f.afterEventCallbacks(e)
	}

	if err = f.leaveStateCallbacks(e); err != nil {
		if _, ok := err.(CanceledError); ok {
			f.transition = nil
		}
		return err
	}

	// Perform the rest of the transition, if not asynchronous.
	f.stateMu.RUnlock()
	err = f.doTransition()
	f.stateMu.RLock()
	if err != nil {
		return InternalError{}
	}

	return e.Err
}

// Transition wraps transitioner.transition.
func (f *FSM) Transition() error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()
	return f.doTransition()
}

// doTransition wraps transitioner.transition.
func (f *FSM) doTransition() error {

	//return f.transitionRunner.transition(f)

	if f.transition == nil {
		return NotInTransitionError{}
	}
	f.transition()
	f.transition = nil
	return nil
}

// beforeEventCallbacks calls the before_ callbacks, first the named then the
// general version.
func (f *FSM) beforeEventCallbacks(e *Event) error {
	if fn, ok := f.callbacksBeforeEvent[cKey(e.Event + 1)]; ok && fn != nil {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}

	if fn, ok := f.callbacksBeforeEvent[CKeyAll]; ok && fn != nil {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the
// general version.
func (f *FSM) leaveStateCallbacks(e *Event) error {
	if fn, ok := f.callbacksLeaveState[cKey(f.current + 1)]; ok && fn != nil {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}

	if fn, ok := f.callbacksLeaveState[CKeyAll]; ok && fn != nil {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}
	return nil
}

// enterStateCallbacks calls the enter_ callbacks, first the named then the
// general version.
func (f *FSM) enterStateCallbacks(e *Event) {
	if fn, ok := f.callbacksEnterState[cKey(f.current + 1)]; ok && fn != nil {
		fn(e)
	}
	if fn, ok := f.callbacksEnterState[CKeyAll]; ok && fn != nil {
		fn(e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the
// general version.
func (f *FSM) afterEventCallbacks(e *Event) {
	if fn, ok := f.callbacksAfterEvent[cKey(e.Event + 1)]; ok && fn != nil {
		fn(e)
	}
	if fn, ok := f.callbacksAfterEvent[CKeyAll]; ok && fn != nil {
		fn(e)
	}
}

const (
	CallbackNone			CallbackType = iota		// 0
	CallbackBeforeEvent		CallbackType = iota		// 1
	CallbackLeaveState		CallbackType = iota		// 2
	CallbackEnterState		CallbackType = iota		// 3
	CallbackAfterEvent		CallbackType = iota		// 4
	CallbackTypeSum			CallbackType = iota		// 5, total of callback types
)

type CallbackType = int

// cKey is a callback key used for indexing the callback corresponding
// to event or state, must start from 1, 0 use for all event and state callback
type cKey = int
const (
	CKeyAll cKey = 0 // Index 0 reserve for all event or status callback
)

// eKey is a struct key used for storing the transition map.
type eKey struct {
	// event is the key of the event
	event EventID

	// src is the source state from where the event can transition.
	src StateID
}

