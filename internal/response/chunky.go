package response

import "fmt"

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	n := len(p)
	_, err := w.writer.Write([]byte(fmt.Sprintf("%x\r\n", n)))
	if err != nil {
		return 0, err
	}
	_, err = w.writer.Write([]byte(fmt.Sprintf("%s\r\n", p)))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}
