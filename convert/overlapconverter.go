package convert

import (
	"image"
	"image/color"
	"math/rand"
	"sort"

	"github.com/spudtrooper/goutil/hist"
	"github.com/thomaso-mirodin/intmath/intgr"
)

func overlapMeanConverter(input string, inputImage image.Image, opts ConvertOptions) (image.Image, error) {
	return genericOverlap(input, inputImage, opts, meanColor)
}

func overlapMedianConverter(input string, inputImage image.Image, opts ConvertOptions) (image.Image, error) {
	return genericOverlap(input, inputImage, opts, medianColor)
}

type colorAggrFn func(inputImage image.Image, startY, endY, startX, endX int) color.Color

func genericOverlap(input string, inputImage image.Image, opts ConvertOptions, aggr colorAggrFn) (image.Image, error) {
	minY, maxY := inputImage.Bounds().Min.Y, inputImage.Bounds().Max.Y
	minX, maxX := inputImage.Bounds().Min.X, inputImage.Bounds().Max.X

	outputImage := image.NewRGBA(image.Rect(minX, minY, maxX, maxY))

	const inc = 10

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
					c := color.RGBA{
						R: uint8(mr + uint32(30-rand.Int()%60)),
						G: uint8(mg + uint32(30-rand.Int()%60)),
						B: uint8(mb + uint32(30-rand.Int()%60)),
						A: uint8(ma + uint32(30-rand.Int()%60)),
					}
					outputImage.Set(x, y, c)
				}
			}
		}
	}

	log.Println("Printing color histogram...\n" + hist.HistString(colorHist))

	return outputImage, nil
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
