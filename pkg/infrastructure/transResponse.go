package infrastructure

import (
	"bytes"
	"fmt"
	"strconv"
)

// transResponse a Trans response in bytes.
type transResponse struct {
	status string
	err    error
	body   []byte
}

// Map returns a new map from a response.
func (r *transResponse) Map() (map[string]string, error) {
	m := make(map[string]string)
	err := r.apply(func(key, value string) {
		if key == "status" {
			r.status = value
		}
		if key == "error" {
			r.err = fmt.Errorf("%s", value)
		}
		m[key] = value
	})
	return m, err
}

// Slice returns a new slice from a response.
func (r *transResponse) Slice() ([]map[string]string, error) {
	m := make(map[string][]string)
	err := r.apply(func(key, value string) {
		if key == "status" {
			r.status = value
		}
		if key == "error" {
			r.err = fmt.Errorf("%s", value)
		}
		m[key] = append(m[key], value)
	})
	if err != nil {
		return nil, err
	}
	slice := make([]map[string]string, getLength(m))
	for k, v := range m {
		for i, elem := range v {
			if len(slice[i]) == 0 {
				slice[i] = make(map[string]string)
			}
			slice[i][k] = elem
		}
	}
	return slice, nil
}

func getLength(m map[string][]string) int {
	for k, v := range m {
		if k != "status" && k != "error" {
			return len(v)
		}
	}
	return 0
}

func (r *transResponse) Status() string {
	return r.status
}

func (r *transResponse) SetStatus(s string) {
	r.status = s
}

func (r *transResponse) Error() error {
	return r.err
}

func (r *transResponse) SetError(err error) {
	r.err = err
}

// apply applies the given function on all key-value pairs of the response.
func (r *transResponse) apply(f func(key, value string)) error {
	n := 0
	for n < len(r.body) {
		blobLen := 0
		// Check if the value is a blob.
		if len(r.body) > n+5 && bytes.Equal(r.body[n:n+5], []byte("blob:")) {
			i := bytes.IndexByte(r.body[n+5:], ':')
			if i == -1 {
				return fmt.Errorf("trans: invalid blob %q", r.body[n:])
			}
			n += 5
			var err error
			blobLen, err = strconv.Atoi(string(r.body[n : n+i]))
			if err != nil {
				return fmt.Errorf("trans: cannot parse blob length: %v", err)
			}
			n += i + 1
		}

		// if current field is blob field - key terminator is newline, not ':'
		var i int
		if blobLen > 0 {
			i = bytes.IndexByte(r.body[n:], '\n')
		} else {
			i = bytes.IndexByte(r.body[n:], ':')
		}

		if i == -1 {
			return fmt.Errorf("trans: invalid key-value format: %q", r.body[n:])
		}

		key := string(r.body[n : n+i])
		n += i + 1

		vl := n + blobLen
		// if current field is not blob field - read until newline, if there is a glob field,
		// we already have value length in blobLen variable
		if blobLen <= 0 {
			i = bytes.IndexByte(r.body[n:], '\n')
			if i == -1 {
				return fmt.Errorf("trans: newline is missing: %q", r.body[n:])
			}
			vl += i
		}

		f(key, string(r.body[n:vl]))
		n = vl + 1
	}
	return nil
}
