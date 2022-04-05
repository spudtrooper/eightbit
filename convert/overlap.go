package convert

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"path"
	"sort"
	"strings"

	"github.com/spudtrooper/goutil/hist"
	"github.com/spudtrooper/goutil/or"
	"github.com/thomaso-mirodin/intmath/intgr"
)

type colorAggrFn func(inputImage image.Image, startY, endY, startX, endX int) color.Color

func overlap(input string, inputImage image.Image, blockSize int, opts ConvertOptions, aggr colorAggrFn, random bool) (ConvertResult, error) {
	minY, maxY := inputImage.Bounds().Min.Y, inputImage.Bounds().Max.Y
	minX, maxX := inputImage.Bounds().Min.X, inputImage.Bounds().Max.X

	outputImage := image.NewRGBA(image.Rect(minX, minY, maxX, maxY))

	inc := or.Int(blockSize, 10)

	colorHist := hist.MakeHistogram()

	for y := minY; y < maxY; y += inc {
		for x := minX; x < maxX; x += inc {
			startY := intgr.Max(y-inc, minY)
			endY := intgr.Min(y+inc, maxY)
			startX := intgr.Max(x-inc, minX)
			endX := intgr.Min(x+inc, maxX)
			mc := aggr(inputImage, startY, endY, startX, endX)
			colorHist.Add(colorName(mc), 1)
			mr, mg, mb, ma := mc.RGBA()
			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					var c color.Color
					if random {
						c = color.RGBA{
							R: uint8(mr + uint32(30-rand.Int()%60)),
							G: uint8(mg + uint32(30-rand.Int()%60)),
							B: uint8(mb + uint32(30-rand.Int()%60)),
							A: uint8(ma + uint32(30-rand.Int()%60)),
						}
					} else {
						c = color.RGBA{
							R: uint8(mr),
							G: uint8(mg),
							B: uint8(mb),
							A: uint8(ma),
						}
					}
					outputImage.Set(x, y, c)
				}
			}
		}
	}

	if opts.ColorHist() {
		log.Println("Printing color histogram...\n" + hist.HistString(colorHist))
	}

	res := makeImageConvertResult(outputImage)
	return res, nil
}

func medianColor(inputImage image.Image, startY, endY, startX, endX int) color.Color {
	var rs, gs, bs, as []int
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			c := inputImage.At(x, y)
			r, g, b, a := c.RGBA()
			rs = append(rs, int(r))
			gs = append(gs, int(g))
			bs = append(bs, int(b))
			as = append(as, int(a))
		}
	}

	sort.Ints(rs)
	sort.Ints(gs)
	sort.Ints(bs)
	sort.Ints(as)

	var mr, mg, mb, ma uint8
	if m := len(rs) / 2; m%2 == 0 {
		mr = uint8(rs[m])
		mg = uint8(gs[m])
		mb = uint8(bs[m])
		ma = uint8(as[m])
	} else {
		mr = uint8((rs[m-1] + rs[m]) / 2)
		mg = uint8((gs[m-1] + gs[m]) / 2)
		mb = uint8((bs[m-1] + bs[m]) / 2)
		ma = uint8((as[m-1] + as[m]) / 2)

	}
	median := color.RGBA{
		R: mr,
		G: mg,
		B: mb,
		A: ma,
	}

	return median
}

func meanColor(inputImage image.Image, startY, endY, startX, endX int) color.Color {
	var sumr, sumb, sumg, suma uint32
	var n uint32
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			c := inputImage.At(x, y)
			r, g, b, a := c.RGBA()
			sumr += r
			sumg += g
			sumb += b
			suma += a
			n++
		}
	}

	mr := sumr / n
	mg := sumg / n
	mb := sumb / n
	ma := suma / n

	mean := color.RGBA{
		R: uint8(mr),
		G: uint8(mg),
		B: uint8(mb),
		A: uint8(ma),
	}

	return mean
}

type overlapConverter struct{ baseConverter }

func (c *overlapConverter) OutputFileName(input string, opts ConvertOptions) string {
	ext := path.Ext(input)
	base := strings.Replace(path.Base(input), ext, "", 1)
	return fmt.Sprintf("%s-%s-%04d%s", base, c.Name(), opts.BlockSize(), ext)
}

func overlapMean(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	return overlap(input, inputImage, opts.BlockSize(), opts, meanColor, true)
}

func overlapMedian(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	return overlap(input, inputImage, opts.BlockSize(), opts, medianColor, true)
}

func blockMean(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	return overlap(input, inputImage, opts.BlockSize(), opts, meanColor, false)
}

func blockMedian(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	return blockMedianFromBlockSize(input, inputImage, opts.BlockSize(), opts)
}

func blockMedianFromBlockSize(input string, inputImage image.Image, blockSize int, opts ConvertOptions) (ConvertResult, error) {
	return overlap(input, inputImage, blockSize, opts, medianColor, false)
}

func init() {
	globalReg.Register(&overlapConverter{
		baseConverter{
			name: "overlap_mean",
			conv: overlapMean,
		}})
	globalReg.Register(&overlapConverter{
		baseConverter{
			name: "overlap_median",
			conv: overlapMedian,
		}})
	globalReg.Register(&overlapConverter{
		baseConverter{
			name: "block_mean",
			conv: blockMean,
		}})
	globalReg.Register(&overlapConverter{
		baseConverter{
			name: "block_median",
			conv: blockMedian,
		}})
}
