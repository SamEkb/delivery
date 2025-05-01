package order

type Status string

const (
	empty     Status = ""
	created   Status = "Created"
	assigned  Status = "Assigned"
	completed Status = "Completed"
)

func (s Status) equals(other Status) bool {
	return s == other
}

func (s Status) isEmpty() bool {
	return s == empty
}

func (s Status) String() string {
	return string(s)
}
