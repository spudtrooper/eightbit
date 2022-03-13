package convert

import (
	"fmt"
	"image"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/noelyahan/mergi"
	"github.com/pkg/errors"
	goutilerrors "github.com/spudtrooper/goutil/errors"
	"github.com/spudtrooper/goutil/or"
)

type blockSizeRange struct {
	start, end, step int
}

func (b blockSizeRange) String() string {
	return fmt.Sprintf("blockSizeRange{start=%d end=%d step=%d}", b.start, b.end, b.step)
}

func MakeBlockSizeRange(start, end, step int) blockSizeRange {
	return blockSizeRange{start: start, end: end, step: step}
}

func animateBlock(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	start, end, step := opts.AnimateBlockSizeRange().start, opts.AnimateBlockSizeRange().end, opts.AnimateBlockSizeRange().step
	if start >= end {
		return nil, errors.Errorf("invalid block size range, start must be < end: %v", opts.AnimateBlockSizeRange())
	}
	if step <= 0 {
		return nil, errors.Errorf("invalid block size range, step must be > 0: %v", opts.AnimateBlockSizeRange())
	}

	log.Printf("animate from %d to %d by %d", start, end, step)

	numBlockSizes := (end - start + 1) / step
	blockSizes := make(chan int, numBlockSizes)
	go func() {
		for i := start; i <= end; i += step {
			blockSizes <- i
		}
		close(blockSizes)
	}()

	type sortedImage struct {
		img       image.Image
		blockSize int
	}

	// TODO: When I write to this channel, I get a deadlock. So write so an array guarded by a mutex.
	var sortableImages []sortedImage
	var imagesMu sync.Mutex
	addImage := func(img image.Image, blockSize int) {
		imagesMu.Lock()
		defer imagesMu.Unlock()
		sortableImages = append(sortableImages, sortedImage{img, blockSize})
	}
	// imgs := make(chan image.Image)
	errs := make(chan error)
	go func() {
		var cur int32
		var wg sync.WaitGroup
		for i, threads := 0, or.Int(opts.AnimateThreads(), 30); i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for blockSize := range blockSizes {
					res, err := blockMedianFromBlockSize(input, inputImage, blockSize, opts)
					atomic.AddInt32(&cur, 1)
					log.Printf("[%3d/%03d %02.2f%%] done with image for blockSize=%d err=%v",
						cur, numBlockSizes, 100.0*(float32(cur)/float32(numBlockSizes)), blockSize, err)
					if err != nil {
						errs <- err
						continue
					}
					if res.Image() == nil {
						errs <- errors.Errorf("nil image")
						continue
					}
					addImage(res.Image(), blockSize)
				}
			}()
		}
		wg.Wait()
		close(errs)
	}()

	ec := goutilerrors.MakeErrorCollector()
	for err := range errs {
		ec.Add(err)
	}
	if !ec.Empty() {
		return nil, ec.Build()
	}

	if opts.AnimateReverse() {
		sort.Slice(sortableImages, func(i, j int) bool {
			return sortableImages[j].blockSize < sortableImages[i].blockSize
		})
	} else {
		sort.Slice(sortableImages, func(i, j int) bool {
			return sortableImages[i].blockSize < sortableImages[j].blockSize
		})
	}
	var images []image.Image
	for _, si := range sortableImages {
		images = append(images, si.img)
	}
	log.Printf("creating gif from %d images", len(sortableImages))
	gif, err := mergi.Animate(images, 10)
	if err != nil {
		return nil, errors.Errorf("mergi.Animate: %v", err)
	}

	res := makeGIFConvertResult(gif)
	return res, nil
}

type animateConverter struct{ baseConverter }

func (c *animateConverter) OutputFileName(input string, opts ConvertOptions) string {
	ext := path.Ext(input)
	base := strings.Replace(path.Base(input), ext, "", 1)
	start, end, step := opts.AnimateBlockSizeRange().start, opts.AnimateBlockSizeRange().end, opts.AnimateBlockSizeRange().step
	return fmt.Sprintf("%s-%s-from_%d-to_%d-by_%d.gif", base, c.Name(), start, end, step)
}

func init() {
	globalReg.Register(&animateConverter{
		baseConverter{
			name: "animate_block",
			conv: animateBlock,
		}})
}
