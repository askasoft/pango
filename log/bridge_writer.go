package log

type BridgeWriter struct {
	Logger Logger
}

func (bw *BridgeWriter) Write(le *Event) {
	bw.Logger.Write(le)
}

func (bw *BridgeWriter) Flush() {
}

func (bw *BridgeWriter) Close() {
}

func NewBridgeWriter(logger Logger) *BridgeWriter {
	return &BridgeWriter{logger}
}
