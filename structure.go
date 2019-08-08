package gopng

var PNG_SIGNATURE []byte = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

type Color int

const (
    CLR_INV Color   = -1
    CLR_GRAY Color  = 0
    CLR_RGB Color   = 2
    CLR_PLT Color   = 3
    CLR_GRAYA Color = 4
    CLR_RGBA Color  = 6
)

var color_types = []Color{
    CLR_GRAY, CLR_INV,
    CLR_RGB, CLR_PLT,
    CLR_GRAYA, CLR_INV,
    CLR_RGBA}

var color_strs = []string{
    "greyscale", "",
    "RGB", "indexed",
    "greyscale with alpha", "",
    "RGBA"}

var filter_methods = []string{
    "None", "Sub", "Up",
    "Average", "Paerth"}

var interlace_methods = []string{
    "None", "Adam7"}

type GREYA struct{
    Black byte
    Alpha byte
}

type RGB8 struct{
    Red byte
    Green byte
    Blue byte
}

type RGB16 struct{
    Red uint16
    Green uint16
    Blue uint16
}

type RGBA struct{
    Red byte
    Green byte
    Blue byte
    Alpha byte
}

type IHDR struct{
    Width uint32
    Height uint32
    Bit_depth byte
    Color_type byte
    Compression_method byte
    Filter_method byte
    Interlace_method byte
}

type tIME struct{
    Year uint16
    Month byte
    Day byte
    Hour byte
    Minute byte
    Second byte
}

type tRNS struct{
    color_greya []GREYA
    color_plt []byte
    color_rgb []RGB8
}

type bKGD struct{
    color_greya GREYA
    color_plt byte
    color_rgb RGB16
}

type PngData struct{
    ihdr IHDR // header info
    color_type Color
    bpp int // bytes per pixel
    width uint32
    height uint32
    plte []RGB8 // palette data
    idat []byte // compressed data
    idat_dcomp []byte // decompressed data
    trns tRNS // simple transparency
    bkgd bKGD // default bg color
    time tIME // last modified timestamp
}
