package pulse

import (
	"errors"
	"unsafe"
)

/*
#cgo pkg-config: libpulse-simple

#include <stdlib.h>
#include <pulse/simple.h>
*/
import "C"

type SampleFormat C.pa_sample_format_t

const (
	SAMPLE_U8        SampleFormat = C.PA_SAMPLE_U8
	SAMPLE_ALAW      SampleFormat = C.PA_SAMPLE_ALAW
	SAMPLE_ULAW      SampleFormat = C.PA_SAMPLE_ULAW
	SAMPLE_S16LE     SampleFormat = C.PA_SAMPLE_S16LE
	SAMPLE_S16BE     SampleFormat = C.PA_SAMPLE_S16BE
	SAMPLE_FLOAT32LE SampleFormat = C.PA_SAMPLE_FLOAT32LE
	SAMPLE_FLOAT32BE SampleFormat = C.PA_SAMPLE_FLOAT32BE
	SAMPLE_S32LE     SampleFormat = C.PA_SAMPLE_S32LE
	SAMPLE_S32BE     SampleFormat = C.PA_SAMPLE_S32BE
	SAMPLE_S24LE     SampleFormat = C.PA_SAMPLE_S24LE
	SAMPLE_S24BE     SampleFormat = C.PA_SAMPLE_S24BE
	SAMPLE_S24_32LE  SampleFormat = C.PA_SAMPLE_S24_32LE
	SAMPLE_S24_32BE  SampleFormat = C.PA_SAMPLE_S24_32BE
	SAMPLE_MAX       SampleFormat = C.PA_SAMPLE_MAX
	SAMPLE_INVALID   SampleFormat = C.PA_SAMPLE_INVALID
)

const (
	SAMPLE_S16NE     SampleFormat = C.PA_SAMPLE_S16NE
	SAMPLE_FLOAT32NE SampleFormat = C.PA_SAMPLE_FLOAT32NE
	SAMPLE_S32NE     SampleFormat = C.PA_SAMPLE_S32NE
	SAMPLE_S24NE     SampleFormat = C.PA_SAMPLE_S24NE
	SAMPLE_S24_32NE  SampleFormat = C.PA_SAMPLE_S24_32NE
	SAMPLE_S16RE     SampleFormat = C.PA_SAMPLE_S16RE
	SAMPLE_FLOAT32RE SampleFormat = C.PA_SAMPLE_FLOAT32RE
	SAMPLE_S32RE     SampleFormat = C.PA_SAMPLE_S32RE
	SAMPLE_S24RE     SampleFormat = C.PA_SAMPLE_S24RE
	SAMPLE_S24_32RE  SampleFormat = C.PA_SAMPLE_S24_32RE
)

const SAMPLE_FLOAT32 SampleFormat = C.PA_SAMPLE_FLOAT32

type StreamDirection C.pa_stream_direction_t

const (
	STREAM_NODIRECTION StreamDirection = C.PA_STREAM_NODIRECTION
	STREAM_PLAYBACK    StreamDirection = C.PA_STREAM_PLAYBACK
	STREAM_RECORD      StreamDirection = C.PA_STREAM_RECORD
	STREAM_UPLOAD      StreamDirection = C.PA_STREAM_UPLOAD
)

type SampleSpec struct {
	Format   SampleFormat
	Rate     uint32
	Channels uint8
}

type Stream struct {
	simple *C.pa_simple
}

// Capture creates a new stream for recording and returns its pointer.
func Capture(clientName, streamName string, spec *SampleSpec) (*Stream, error) {
	return NewStream("", clientName, STREAM_RECORD, "", streamName, spec, nil, nil)
}

// Playback creates a new stream for playback and returns its pointer.
func Playback(clientName, streamName string, spec *SampleSpec) (*Stream, error) {
	return NewStream("", clientName, STREAM_PLAYBACK, "", streamName, spec, nil, nil)
}

type ChannelMap int // FIXME: STUB
type BufferAttr int // FIXME: STUB

func NewStream(
	serverName, clientName string,
	dir StreamDirection,
	deviceName, streamName string,
	spec *SampleSpec,
	cmap *ChannelMap,
	battr *BufferAttr,
) (*Stream, error) {

	s := new(Stream)

	var server *C.char
	if serverName != "" {
		server = C.CString(serverName)
		defer C.free(unsafe.Pointer(server))
	}

	var dev *C.char
	if deviceName != "" {
		dev := C.CString(deviceName)
		defer C.free(unsafe.Pointer(dev))
	}

	name := C.CString(clientName)
	defer C.free(unsafe.Pointer(name))
	stream_name := C.CString(streamName)
	defer C.free(unsafe.Pointer(stream_name))

	var err C.int

	ss := C.pa_sample_spec{
		format:   C.pa_sample_format_t(spec.Format),
		rate:     C.uint32_t(spec.Rate),
		channels: C.uint8_t(spec.Channels),
	}

	s.simple = C.pa_simple_new(
		server,
		name,
		C.pa_stream_direction_t(dir),
		dev,
		stream_name,
		&ss,
		nil,
		nil,
		&err,
	)

	if err == C.PA_OK {
		return s, nil
	}
	return s, errors.New("whatever")
}

// Stream.Free closes the stream and frees the associated memory.
// The stream becomes invalid after this has been called.
// This should usually be deferred immediately after obtaining a stream.
func (s *Stream) Free() {
	C.pa_simple_free(s.simple)
}

// Stream.Drain blocks until all buffered data has finished playing.
func (s *Stream) Drain() (int, error) {
	var err C.int
	written := C.pa_simple_drain(s.simple, &err)
	if err == C.PA_OK {
		return int(written), nil
	}
	return int(written), errors.New("summt went wrong")
}

// Stream.Flush flushes the playback buffer, discarding any audio therein
func (s *Stream) Flush() (int, error) {
	var err C.int
	flushed := C.pa_simple_flush(s.simple, &err)
	if err == C.PA_OK {
		return int(flushed), nil
	}
	return int(flushed), errors.New("summt went wrong")
}

// Stream.Write writes the given data to the stream,
// blocking until the data has been written.
func (s *Stream) Write(data []byte) (int, error) {
	var err C.int
	written := C.pa_simple_write(
		s.simple,
		unsafe.Pointer(&data[0]),
		C.size_t(len(data)),
		&err,
	)
	if err == C.PA_OK {
		return int(written), nil
	}
	return int(written), errors.New("summt went wrong")
}

// Stream.Read reads data from the stream,
// blocking until it has filled the provided slice.
func (s *Stream) Read(data []byte) (int, error) {
	var err C.int
	written := C.pa_simple_read(
		s.simple,
		unsafe.Pointer(&data[0]),
		C.size_t(len(data)),
		&err,
	)
	if err == C.PA_OK {
		return int(written), nil
	}
	return int(written), errors.New("summt went wrong")
}

// Stream.Latency returns the playback latency in microseconds.
func (s *Stream) Latency() (uint64, error) {
	var err C.int
	lat := C.pa_simple_get_latency(s.simple, &err)
	if err == C.PA_OK {
		return uint64(lat), nil
	}
	return uint64(lat), errors.New("summt went wrong")
}
