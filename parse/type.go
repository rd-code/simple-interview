package parse

import (
    "bytes"
    "errors"
    "io"
    "strconv"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2019-03-18 23:27
 **/

type Parser interface {
    Parse(io.Reader) (interface{}, error)
}

type bulkString struct {
}

func (b *bulkString) Parse(r io.Reader) (interface{}, error) {
    if r == nil {
        return nil, errors.New("the reader cannot be nil")
    }
    return b.parse(r)
}

func (b *bulkString) parse(r io.Reader) (res string, err error) {
    var size int64
    if size, err = parseSize(r); err != nil {
        return
    }

    res, err = b.parseBody(r, size)
    return
}

func parseSize(r io.Reader) (res int64, err error) {
    var array []byte
    var bt byte
    if array, bt, err = readToFlag(r, '\r', '\n'); err != nil {
        return
    }
    if bt != '\r' {
        err = errors.New("invalid input flag '\r'")
        return
    }
    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, r, 1); err != nil {
        return
    }
    if bt, err = buffer.ReadByte(); err != nil {
        return
    }
    if bt != '\n' {
        err = errors.New("the size must end with '\n'")
    }
    res, err = strconv.ParseInt(string(array), 10, 64)
    return

}

func (b *bulkString) parseBody(r io.Reader, size int64) (res string, err error) {
    if size == -1 {
        return
    }
    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, r, size); err != nil {
        return
    }
    res = buffer.String()
    buffer.Reset()
    if _, err = io.CopyN(buffer, r, 2); err != nil {
        return
    }
    var bt byte

    if bt, err = buffer.ReadByte(); err != nil {
        return
    }
    if bt != '\r' {
        err = errors.New("body must end with '\r''\n'")
        return
    }
    if bt, err = buffer.ReadByte(); err != nil {
        return
    }
    if bt != '\n' {
        err = errors.New("body must end with '\r''\n'")
        return
    }
    return
}

type integerString struct {
}

func (i *integerString) Parse(r io.Reader) (interface{}, error) {
    if r == nil {
        return nil, errors.New("the input cannot be nil")
    }
    return i.parse(r)
}

func (i *integerString) parse(r io.Reader) (res int64, err error) {
    var array []byte
    var b byte
    if array, b, err = readToFlag(r, '\r', '\n'); err != nil {
        return
    }
    if b != '\r' {
        err = errors.New("the message must end with '\r''\n'")
        return
    }
    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, r, 1); err != nil {
        return
    }
    if b, err = buffer.ReadByte(); err != nil {
        return
    }
    if b != '\n' {
        err = errors.New("the message must end wieth '\r' '\n'")
        return
    }
    res, err = strconv.ParseInt(string(array), 10, 64)
    return
}

type errString struct {
}

func (e *errString) Parse(r io.Reader) (interface{}, error) {
    if r == nil {
        return nil, errors.New("the input cannot be nil")
    }
    return e.parse(r)
}

func (e *errString) parse(r io.Reader) (res, err error) {
    var array []byte
    var b byte
    if array, b, err = readToFlag(r, '\r', '\n'); err != nil {
        return
    }
    if b != '\r' {
        err = errors.New("the message must end with '\r''\n'")
        return
    }
    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, r, 1); err != nil {
        return
    }
    if b, err = buffer.ReadByte(); err != nil {
        return
    }
    if b != '\n' {
        err = errors.New("the message must end wieth '\r' '\n'")
        return
    }
    res = errors.New(string(array))
    return
}

type simpleString struct {
}

func (s *simpleString) Parse(r io.Reader) (interface{}, error) {
    if r == nil {
        return nil, errors.New("the input cannot be nil")
    }
    return s.parse(r)
}

func (s *simpleString) parse(r io.Reader) (res string, err error) {
    var array []byte
    var b byte
    if array, b, err = readToFlag(r, '\r', '\n'); err != nil {
        return
    }
    if b != '\r' {
        err = errors.New("the message must end with '\r''\n'")
        return
    }
    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, r, 1); err != nil {
        return
    }
    if b, err = buffer.ReadByte(); err != nil {
        return
    }
    if b != '\n' {
        err = errors.New("the message must end wieth '\r' '\n'")
        return
    }
    res = string(array)
    return
}

type arraysString struct {
}

func (a *arraysString) Parse(r io.Reader) (interface{}, error) {
    if r == nil {
        return nil, errors.New("the input cannot be nil")
    }
    return a.parse(r)
}

func (a *arraysString) parse(r io.Reader) (res []interface{}, err error) {
    var size int64
    if size, err = parseSize(r); err != nil {
        return
    }
    res = make([]interface{}, 0, size)
    var i int64
    for i = 0; i < size; i++ {
        var v interface{}
        if v, err = Parse(r); err != nil {
            return
        }
        res = append(res, v)
    }
    return
}

func readToFlag(r io.Reader, array ...byte) (res []byte, flag byte, err error) {
    if len(array) == 0 {
        err = errors.New("the bytes cannot empty")
        return
    }
    data := make(map[byte]struct{}, len(array))
    for _, b := range array {
        data[b] = struct{}{}
    }
    sw := &StringWriter{
        buffer: &bytes.Buffer{},
        flags:  data,
    }
    res, flag, err = sw.parse(r)
    return
}

type StringWriter struct {
    buffer *bytes.Buffer
    flags  map[byte]struct{}
}

func (sw *StringWriter) write(p byte) (ok bool) {
    if _, ok = sw.flags[p]; ok {
        return
    }
    sw.buffer.WriteByte(p)
    return
}

func (sw *StringWriter) parse(r io.Reader) (res []byte, flag byte, err error) {
    if r == nil {
        err = errors.New("the input cannot be nil")
        return
    }
    buffer := &bytes.Buffer{}
    var f byte
    for {
        if _, err = io.CopyN(buffer, r, 1); err != nil {
            return
        }
        if f, err = buffer.ReadByte(); err != nil {
            return
        }
        if ok := sw.write(f); ok {
            flag = f
            res = sw.buffer.Bytes()
            return
        }
    }
}
