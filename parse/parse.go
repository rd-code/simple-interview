package parse

import (
    "bytes"
    "errors"
    "io"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2019-03-18 23:10
 **/
func Parse(r io.Reader) (interface{}, error) {
    return parse(r)
}

var impl map[byte]Parser

func init() {
    impl = map[byte]Parser{
        '$': &bulkString{},
        ':': &integerString{},
        '-': &errString{},
        '+': &simpleString{},
        '*': &arraysString{},
    }
}

func parse(r io.Reader) (resp interface{}, err error) {
    if r == nil {
        err = errors.New("the read cannot be nil")
        return
    }
    var flag byte
    if flag, err = getFirstByte(r); err != nil {
        return
    }
    return impl[flag].Parse(r)
}

func getFirstByte(r io.Reader) (byte, error) {
    buffer := &bytes.Buffer{}
    if _, err := io.CopyN(buffer, r, 1); err != nil {
        return 0, err
    }
    if buffer.Len() != 1 {
        return 0, errors.New("read first byte failed")
    }
    return buffer.Bytes()[0], nil
}
