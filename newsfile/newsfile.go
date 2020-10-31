// Copyright © 2020 The Stentor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package newsfile

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	readLength = 1024
)

func WriteFragments(fn, startComment string, data []byte, keepHeader bool) error {
	tf, err := writeFragments(fn, startComment, data, keepHeader)
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

func writeFragments(fn, startComment string, data []byte, keepHeader bool) (string, error) {
	dst, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer dst.Close()

	src, err := os.Open(fn)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		// news file doesn't exist, so just write data to it
		if _, err := dst.Write(data); err != nil {
			return "", err
		}

		return dst.Name(), nil
	}
	defer src.Close()

	if err := copyIntoFile(dst, src, startComment, data, keepHeader); err != nil {
		return "", err
	}

	return dst.Name(), nil
}

// nolint:gocognit // try to simplify this at some point
func copyIntoFile(dst io.Writer, src io.Reader, startComment string, data []byte, keepHeader bool) error {
	var partialBuf []byte = make([]byte, 0)
	var startFound bool
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

		if !startFound {
			// find index of startComment
			idx := bytes.Index(readBuf, []byte(startComment))
			if idx >= 0 {
				idx += len(startComment)
				if keepHeader {
					// need to keep the existing header,
					// so splice data into readBuf
					readBuf = append(readBuf[:idx], append(data, readBuf[idx:]...)...)
				} else {
					readBuf = append(data, readBuf[idx:]...)
				}
				startFound = true
			} else {
				// trim off everything after the last full line
				lastIdx := bytes.LastIndex(readBuf, []byte("\n"))
				if lastIdx >= 0 && lastIdx < len(readBuf)-1 {
					readBuf, partialBuf = readBuf[:idx+1], readBuf[idx+1:]
				}
			}
		}

		// write readBuf to dst
		// if we already found the start comment
		// or we're keeping the header
		if startFound || keepHeader {
			if _, err := dst.Write(readBuf); err != nil {
				return err
			}
		}

		if rerr == io.EOF {
			// we're done so flush anything left in partialBuf to dst
			if len(partialBuf) > 0 {
				if _, err := dst.Write(partialBuf); err != nil {
					return err
				}
			}
			break
		}
	}

	if !startFound {
		return errors.New("no start comment found")
	}

	return nil
}
