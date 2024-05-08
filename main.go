package sentryLoggerGo

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"reflect"
	"time"
)

type SentryLogMsg struct {
	ErrorMsg    string                   `json:"error_msg,omitempty"`
	Error       error                    `json:"error,omitempty"`
	ApiEndPoint string                   `json:"api_end_point"`
	RabbitMq    string                   `json:"rabbit_mq"`
	StackTrace  []map[string]interface{} `json:"stack_trace"`
}

func SentryInit(sentryDsn string) {
	if sentryDsn == "" {
		fmt.Println("sentry connection error | sentryDsn not found")
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryDsn,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return event
		},
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	} else {
		log.Println("SentryInit success")
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)
}

func (sentryLogMsg SentryLogMsg) captureError() {

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.AddEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			event.Exception[0].Type = sentryLogMsg.ErrorMsg
			return event
		})
		scope.SetExtra("Api Endpoint", sentryLogMsg.ApiEndPoint)
		scope.SetLevel(sentry.LevelError)
		sentry.CaptureException(sentryLogMsg.Error)
	})
}

func SentryMiddleware(c *gin.Context) {
	defer func() {
		var sentryLogMsg = SentryLogMsg{}
		if err := recover(); err != nil {

			//API Response
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusOK, gin.H{
					"message":     "Document Not Found",
					"data":        nil,
					"status_code": http.StatusNotFound,
				})
			} else {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"message":     err,
					"data":        nil,
					"status_code": http.StatusUnprocessableEntity,
				})
			}

			//Error capture
			sentryLogMsg.ApiEndPoint = c.Request.RequestURI
			if e, ok := err.(error); ok {
				sentryLogMsg.Error = e
			} else {
				errMsg := fmt.Sprintf("%v", err)
				sentryLogMsg.Error = errors.New(errMsg)
			}
			sentryLogMsg.ErrorMsg = reflect.TypeOf(sentryLogMsg.Error).String()
			sentryLogMsg.captureError()
		}
	}()
	c.Next()
}

func TestErrorLog(errorMsg string) {
	if errorMsg == "" {
		errorMsg = "Test Error Log"
	}
	sentry.CaptureException(errors.New(errorMsg))
}

func main() {
	fmt.Println("Sentry logger calling")
}
