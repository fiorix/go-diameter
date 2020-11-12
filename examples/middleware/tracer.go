package middleware

import (
	"strings"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/opentracing/opentracing-go"
)

type Tracer struct {
	h diam.Handler
}

func getOp(m *diam.Message) string {
	op := "diameter_"
	cmd, err := m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode,
	)
	if err != nil {
		op += "unknown_command"
	} else {
		op += strings.ToLower(strings.Replace(cmd.Name, "-", "_", 1))
	}
	if m.Header.CommandFlags&diam.RequestFlag == diam.RequestFlag {
		op += "r"
	} else {
		op += "a"
	}
	return op
}

func (t *Tracer) ServeDIAM(c diam.Conn, m *diam.Message) {
	span, ctx := opentracing.StartSpanFromContext(m.Context(), getOp(m))
	defer span.Finish()
	m.SetContext(ctx)
	t.h.ServeDIAM(c, m)
}

func NewTracer(h diam.Handler) diam.Handler {
	return &Tracer{
		h: h,
	}
}

func TracerFunc(f diam.HandlerFunc) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		span, ctx := opentracing.StartSpanFromContext(m.Context(), getOp(m))
		defer span.Finish()
		m.SetContext(ctx)
		f(c, m)
	}
}
