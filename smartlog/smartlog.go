package smartlog

import (
	"context"
	"fmt"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

func init() {
	server.RegisterExtension(&myLogExt{})
}

type myLogExt struct {

}

func (e *myLogExt) Name() string {
	return "github.com/apexcz/fnExtensions/smartlog"
}

func (e *myLogExt) Setup(s fnext.ExtServer) error {
	s.AddCallListener(&LogListener{})
	return nil
}

type LogListener struct {

}

func (l *LogListener) BeforeCall(ctx context.Context, call *models.Call) error {
	fmt.Println("This is an interception that occurs before a function is called")
	return nil
}

func (l *LogListener) AfterCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Hahahaha YO! And this is an annoying message that will happen AFTER every time a function is called.")
	return nil
}