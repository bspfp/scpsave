package gzipio

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

func NewCompressReader(infile string) (io.Reader, error) {
	f, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)
	defer w.Close()

	if _, err := io.Copy(w, f); err != nil {
		return nil, err
	}
	w.Close()
	return bytes.NewBuffer(buf.Bytes()), nil
}

func NewDecompressWriter(outfile string) io.WriteCloser {
	return &decompressWriter{outfile: outfile}
}

type decompressWriter struct {
	outfile string
	buf     bytes.Buffer
}

func (w *decompressWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *decompressWriter) Close() error {
	r, err := gzip.NewReader(bytes.NewReader(w.buf.Bytes()))
	if err != nil {
		return err
	}
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return os.WriteFile(w.outfile, decompressed, 0644)
}
