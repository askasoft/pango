package log

type BridgeWriter struct {
	Logger Logger
}

func (bw *BridgeWriter) Write(le *Event) error {
	bw.Logger.Write(*le)
	return nil
}

func (bw *BridgeWriter) Flush() {
}

func (bw *BridgeWriter) Close() {
}

func NewBridgeWriter(logger Logger) Writer {
	return &BridgeWriter{logger}
}
