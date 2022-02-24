package logh

import "io"

type stringWriter struct {
	buffer []byte
	index  uint64
	size   uint64
}

func NewStringWriter(size uint64) *stringWriter {

	return &stringWriter{
		buffer: make([]byte, size),
		index:  0,
		size:   size,
	}
}

// Write - implements the io.Writer interface
func (sb *stringWriter) Write(p []byte) (n int, err error) {

	for i := 0; i < len(p); i++ {

		if sb.index >= sb.size {
			return n, io.EOF
		}

		sb.buffer[sb.index] = p[i]
		sb.index++
		n++
	}

	return
}

// Write - implements the io.Writer interface
func (sb *stringWriter) Reset() {

	sb.index = 0
}

// Bytes - return the stored bytes
func (sb *stringWriter) Bytes() []byte {

	return sb.buffer[0:sb.index]
}
