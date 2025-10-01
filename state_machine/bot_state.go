package statemachine

type StateKey int64
type StateUnit int

type BotState struct {
	stateMap map[StateKey]UserState
}

func (s *BotState) SetState(key StateKey, state StateUnit, value any) {
	currentState := s.stateMap[key]
	currentState.State = state
	currentState.Data = value
	s.stateMap[key] = currentState
}

func (s *BotState) UpdateValue(key StateKey, updater func(value any) any) {
	currentState := s.GetState(key)
	s.SetState(key, currentState.State, updater(currentState.Data))
}

func (s *BotState) GetState(key StateKey) UserState {
	return s.stateMap[key]
}

func NewBotState() *BotState {
	return &BotState{
		stateMap: make(map[StateKey]UserState),
	}
}
