package mail

type ReminderFields struct {
	Organization string
	Service      string
	SessionTime  string
	SessionDate  string
}

type Message struct {
	Email string
	Type  string
	Value interface{}
}
