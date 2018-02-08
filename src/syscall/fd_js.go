// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

type jsFile struct {
	defaultFileImpl
	jsFD    int
	path    string
	entries []string
	pos     int64
	seeked  bool
}

func init() {
	newFD(&jsFile{jsFD: 0})
	newFD(&jsFile{jsFD: 1})
	newFD(&jsFile{jsFD: 2})
	setMapFD(func(fd int) int {
		files.Lock()
		f := files.tab[fd]
		files.Unlock()
		return f.impl.(*jsFile).jsFD
	})
}

func setMapFD(f func(fd int) int)

func (f *jsFile) stat(st *Stat_t) error {
	jsSt, err := fsCall("fstatSync", f.jsFD)
	if err != nil {
		return err
	}
	setStat(st, jsSt)
	return nil
}

func (f *jsFile) read(b []byte) (int, error) {
	if f.seeked {
		n, err := f.pread(b, f.pos)
		f.pos += int64(n)
		return n, err
	}

	n, err := fsCall("readSync", f.jsFD, b, 0, len(b))
	if err != nil {
		return 0, err
	}
	n2 := n.Int()
	f.pos += int64(n2)
	return n2, err
}

func (f *jsFile) write(b []byte) (int, error) {
	if f.seeked {
		n, err := f.pwrite(b, f.pos)
		f.pos += int64(n)
		return n, err
	}

	n, err := fsCall("writeSync", f.jsFD, b, 0, len(b))
	if err != nil {
		return 0, err
	}
	n2 := n.Int()
	f.pos += int64(n2)
	return n2, err
}

func (f *jsFile) pread(b []byte, offset int64) (int, error) {
	n, err := fsCall("readSync", f.jsFD, b, 0, len(b), offset)
	if err != nil {
		return 0, err
	}
	return n.Int(), nil
}

func (f *jsFile) pwrite(b []byte, offset int64) (int, error) {
	n, err := fsCall("writeSync", f.jsFD, b, 0, len(b), offset)
	if err != nil {
		return 0, err
	}
	return n.Int(), nil
}

func (f *jsFile) close() error {
	_, err := fsCall("closeSync", f.jsFD)
	f.jsFD = -1
	return err
}

func (f *jsFile) seek(offset int64, whence int) (int64, error) {
	f.seeked = true
	switch whence {
	case 0:
		f.pos = offset
		return f.pos, nil
	case 1:
		f.pos += offset
		return f.pos, nil
	case 2:
		var st Stat_t
		if err := f.stat(&st); err != nil {
			return 0, err
		}
		f.pos = st.Size + offset
		return f.pos, nil
	default:
		return 0, EINVAL
	}
}
