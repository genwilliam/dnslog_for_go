package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateTraceID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.NewString())
}
