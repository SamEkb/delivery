package order

type Status string

const (
	Empty     Status = ""
	Created   Status = "Created"
	Assigned  Status = "Assigned"
	Completed Status = "Completed"
)

func (s Status) equals(other Status) bool {
	return s == other
}

func (s Status) isEmpty() bool {
	return s == Empty
}

func (s Status) String() string {
	return string(s)
}
