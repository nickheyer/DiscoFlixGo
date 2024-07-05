package tasks

import (
	"github.com/nickheyer/DiscoFlixGo/pkg/services"
)

// Register registers all task queues with the task client
func Register(c *services.Container) {
	c.Tasks.Register(NewExampleTaskQueue(c))
}
