package gopng

import (
    "fmt"
    "math"
)

const (
    FLT_INVALID int = -1
    FLT_NONE    int = 0
    FLT_SUB     int = 1
    FLT_UP      int = 2
    FLT_AVG     int = 3
    FLT_PTH     int = 4
)

func paeth_predict(_a byte, _b byte, _c byte) int{
    a := int(_a)
    b := int(_b)
    c := int(_c)
    //fmt.Printf("a b c: %x %x %x\n", a, b, c)
    p := a + b - c
    pa := math.Abs(float64(p - a))
    pb := math.Abs(float64(p - b))
    pc := math.Abs(float64(p - c))

    if pa <= pb && pa <= pc{
        return a
    } else if pb <= pc{
        return b
    } else {
        return c
    }
}

type ImgAttrs struct{
    len_sl int
    bpp int
}

func unsub(sl []byte, attrs ImgAttrs) []byte{
    raw_sl := make([]byte, 0)
    for slx := 0; slx < attrs.len_sl; slx++{
        c_byte := sl[slx]
        var p_byte byte
        if len(raw_sl) >= attrs.bpp{
            p_byte = raw_sl[slx-attrs.bpp]
        }

        result := c_byte+p_byte
        raw_sl = append(raw_sl, result)
    }

    return raw_sl
}

func unup(sl []byte, prev_sl []byte, attrs ImgAttrs) []byte{
    raw_sl := make([]byte, 0)
    for slx := 0; slx < attrs.len_sl; slx++{
        c_byte := sl[slx]
        p_byte := prev_sl[slx]

        result := c_byte+p_byte
        raw_sl = append(raw_sl, result)
    }

    return raw_sl
}

func unavg(sl []byte, prev_sl []byte, attrs ImgAttrs) []byte{
    raw_sl := make([]byte, 0)
    for slx := 0; slx < attrs.len_sl; slx++{
        c_byte := sl[slx]
        var prev uint16
        if len(raw_sl) >= attrs.bpp{
            prev = uint16(raw_sl[slx-attrs.bpp])
        }
        up := uint16(prev_sl[slx])
        p_byte := byte((prev+up)/2)

        result := c_byte+p_byte
        raw_sl = append(raw_sl, result)
    }

    return raw_sl
}

func unpth(sl []byte, prev_sl []byte, attrs ImgAttrs) []byte{
    raw_sl := make([]byte, 0)
    for slx := 0; slx < attrs.len_sl; slx++{
        c_byte := sl[slx]
        var prev byte
        var prevPrior byte
        if len(raw_sl) >= attrs.bpp{
            prev = raw_sl[slx-attrs.bpp]
            prevPrior = prev_sl[slx-attrs.bpp]
        }
        prior := prev_sl[slx]
        p_byte := byte(paeth_predict(prev, prior, prevPrior))

        result := c_byte + p_byte
        raw_sl = append(raw_sl, result)
    }

    fmt.Printf("")
    return raw_sl
}

func do_unfilter(data PngData) []byte{
    filtdata := data.idat_dcomp
    width := int(data.ihdr.Width)
    bpp := data.bpp

    len_data := len(filtdata)
    len_sl := width*bpp
    attrs := ImgAttrs{len_sl, bpp}
    filter_count := []int{0, 0, 0, 0, 0}
    rawdata := make([]byte, 0)
    for x := 0; x < len_data; x += len_sl{
        filter_type := int(filtdata[x])
        x++
        sl := filtdata[x:x+len_sl]
        var prev_sl []byte
        var raw_sl []byte

        if filter_type > FLT_SUB{
            if len(rawdata) == 0{
                prev_sl = make([]byte, len_sl)
            } else {
                prev_sl = rawdata[len(rawdata)-len_sl:]
            }
        }

        switch filter_type{
        case FLT_NONE:
            raw_sl = sl
            filter_count[0]++
        case FLT_SUB:
            raw_sl = unsub(sl, attrs)
            filter_count[1]++
        case FLT_UP:
            raw_sl = unup(sl, prev_sl, attrs)
            filter_count[2]++
        case FLT_AVG:
            raw_sl = unavg(sl, prev_sl, attrs)
            filter_count[3]++
        case FLT_PTH:
            raw_sl = unpth(sl, prev_sl, attrs)
            filter_count[4]++
        }
        rawdata = append(rawdata, raw_sl...)
    }

    return rawdata
}

