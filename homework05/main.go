package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"unsafe"

	"github.com/InternatBlackhole/cudago/cuda"
)

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return gray
}

func main() {
	// Read command line arguments
	inputImageStr := flag.String("i", "", "input image")
	outputImageStr := flag.String("o", "", "output image")
	flag.Parse()
	if *inputImageStr == "" || *outputImageStr == "" {
		panic("Missing input or output image arguments\nUsage: go run main.go -i input.png -o output.png")
	}

	// Initialize CUDA API on OS thread

}

func proccessImageOnGPU(inputImageStr, outputImageStr string) {
	var err error
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// Open image file
	inputFile, err := os.Open(inputImageStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Read image " + inputImageStr)
	defer inputFile.Close()

	// Decode image
	inputImage, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	// Convert image to grayscale
	inputImageGray := rgbaToGray(inputImage)

	// Get image size
	imgSize := inputImageGray.Bounds().Size()
	size := uint64(imgSize.X * imgSize.Y)

	// Allocate memory on the device for input and output image
	imgInDevice, err := cuda.DeviceMemAlloc(size)
	if err != nil {
		panic(err)
	}
	defer imgInDevice.Free()

	imgOutDevice, err := cuda.DeviceMemAlloc(size)
	if err != nil {
		panic(err)
	}
	defer imgOutDevice.Free()

	// Copy image to device
	err = imgInDevice.MemcpyToDevice(unsafe.Pointer(&inputImageGray.Pix[0]), size)
	if err != nil {
		panic(err)
	}

	// Specify grid and block size
	dimBlock := cuda.Dim3{X: 1, Y: 1, Z: 1}
	dimGrid := cuda.Dim3{X: 1, Y: 1, Z: 1}

	// Call the kernel to execute on the device
	err = cudaMedianFilter.Process(dimGrid, dimBlock, imgInDevice.Ptr, imgOutDevice.Ptr, int32(imgSize.X), int32(imgSize.Y))
	if err != nil {
		panic(err)
	}

	// Copy image back to host
	imgOutHost := make([]byte, size)
	imgOutDevice.MemcpyFromDevice(unsafe.Pointer(&imgOutHost[0]), size)

	// Save image to file
	outputFile, err := os.Create(outputImageStr)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	outputImage := image.NewGray(inputImageGray.Bounds().Bounds())
	outputImage.Pix = imgOutHost

	err = png.Encode(outputFile, outputImage)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image saved to " + outputImageStr)
}

func processImageOnCPU(inputImageStr, outputImageStr string) {
	inputFile, err := os.Open(inputImageStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Read image " + inputImageStr)
	defer inputFile.Close()

	inputImage, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	inputImageGray := rgbaToGray(inputImage)
	outputImage := medianFilter(inputImageGray)

	outputFile, err := os.Create(outputImageStr)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outputImage)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image saved to " + outputImageStr)
}

func medianFilter(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)
	windowSize := 3
	offset := windowSize / 2

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			var window []uint8
			for wy := -offset; wy <= offset; wy++ {
				for wx := -offset; wx <= offset; wx++ {
					clampedX := max(0, min(x+wx, bounds.Max.X-1))
					clampedY := max(0, min(y+wy, bounds.Max.Y-1))

					pixel := img.GrayAt(clampedX, clampedY).Y
					window = append(window, pixel)
				}
			}

			sort.Slice(window, func(i, j int) bool { return window[i] < window[j] })
			median := window[len(window)/2]

			output.SetGray(x, y, color.Gray{Y: median})
		}
	}

	return output
}
