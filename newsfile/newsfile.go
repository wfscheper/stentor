package newsfile

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var (
	startCommentRST      = ".. stentor output starts\n"
	startCommentMarkdown = "<!-- stentor output starts -->\n"

	readLength = 1024
)

func WriteFragments(fn string, startComment, data []byte) error {
	tf, err := writeFragments(fn, startComment, data)
	if err != nil {
		return err
	}

	// attempt to move temp file over the top of fn
	if err := os.Rename(tf, fn); err != nil {
		defer os.Remove(tf)
		return err
	}

	return nil
}

func writeFragments(fn string, startComment, data []byte) (string, error) {
	src, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if err := copyIntoFile(dst, src, startComment, data); err != nil {
		return "", err
	}

	return dst.Name(), nil
}

// nolint:gocognit // try to simplify this at some point
func copyIntoFile(dst io.Writer, src io.Reader, startComment, data []byte) error {
	var partialBuf []byte = make([]byte, 0)
	var dataWritten bool
	for {
		// read from source
		readBuf := make([]byte, readLength)
		n, rerr := src.Read(readBuf)
		if rerr != nil && rerr != io.EOF {
			return rerr
		}

		// trim off zero bytes if we under read
		readBuf = readBuf[:n]

		if len(partialBuf) > 0 {
			// we had a partial read so prepend it to readBuf
			readBuf = append(partialBuf, readBuf...)
			partialBuf = make([]byte, 0)
		}

		if !dataWritten {
			// find index of startComment
			idx := bytes.Index(readBuf, startComment)
			if idx >= 0 {
				// splice in data
				idx += len(startComment)
				readBuf = append(readBuf[:idx], append(data, readBuf[idx:]...)...)
				dataWritten = true
			} else {
				// trim off everything after the last full line
				lastIdx := bytes.LastIndex(readBuf, []byte("\n"))
				if lastIdx >= 0 && lastIdx < len(readBuf)-1 {
					readBuf, partialBuf = readBuf[:idx+1], readBuf[idx+1:]
				}
			}
		}

		// write to dst
		if _, err := dst.Write(readBuf); err != nil {
			return err
		}

		if rerr == io.EOF {
			// we're done
			if len(partialBuf) > 0 {
				if _, err := dst.Write(partialBuf); err != nil {
					return err
				}
			}
			break
		}
	}

	if !dataWritten {
		return errors.New("no start comment found")
	}

	return nil
}
