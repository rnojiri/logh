package logh

import "io"

// StringWriter - writes in memory
type StringWriter struct {
	buffer []byte
	index  uint64
	size   uint64
}

func NewStringWriter(size uint64) *StringWriter {

	return &StringWriter{
		buffer: make([]byte, size),
		index:  0,
		size:   size,
	}
}

// Write - implements the io.Writer interface
func (sb *StringWriter) Write(p []byte) (n int, err error) {

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
func (sb *StringWriter) Reset() {

	sb.index = 0
}

// Bytes - return the stored bytes
func (sb *StringWriter) Bytes() []byte {

	return sb.buffer[0:sb.index]
}
