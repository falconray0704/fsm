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

import "sync"

type FSM interface {
	Event(eventID int, args ...interface{}) error
	Transition() error
	Current() (int, string)
	Is(stateID int) bool
	SetState(stateID int) error
	Can(eventID int) bool
	Cannot(eventID int) bool
	AvailableTransitions() []string
}

// transitioner is an interface for the FSM's transition function.
type transitioner interface {
	transition(*fsmGo) error
}

type CallbackType = int
const (
	CallbackTypeNone		CallbackType = iota		// 0
	CallbackBeforeEvent		CallbackType = iota		// 1
	CallbackLeaveState		CallbackType = iota		// 2
	CallbackEnterState		CallbackType = iota		// 3
	CallbackAfterEvent		CallbackType = iota		// 4
	CallbackTypeSum			CallbackType = iota		// 5, total of callback types
)

// CBKey is a callback key used for indexing the callback corresponding
// to event ID or state ID defined , must start from 1.
// 0, reserve for all event id and state id callback.
type CBKey = int
const (
	EventStartID CBKey = 0 // Event ID start Index, reserve for all event callback index
	EventStartStr = "EventStart"
	StateStartID CBKey = 0 // State ID start Index, reserve for all state callback index
	StateStartStr = "StateStart"
)

// eKey is a struct key used for storing the transition map.
type eKey struct {
	// event is the key of the event
	event EventID
	// src is the source state from where the event can transition.
	src StateID
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

	// Dst is the destination state that the FSM will be in if the transition success.
	IDDst StateID
}

// Events is a shorthand for defining the transition map in NewFSM.
type Events []EventDesc


// Callback is a interface that internal use to call user's callback function.
// Event is the current event info as the callback happens.
type Callback interface {
	EventCall(event *Event)
}

// CallbackFunc is the event/state callback function type that user implemented.
type CallbackFunc func(*Event)

func (f CallbackFunc) EventCall(event *Event) {
	f(event)
}

// CallBackDesc represents an callback defining for happening of event/state when initializing the FSM.
type CallBackDesc struct {
	IDCallbackType	CallbackType	// type of callback, only supports:
	// CallbackBeforeEvent, CallbackLeaveState, CallbackEnterState, CallbackAfterEvent

	ID				int 			// ID of event or state
}

// Callbacks is a shorthand for defining the callbacks in NewFSM.a
type Callbacks map[CallBackDesc]CallbackFunc

type StateID = int
type EventID = int

type EventMap map[EventID]string
type StateMap map[StateID]string

type CallbackMap map[CBKey]Callback


// FSM is the state machine that holds the current state.
//
// It has to be created with NewFSM to function properly.
type fsmGo struct {
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
	transitionerObj transitioner

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex


	// map of event ID to event name string
	eventMap	EventMap
	// map of state ID to state name string
	stateMap	StateMap
	// map of event ID to before event callback
	callbacksBeforeEvent	CallbackMap
	// map of state ID to leave state callback
	callbacksLeaveState		CallbackMap
	// map of state ID to enter state callback
	callbacksEnterState		CallbackMap
	// map of event ID to after event callback
	callbacksAfterEvent		CallbackMap
}

// NewFSM constructs a FSM from maps of events, states, and callbacks, ane set of transitions.
//
// The eventMap and stateMap are specified as a map of event/state ID to event/state name string.
// Each Event "EventDesc.IDEvent" is mapped to one or more internal transitions from "EventDesc.IDsSrc" to "EventDesc.IDDst".
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
// 1. <StateStartID> - called after entering <StateStartID>
//
// 2. <EventStartID> - called after event named <EventStartID>
//
func NewFSM(initState StateID, eventMap EventMap, stateMap StateMap,
	eventTransitions []EventDesc,
	callbackDescMap map[CallBackDesc]CallbackFunc) (FSM, error) {
		return newGoFSM(initState, eventMap, stateMap, eventTransitions, callbackDescMap)
}

