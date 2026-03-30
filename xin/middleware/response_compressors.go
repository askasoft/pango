package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"sync"
)

type Compressor interface {
	io.Writer
	Reset(w io.Writer)
	Flush() error
	Close() error
}

type passCompressor struct {
	W io.Writer
}

func (pc passCompressor) Write(p []byte) (n int, err error) {
	return pc.W.Write(p)
}

func (pc passCompressor) Reset(w io.Writer) {
}

func (pc passCompressor) Flush() error {
	return nil
}

func (pc passCompressor) Close() error {
	return nil
}

type CompressorProvider interface {
	GetCompressor() Compressor
	PutCompressor(Compressor)
}

type GzipCompressorProvider struct {
	pool *sync.Pool
}

func NewGzipCompressorProvider() *GzipCompressorProvider {
	gcws := &GzipCompressorProvider{}
	gcws.pool = &sync.Pool{New: gcws.NewCompressor}
	return gcws
}

func (gcws *GzipCompressorProvider) NewCompressor() any {
	return gzip.NewWriter(io.Discard)
}

func (gcws *GzipCompressorProvider) GetCompressor() Compressor {
	return gcws.pool.Get().(Compressor)
}

func (gcws *GzipCompressorProvider) PutCompressor(cw Compressor) {
	cw.Reset(io.Discard)
	gcws.pool.Put(cw)
}

type ZlibCompressorProvider struct {
	pool *sync.Pool
}

func NewZlibCompressorProvider() *ZlibCompressorProvider {
	zcws := &ZlibCompressorProvider{}
	zcws.pool = &sync.Pool{New: zcws.NewCompressor}
	return zcws
}

func (zcws *ZlibCompressorProvider) NewCompressor() any {
	return zlib.NewWriter(io.Discard)
}

func (zcws *ZlibCompressorProvider) GetCompressor() Compressor {
	return zcws.pool.Get().(Compressor)
}

func (zcws *ZlibCompressorProvider) PutCompressor(cw Compressor) {
	cw.Reset(io.Discard)
	zcws.pool.Put(cw)
}
