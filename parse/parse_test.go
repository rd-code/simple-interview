package parse

import (
    "errors"
    "strings"
    "testing"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2019-03-19 21:08
 **/

func TestParse(t *testing.T) {
    //s := "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n"
    //res := Parse(strings.NewReader(s))

    type Demo struct {
        input  string
        output interface{}
    }

    params := []Demo{
        {
            input:  "$6\r\nfoobar\r\n",
            output: "foobar",
        },
        {
            input:  "$0\r\n\r\n",
            output: "",
        },
        {
            input:  "$-1\r\n",
            output: "",
        },
        {
            input:  ":-1000\r\n",
            output: int64(-1000),
        },
        {
            input:  "-Error message\r\n",
            output: errors.New("Error message"),
        },
        {
            input:  "+OK\r\n",
            output: "OK",
        },
        {
            input:  "*6\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n-Error message\r\n",
            output: []interface{}{int64(1), int64(2), int64(3), int64(4), "foobar", errors.New("Error message")},
        },
    }
    for _, param := range params {
        input := strings.NewReader(param.input)
        v, err := Parse(input)
        if err != nil {
            t.Error("happened err", param.input, errString{})
        }
        switch v.(type) {
        case int64, string:
            if v != param.output {
                t.Error("expected", param.output, "actual", v, "input", param.input)
            }
        case error:
            t1 := v.(error)
            t2 := param.output.(error)
            if t1.Error() != t2.Error() {
                t.Error("expected", param.output, "actual", v, "input", param.input)
            }
        default:
            array1 := v.([]interface{})
            array2 := param.output.([]interface{})
            for i := range array1 {
                v := array1[i]
                switch v.(type) {
                case int64, string:
                    if v != array2[i] {
                        t.Error("expected", array2[i], "actual", v, "input", param.input)
                    }
                case error:
                    t1 := v.(error)
                    t2 := array2[i].(error)
                    if t1.Error() != t2.Error() {
                        t.Error("expected", array2[i], "actual", v, "input", param.input)
                    }
                default:
                    t.Error("unknow type", param.input)
                }
            }
        }
    }
}
