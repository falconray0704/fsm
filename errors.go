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
	"strconv"
)

const (
	UnknownStateStr = "unknownState"
	UnknownEventStr = "unknownEvent"
)

type DuplicateCallbackEnterStateError struct {
	state string
}

func (e DuplicateCallbackEnterStateError) Error() string {
	return "duplicate callback of enter state " + e.state
}

type DuplicateCallbackLeaveStateError struct {
	state string
}

func (e DuplicateCallbackLeaveStateError) Error() string {
	return "duplicate callback of leave state " + e.state
}

type DuplicateCallbackAfterEventError struct {
	event string
}

func (e DuplicateCallbackAfterEventError) Error() string {
	return "duplicate callback of after event " + e.event
}

type DuplicateCallbackBeforeEventError struct {
	event string
}

func (e DuplicateCallbackBeforeEventError) Error() string {
	return "duplicate callback of before event " + e.event
}

type DuplicateCallbackError struct {
}

func (e DuplicateCallbackError) Error() string {
	return "duplicate callback register"
}

type DuplicateTransitionError struct {
	event string
	state string
}

func (e DuplicateTransitionError) Error() string {
	return "duplicate transition from { " + e.event + ", " + e.state + " }"
}

type StateStartReserveError struct {
}

func (e StateStartReserveError) Error() string {
	return `State {ID:0, Name:"StateStart"} reserve for all state callback Index`
}

type StateStartReserveMissingError struct {
}

func (e StateStartReserveMissingError) Error() string {
	return `State {StateID:fsm.StateStartID, Name:fsm.StateStartStr} reserve and should be included in stateMap`
}

type EventStartReserveMissingError struct {
}

func (e EventStartReserveMissingError) Error() string {
	return `Event {EventID:fsm.EventStartID, Name:fsm.EventStartStr} reserve and should be included in eventMap`
}

type EventStartReserveError struct {
}

func (e EventStartReserveError) Error() string {
	return `Event {ID:0, Name:"EventStart"} reserve for all event callback Index`
}

type EventOutOfRangeError struct {
	ID EventID
}

func (e EventOutOfRangeError) Error() string {
	return "eventID " + strconv.Itoa(e.ID) + " out of the range"
}

type StateOutOfRangeError struct {
	ID StateID
}

func (e StateOutOfRangeError) Error() string {
	return "stateID " + strconv.Itoa(e.ID) + " out of the range"
}

type CallbackTypeOutOfRangeError struct {
	Type int
}

func (e CallbackTypeOutOfRangeError) Error() string {
	return "callback type " + strconv.Itoa(e.Type) +  " out of the range"
}

// InvalidEventError is returned by FSM.Event() when the event cannot be called
// in the current state.
type InvalidEventError struct {
	Event string
	State string
}

func (e InvalidEventError) Error() string {
	return "event " + e.Event + " inappropriate in current state " + e.State
}

// UnknownEventError is returned by FSM.Event() when the event is not defined.
type UnknownEventError struct {
	Event string
}

func (e UnknownEventError) Error() string {
	return "event " + e.Event + " does not exist"
}

// InTransitionError is returned by FSM.Event() when an asynchronous transition
// is already in progress.
type InTransitionError struct {
	Event string
}

func (e InTransitionError) Error() string {
	return "event " + e.Event + " inappropriate because previous transition did not complete"
}

// NotInTransitionError is returned by FSM.Transition() when an asynchronous
// transition is not in progress.
type NotInTransitionError struct{}

func (e NotInTransitionError) Error() string {
	return "transition inappropriate because no state change in progress"
}

// NoTransitionError is returned by FSM.Event() when no transition have happened,
// for example if the source and destination states are the same.
type NoTransitionError struct {
	Err error
}

func (e NoTransitionError) Error() string {
	if e.Err != nil {
		return "no transition with error: " + e.Err.Error()
	}
	return "no transition"
}

// CanceledError is returned by FSM.Event() when a callback have canceled a
// transition.
type CanceledError struct {
	Err error
}

func (e CanceledError) Error() string {
	if e.Err != nil {
		return "transition canceled with error: " + e.Err.Error()
	}
	return "transition canceled"
}

// AsyncError is returned by FSM.Event() when a callback have initiated an
// asynchronous state transition.
type AsyncError struct {
	Err error
}

func (e AsyncError) Error() string {
	if e.Err != nil {
		return "async started with error: " + e.Err.Error()
	}
	return "async started"
}

// InternalError is returned by FSM.Event() and should never occur. It is a
// probably because of a bug.
type InternalError struct{}

func (e InternalError) Error() string {
	return "internal error on state transition"
}
