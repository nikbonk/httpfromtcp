package response

import "github.com/nikbonk/httpfromtcp/internal/headers"

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		_, err := w.writer.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
