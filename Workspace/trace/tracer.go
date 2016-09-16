package trace

//Interface for tracing code events

type Tracer interface {
	Trace(...interface{})
}