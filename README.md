# eightbit

A converter to create shitty 8-bit like images. It now creates other types of shitty images, too.

## Usage

To install:

```
go install github.com/spudtrooper/eightbit
```

To run:

```
eightbit --input <INPUT> --output <OUTPUT>
```

## Examples

| In                                                         | Out                                                          |
| ---------------------------------------------------------- | ------------------------------------------------------------ |
| ![in/magritte](./examples/in/magritte.jpg)                 | ![out/magritte](./examples/out/magritte.jpg)                 |
| ![in/la-belle-lurette](./examples/in/la-belle-lurette.jpg) | ![out/la-belle-lurette](./examples/out/la-belle-lurette.jpg) |
| ![in/starry-night](./examples/in/starry-night.jpg)         | ![out/starry-night](./examples/out/starry-night.jpg)         |

## Animation

You can create an animation ([like this](./examples/animation/stella.gif)) with the `block_median` converter, with something like this:

```bash
for f in `seq 1 100`; do
  eightbit --input <input-image> --converters block_median --force --block_size $f
done
ls -1 *.png > up.txt
ls -1 *.png | sort -r > down.txt
convert -delay 1 -loop 0 `cat up.txt | xargs` `cat down.txt | xargs` output.gif
```

Example: ![animation](./examples/animation/stella.gif)
