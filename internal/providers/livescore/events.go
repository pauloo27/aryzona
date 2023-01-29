package livescore

type EventType int

const (
	EventTypeGoal             = 36
	EventTypeFoulPenaltyGoal  = 37
	EventTypeOvertimeGoal     = 47
	EventTypeYellowCard       = 43
	EventTypeRedCard          = 44
	EventTypeDoubleYellowCard = 45
	EventTypePenaltyGoal      = 41
	EventTypePenaltyMissed    = 40
)
