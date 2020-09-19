package eventlog

type destinationFace interface {
	Write(*EventLog) (*EventLog, error)
}
