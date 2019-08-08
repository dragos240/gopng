package gopng

import (
    "fmt"
    "os"
    "encoding/binary"
    "bytes"
    "strings"
    _ "io/ioutil"
)

func parse_optional_chunk(ctype string, length uint32,
                          fp *os.File, data *PngData) bool{
    // parse all the chunks we don't do anything with
    chunk_strs := []string{
            "tEXt",
            "zTXt",
            "iCCP",
            "pHYs",
            "iTXt"}
    for _, c := range chunk_strs{
        if ctype == c{
            fmt.Printf("\tFound type %s\n", c)
            fp.Seek(int64(length), 1)
            return true
        }
    }
    // parse all the chunks we do something with
    switch(ctype){
    case "bKGD":
        // r/g/b each contain 2 bytes
        fmt.Printf("\tFound type bKGD\n")
        color_type := data.color_type
        if color_type == CLR_PLT{
            buf := make([]byte, 1)
            _, err := fp.Read(buf)
            check(err)
            data.bkgd.color_plt = buf[0]
        } else if (color_type == CLR_GRAY ||
                color_type == CLR_GRAYA){
            err := binary.Read(fp, binary.BigEndian, &data.bkgd.color_greya)
            check(err)
        } else if (color_type == CLR_RGB ||
                color_type == CLR_RGBA){
            err := binary.Read(fp, binary.BigEndian, &data.bkgd.color_rgb)
            check(err)
        }
    case "tIME":
        fmt.Printf("\tFound type pHYs\n")
        err := binary.Read(fp, binary.BigEndian, &data.time)
        check(err)
        fmt.Printf("Last modified: %04d-%02d-%02d %02d:%02d:%02d\n",
                data.time.Year, data.time.Month, data.time.Day,
                data.time.Hour, data.time.Minute, data.time.Second)
    default:
        pos, _ := fp.Seek(0, 1)
        fmt.Printf("\tType: %s at 0x%x\n", ctype, pos)
        fp.Seek(int64(length), 1)
        return false
    }
    return true
}

func parse_IHDR(fp *os.File) IHDR{
    data := IHDR{}
    err := binary.Read(fp, binary.BigEndian, &data)
    check(err)
    return data
}

func Parse(fname string) int{
    var palettes []RGB8
    data := PngData{}
    fmt.Printf("Extracting image metadata from %s\n\n", fname)

    fp, err := os.Open(fname)
    check(err)
    defer fp.Close()

    // Seek past file signature
    signature := make([]byte, 8)
    fp.Read(signature)
    if !compare_slice(signature, PNG_SIGNATURE){
        fmt.Println("Invalid PNG file signature!")
        return 1
    }

    // loop through each block
    for{
        // length, type, data, crc32
        var length uint32
        ctype := make([]byte, 4)
        binary.Read(fp, binary.BigEndian, &length)
        fp.Read(ctype)
        pos, _ := fp.Seek(0, 1)
        switch string(ctype){
        case "IHDR":
            fmt.Printf("\tFound type IHDR\n")
            ihdr := parse_IHDR(fp)
            data.ihdr = ihdr
            data.color_type = color_types[ihdr.Color_type]
            data.bpp = get_bpp(data.color_type, ihdr.Bit_depth)
            data.width = ihdr.Width
            data.height = ihdr.Height
            fmt.Printf("width: %d, height: %d\n", ihdr.Width, ihdr.Height)
            fmt.Printf("bit depth: %d, ", ihdr.Bit_depth)
            fmt.Printf("color type: %s\n", color_strs[ihdr.Color_type])
            fmt.Printf("compression method: zlib, ")
            fmt.Printf("filter method: %s\n", filter_methods[ihdr.Filter_method])
            fmt.Printf("interlace method: %s\n", interlace_methods[ihdr.Interlace_method])
        case "PLTE":
            fmt.Printf("\tFound type PLTE\n")
            palettes = make([]RGB8, length/3)
            palette := RGB8{}
            var i int
            for i = 0; i < int(length)/3; i++{
                err := binary.Read(fp, binary.BigEndian, &palette)
                check(err)
                palettes[i] = palette
                //fmt.Printf("%d. #%02x%02x%02x\n", i, palette.Red, palette.Green, palette.Blue)
            }
            newpos, _ := fp.Seek(0, 1)
            if newpos - pos != int64(length){
                fmt.Printf("Length read (%d) is not equal to length of chunk (%d)!\n", (newpos - pos),
                        length)
                return 2
            }
            fmt.Printf("Read %d colors\n", i)
            data.plte = palettes
        case "IDAT":
            fmt.Printf("\tFound type IDAT\n")
            buf := make([]byte, length)
            fp.Read(buf)
            if len(data.idat) == 0{
                // no previous IDAT chunks written
                data.idat = buf
            } else {
                // previous IDAT chunks written
                data.idat = append(data.idat, buf...)
            }
            dcomp_data := Decompress(data.idat)
            if len(dcomp_data) != 0{
                data.idat_dcomp = dcomp_data
            }
        case "IEND":
            _rname := "raw_" + strings.Split(fname, ".")[0]
            var rname string
            if data.color_type == CLR_GRAY || data.color_type == CLR_GRAYA{
                rname = _rname + ".pgm"
            } else {
                rname = _rname + ".ppm"
            }
            raw_fp, err := os.Create(rname)
            check(err)
            defer raw_fp.Close()

            switch(data.color_type){
            case CLR_PLT:
                buf := new(bytes.Buffer)
                for key, value := range data.idat_dcomp{
                    _ = key
                    // first pixel on every line is extra, skip
                    if key%(int(data.width)+1) == 0{
                        continue
                    }
                    color := data.plte[value]
                    err := binary.Write(buf, binary.BigEndian, color)
                    check(err)
                }
                //raw_fp.Write(buf.Bytes())
                ppm_data := data_to_netpbm(buf.Bytes(), data)
                raw_fp.Write(ppm_data)
            case CLR_GRAY:
                fallthrough
            case CLR_GRAYA:
                raw_data := do_unfilter(data)
                pgm_data := data_to_netpbm(raw_data, data)
                fmt.Printf("\tWriting pgm image %s\n", rname)
                raw_fp.Write(pgm_data)
            case CLR_RGB:
                fallthrough
            case CLR_RGBA:
                raw_data := do_unfilter(data)
                ppm_data := data_to_netpbm(raw_data, data)
                fmt.Printf("\tWriting ppm image %s\n", rname)
                raw_fp.Write(ppm_data)
            case CLR_INV:
                fallthrough
            default:
                fmt.Println("FIXME: This shouldn't ever happen.")
                return 2
            }
            return 0
        default:
            ret := parse_optional_chunk(string(ctype), length, fp, &data)
            if !ret{
                fmt.Printf("Failed to parse optional chunk %s!\n", string(ctype))
                return 2
            }
        }
        // Seek past crc
        fp.Seek(4, 1)
    }
}
