package containerd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/containerd/containerd/cio"
	"github.com/weaveworks/ignite/pkg/util"
)

// ^P^Q bytes taken from Docker
var detachBytes = []byte{16, 17}

// igniteIO provides handles for duplicating stdout and stderr to
// a file for logging as well as an detach capturer for detecting
// a detach sequence from stdin input
type igniteIO struct {
	input  *detachCapturer
	output *fileCloner
}

func newIgniteIO(logFile string) (*igniteIO, error) {
	var err error
	i := &igniteIO{
		input: newDetachCapturer(os.Stdin, detachBytes),
	}

	if i.output, err = newFileCloner(os.Stdout, logFile); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *igniteIO) Opt() cio.Opt {
	return cio.WithStreams(i.input, i.output, i.output)
}

func (i *igniteIO) Detach() <-chan struct{} {
	return i.input.detachC
}

func (i *igniteIO) Close() (err error) {
	if err = i.output.Close(); err != nil {
		return
	}

	return
}

// FileCloner is a io.WriteCloser implementation that
// clones everything it's being written to to both
// the given file and io.Writer
type fileCloner struct {
	writer io.Writer
	file   *os.File
}

func newFileCloner(writer io.Writer, fileName string) (c *fileCloner, err error) {
	c = &fileCloner{writer: writer}
	c.file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return
}

var _ io.WriteCloser = &fileCloner{}

func (c *fileCloner) Write(p []byte) (n1 int, err error) {
	var n2 int
	if n1, err = c.file.Write(p); err == nil {
		if n2, err = c.writer.Write(p); err == nil && n1 != n2 {
			err = fmt.Errorf("write mismatch; file: %d, writer: %d", n1, n2)
		}
	}

	return
}

func (c *fileCloner) Close() error {
	return c.file.Close()
}

// detachCapturer is an io.Reader implementation that
// can fire a signal using a channel when a specific
// sequence is input (excluded from reader output)
type detachCapturer struct {
	reader    io.Reader
	detachSeq []byte
	progress  []byte
	detachC   chan struct{}
}

func newDetachCapturer(reader io.Reader, detachSeq []byte) *detachCapturer {
	return &detachCapturer{
		reader:    reader,
		detachSeq: detachSeq,
		progress:  make([]byte, 0, len(detachSeq)),
		detachC:   make(chan struct{}),
	}
}

var _ io.Reader = &detachCapturer{}

func (c *detachCapturer) Read(p []byte) (n int, err error) {
	// Use the given reader to read to a buffer
	nRead, pRead := 0, make([]byte, len(p))
	if nRead, err = c.reader.Read(pRead); err != nil {
		return
	}

	// appendP overwrites indices in the given p buffer by counting with n
	appendP := func(bs ...byte) {
		for _, b := range bs {
			p[n] = b
			n++
		}
	}

	for i := 0; i < nRead; i++ {
		b := pRead[i] // Retrieve a read byte

		// If the byte matches the current progress in the
		// detach sequence, add it to the progress
		if b == c.detachSeq[len(c.progress)] {
			c.progress = append(c.progress, b)

			// If the progress is complete, clear it and signal a detach
			if len(c.progress) == len(c.detachSeq) {
				c.progress = c.progress[:0]
				c.detachC <- struct{}{}
			}
		} else if len(c.progress) > 0 {
			// The detach sequence was interrupted, append any
			// cached bytes to p and clear the progress
			appendP(c.progress...)
			c.progress = c.progress[:0]
		} else {
			// The byte is not part of the detach sequence and
			// no progress needs to be cleared, just append
			appendP(b)
		}
	}

	return
}

type dummyReader struct{}

var _ io.Reader = &dummyReader{}

func (r *dummyReader) Read(_ []byte) (int, error) {
	return 0, nil
}

type logRetriever struct {
	reader *os.File
	writer *os.File
	output *fileCloner
}

func newlogRetriever(logFile string) (l *logRetriever, err error) {
	l = &logRetriever{}
	if l.reader, l.writer, err = os.Pipe(); err != nil {
		return
	}

	if l.output, err = newFileCloner(l.writer, logFile); err != nil {
		return
	}

	if util.FileExists(logFile) {
		var reader io.ReadCloser
		if reader, err = os.Open(logFile); err != nil {
			return
		}
		defer util.DeferErr(&err, reader.Close)

		if _, err = io.Copy(l.writer, reader); err != nil {
			return
		}
	}

	return
}

var _ io.ReadCloser = &logRetriever{}

type fakeIO struct {
	config cio.Config
}

func (f *fakeIO) Config() cio.Config {
	return f.config
}

func (f *fakeIO) Cancel() {}

func (f *fakeIO) Wait() {}

func (f *fakeIO) Close() error { return nil }

func (l *logRetriever) Attach() cio.Attach {
	return func(fifos *cio.FIFOSet) (i cio.IO, err error) {
		time.Sleep(time.Second)

		var stdout io.ReadCloser
		stdout, err = os.Open(fifos.Stdout)
		if err != nil {
			return
		}
		defer util.DeferErr(&err, stdout.Close)

		async, ok := stdout.(syscall.Conn)
		if !ok {
			err = fmt.Errorf("cannot cast to syscall.Conn")
			return
		}

		sc, err := async.SyscallConn()
		if err != nil {
			return
		}

		if err = sc.Control(func(_ uintptr) {}); err != nil {
			return
		}

		//if err = syscall.SetNonblock(int(stdout.Fd()), true); err != nil {
		//	return
		//}

		if err = syscall.SetNonblock(int(os.Stdout.Fd()), true); err != nil {
			return
		}

		buf := bufio.NewReader(stdout)

		fmt.Println(buf.Buffered())

		for {
			//fmt.Println("Read")
			if _, err = io.CopyN(os.Stdout, buf, int64(buf.Buffered())); err == io.EOF {
				break
			} else if err != nil {
				return
			}
		}
		//buf := make([]byte, 100)
		//_, _ = io.ReadFull(stdout, buf)
		//io.

		//fi, err := stdout.Stat()
		//fmt.Println("STDOUT SIZE:", fi.Size(), err)
		//fmt.Println(string(buf))

		return &fakeIO{fifos.Config}, nil
	}
}

func (l *logRetriever) Ready(io cio.IO) error {
	waitPipes := []string{io.Config().Stdout, io.Config().Stderr}

	for {
		var size int64
		for _, pipe := range waitPipes {
			fi, err := os.Stat(pipe)
			if err != nil {
				return err
			}

			fmt.Println("File:", pipe)
			fmt.Println("Size:", fi.Size())
			size += fi.Size()
		}

		if size == 0 {
			break
		}
	}

	return l.writer.Close()
}

func (l *logRetriever) Opt() cio.Opt {
	return cio.WithStreams(&dummyReader{}, l.output, l.output)
}

func (l *logRetriever) Read(p []byte) (n int, err error) {
	return l.reader.Read(p)
}

func (l *logRetriever) Close() error {
	return l.reader.Close()
}
