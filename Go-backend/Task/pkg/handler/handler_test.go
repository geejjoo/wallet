package handler_test

import (
	"reflect"
	"testing"

	"github.com/geejjoo/task/pkg/handler"
	"github.com/geejjoo/task/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestHandler_InitRoutes(t *testing.T) {
	svc := &service.Service{}
	h := handler.NewHandler(svc)
	exp := reflect.TypeOf(h).String()

	assert.Equal(t, exp, "*handler.Handler")
}