func newGoFSM(initState StateID, eventMap EventMap, stateMap StateMap,
			eventTransitions []EventDesc,
			callbackDescMap map[CallBackDesc]CallbackFunc) (*fsmGo, error) {
	var (
		ok bool
		err error
	)

	if err = validateEventMap(eventMap); err != nil {
		return nil, err
	}

	if err = validateStateMap(stateMap); err != nil {
		return nil, err
	}

	if initState == StateStartID {
		return nil, StateStartReserveError{}
	} else if  _, ok = stateMap[initState]; !ok {
		return nil, StateOutOfRangeError{ID: initState}
	}

	if err = validateEventTransitionsMap(eventTransitions, eventMap, stateMap); err != nil {
		return nil, err
	}

	if err = validateCallbackMap(callbackDescMap, eventMap, stateMap); err != nil {
		return nil, err
	}


	f := &fsmGo{
		eventMap:		eventMap,
		stateMap:		stateMap,
		current:		initState,
		transitions:	make(map[eKey]StateID),
		transitionerObj: &transitionerStruct{},
		callbacksBeforeEvent:		make(CallbackMap),
		callbacksLeaveState:		make(CallbackMap),
		callbacksEnterState:		make(CallbackMap),
		callbacksAfterEvent:		make(CallbackMap),
	}

	if err = f.buildUpTransitions(eventTransitions); err != nil {
		return nil, err
	}

	f.buildUpCallbackMap(callbackDescMap)

	return f, nil
}

func (f *fsmGo) buildUpCallbackMap(callbackDescMap map[CallBackDesc]CallbackFunc) {
	var (
		desc CallBackDesc
		cbf CallbackFunc
	)

	for desc, cbf = range callbackDescMap {
		switch desc.IDCallbackType {
		case CallbackBeforeEvent:
			f.callbacksBeforeEvent[desc.ID] = CallbackFunc(cbf)
		case CallbackLeaveState:
			f.callbacksLeaveState[desc.ID] = CallbackFunc(cbf)
		case CallbackEnterState:
			f.callbacksEnterState[desc.ID] = CallbackFunc(cbf)
		case CallbackAfterEvent:
			f.callbacksAfterEvent[desc.ID] = CallbackFunc(cbf)
		}
	}
}

func (f *fsmGo) buildUpTransitions(eventTransitions []EventDesc) error {
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

func validateStateMap(stateMap StateMap) error {
	name, ok := stateMap[StateStartID]
	if !ok || name != StateStartStr {
		return StateStartReserveMissingError{}
	}
	return nil
}

func validateEventMap(eventMap EventMap) error {
	name, ok := eventMap[EventStartID]
	if !ok || name != EventStartStr {
		return EventStartReserveMissingError{}
	}
	return nil
}

func validateEventTransitionsMap(eventTransitions []EventDesc, eventMap EventMap, stateMap StateMap) error {
	var (
		ok bool
		ed EventDesc
	)
	for _, ed = range eventTransitions {
		if ed.IDEvent == EventStartID {
			return EventStartReserveError{}
		} else if ed.IDDst == StateStartID {
			return StateStartReserveError{}
		} else if _, ok = eventMap[ed.IDEvent]; !ok {
			return EventOutOfRangeError{ID: ed.IDEvent}
		} else if _, ok = stateMap[ed.IDDst]; !ok {
			return StateOutOfRangeError{ID: ed.IDDst}
		}

		for _, srcState := range ed.IDsSrc {
			if srcState == StateStartID {
				return StateStartReserveError{}
			} else if _, ok = stateMap[srcState]; !ok {
				return StateOutOfRangeError{ID: srcState}
			}
		}
	}

	return nil
}

func validateCallbackMap(callbackDescMap map[CallBackDesc]CallbackFunc, eventMap EventMap, stateMap StateMap) error {
	var (
		ok bool
	)
	for desc := range callbackDescMap {
		switch desc.IDCallbackType {
		case CallbackBeforeEvent, CallbackAfterEvent:
			if _, ok = eventMap[desc.ID]; !ok {
				return EventOutOfRangeError{ID: desc.ID}
			}
		case CallbackLeaveState, CallbackEnterState:
			if _, ok = stateMap[desc.ID]; !ok {
				return StateOutOfRangeError{ID: desc.ID}
			}
		default:
			return CallbackTypeOutOfRangeError{Type: desc.IDCallbackType}
		}
	}
	return nil
}

// Current returns the current state of the FSM.
func (f *fsmGo) Current() (int, string) {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.current, f.stateMap[f.current]
}

// Is returns true if state is the current state.
func (f *fsmGo) Is(stateID int) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return stateID == f.current
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *fsmGo) SetState(stateID int) error {

	if stateID == StateStartID {
		return StateStartReserveError{}
	}

	if _, ok :=  f.stateMap[stateID]; !ok {
		return StateOutOfRangeError{ID: stateID}
	}

	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = stateID
	return nil
}

