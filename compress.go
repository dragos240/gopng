package gopng

import (
    "bytes"
    "compress/zlib"
    "io"
)

func Decompress(buf []byte) []byte{
    var imgdata bytes.Buffer
    b := bytes.NewReader(buf)

    r, err := zlib.NewReader(b)
    if err != nil{
        return []byte{}
    }
    defer r.Close()

    io.Copy(&imgdata, r)

    return imgdata.Bytes()
}
