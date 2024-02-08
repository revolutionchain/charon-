package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/notifier"
)

func getNotifier(c echo.Context) *notifier.Notifier {
	storedValue := c.Get("notifier")
	notifier, ok := storedValue.(*notifier.Notifier)
	if !ok {
		return nil
	}
	return notifier
}
