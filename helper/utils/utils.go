package utils

type Utils struct {
	JaegerTracer         JaegerTracer
	GenerateTicketNumber GenerateTicketNumber
}

func NewUtils() *Utils {
	return &Utils{
		JaegerTracer:         *NewJaegerTracer(),
		GenerateTicketNumber: *NewGenerateTicketNumber(),
	}
}
