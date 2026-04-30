package events

type UserRegisteredEvent struct {
	UserID    int32
	Email     string
	FirstName string
	LastName  string
	Role      string
}

type UserRegisteredEventHandler struct {
	
}

