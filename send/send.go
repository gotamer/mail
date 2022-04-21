package send

const dirQueue = "/var/spool/queue/"

// Sender is an interface with a Send method, that dispatches a single email envelop
type Sender interface {
	Send() error
}
