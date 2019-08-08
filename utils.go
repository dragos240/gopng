package gopng

func compare_slice(ba1 []byte, ba2 []byte) bool{
    matches := true

    for i := 0; i < len(ba1); i++{
        if ba1[i] != ba2[i]{
            matches = false
            break
        }
    }

    return matches
}

func get_bpp(color_type Color, depth byte) int{
    var bpp int
    var multiplier int

    if depth == 16{
        multiplier = 2
    } else {
        // anything under 8 is rounded up
        multiplier = 1
    }

    switch color_type{
    case CLR_GRAY:
        bpp = 1 * multiplier
    case CLR_RGB:
        bpp = 3 * multiplier
    case CLR_PLT:
        bpp = 1 * multiplier
    case CLR_GRAYA:
        bpp = 2 * multiplier
    case CLR_RGBA:
        bpp = 4 * multiplier
    }
    
    return bpp
}

func check(err error){
    if err != nil{
        panic(err)
    }
}

func file_export_available(color Color) bool{
    //supported_color_types := []Color{CLR_PLT, CLR_RGBA}

    //for _, c := range supported_color_types{
    //    if color == c{
    //        return true
    //    }
    //}

    //return false
    return true
}
