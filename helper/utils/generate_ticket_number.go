package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GenerateTicketNumber struct {
}

func NewGenerateTicketNumber() *GenerateTicketNumber {
	return &GenerateTicketNumber{}
}

func (s *GenerateTicketNumber) UUID() string {
	return fmt.Sprintf("TCK-%s-%s", time.Now().Format("20060102"), uuid.New().String()[:8])
}