// Can returns true if event can occur in the current state.
func (f *fsmGo) Can(eventID int) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	_, ok := f.transitions[eKey{eventID, f.current}]
	return ok && (f.transition == nil)
}

// AvailableTransitions returns a list of transitions avilable in the
// current state.
func (f *fsmGo) AvailableTransitions() []string {
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
func (f *fsmGo) Cannot(eventID int) bool {
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
func (f *fsmGo) Event(eventID int, args ...interface{}) error {
	var (
		ok bool
		eventName string
	)

	if eventName, ok = f.eventMap[eventID]; !ok {
		return EventOutOfRangeError{ID: eventID}
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
func (f *fsmGo) Transition() error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()
	return f.doTransition()
}

// doTransition wraps transitioner.transition.
func (f *fsmGo) doTransition() error {
	return f.transitionerObj.transition(f)
}

// transitionerStruct is the default implementation of the transitioner
// interface. Other implementations can be swapped in for testing.
type transitionerStruct struct{}

// Transition completes an asynchrounous state change.
//
// The callback for leave_<STATE> must prviously have called Async on its
// event to have initiated an asynchronous state transition.
func (t transitionerStruct) transition(f *fsmGo) error {
	if f.transition == nil {
		return NotInTransitionError{}
	}
	f.transition()
	f.transition = nil
	return nil
}


// beforeEventCallbacks calls the before_ callbacks, first the named then the
// general version.
func (f *fsmGo) beforeEventCallbacks(e *Event) error {
	if fn, ok := f.callbacksBeforeEvent[CBKey(e.Event)]; ok && fn != nil {
		fn.EventCall(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}

	if fn, ok := f.callbacksBeforeEvent[StateStartID]; ok && fn != nil {
		fn.EventCall(e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the
// general version.
func (f *fsmGo) leaveStateCallbacks(e *Event) error {
	if fn, ok := f.callbacksLeaveState[CBKey(f.current)]; ok && fn != nil {
		fn.EventCall(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}

	if fn, ok := f.callbacksLeaveState[StateStartID]; ok && fn != nil {
		fn.EventCall(e)
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
func (f *fsmGo) enterStateCallbacks(e *Event) {
	if fn, ok := f.callbacksEnterState[CBKey(f.current)]; ok && fn != nil {
		fn.EventCall(e)
	}
	if fn, ok := f.callbacksEnterState[StateStartID]; ok && fn != nil {
		fn.EventCall(e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the
// general version.
func (f *fsmGo) afterEventCallbacks(e *Event) {
	if fn, ok := f.callbacksAfterEvent[CBKey(e.Event)]; ok && fn != nil {
		fn.EventCall(e)
	}
	if fn, ok := f.callbacksAfterEvent[EventStartID]; ok && fn != nil {
		fn.EventCall(e)
	}
}


