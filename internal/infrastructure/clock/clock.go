// Package clock — реальное время для порта commands.Clock.
package clock

import (
	"time"

	"github.com/Game-WW/users-be/internal/core/application/commands"
)

type Clock struct{}

var _ commands.Clock = Clock{}

func (Clock) Now() time.Time {
	return time.Now()
}
