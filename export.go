package gopng

import (
    "fmt"
)

func alpha_to_white(src_slice []byte) []byte{
    white := byte(0xff)
    len_slice := len(src_slice)
    masking_factor := float32(src_slice[len_slice-1])/255
    new_slice := make([]byte, len_slice)
    for i := 0; i < len_slice; i++{
        src_val := src_slice[i]
        new_val := byte(
                float32(src_val) * masking_factor +
                float32(white) * (1.0 - masking_factor))
        new_slice[i] = new_val
    }

    return new_slice
}

func data_to_netpbm(rawdata []byte, pngdata PngData) []byte{
    width := pngdata.width
    height := pngdata.height
    bpp := pngdata.bpp
    color_type := pngdata.color_type
    processed_data := rawdata
    if color_type >= CLR_GRAYA{
        nvals := bpp-1
        factor := float64(nvals)/float64(bpp)
        processed_data = make([]byte, int(float64(len(rawdata))*factor))
        j := 0
        for i := 0; i < len(rawdata); i+=bpp{
            slice := rawdata[i:i+nvals]
            if rawdata[i+nvals] < 0xff{
                slice = alpha_to_white(rawdata[i:i+bpp])
            }
            copy(processed_data[j:], slice)
            j += nvals
        }
    }

    p := 6
    if color_type == 0 || color_type == 4{
        p = 5
    }
    header := []byte(fmt.Sprintf("P%d\n%d %d\n255\n", p, width, height))
    pbmdata := make([]byte, len(header)+len(processed_data))
    copy(pbmdata, header)
    copy(pbmdata[len(header):], processed_data)

    return pbmdata
}
