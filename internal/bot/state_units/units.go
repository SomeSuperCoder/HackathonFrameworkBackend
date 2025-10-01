package stateunits

import statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"

const (
	STATE_NONE statemachine.StateUnit = iota
	STATE_ENTER_NAME
	STATE_ENTER_BIRTHDATE
)
