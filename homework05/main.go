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
	"time"

	"homework05/cudaMedianFilter"
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
	inputImageStr := flag.String("i", "", "input image")
	outputImageStr := flag.String("o", "", "output image")
	flag.Parse()
	if *inputImageStr == "" || *outputImageStr == "" {
		panic("Missing input or output image arguments\nUsage: go run main.go -i input.png -o output.png")
	}

	start := time.Now()
	processImageOnCPU(*inputImageStr, *outputImageStr)
	elapsed_cpu := time.Since(start)
	fmt.Printf("Time on CPU: %fs\n", elapsed_cpu.Seconds());

	start = time.Now()
	processImageOnGPU(*inputImageStr, *outputImageStr)
	elapsed_gpu := time.Since(start)
	fmt.Printf("Time on GPU: %fs\n", elapsed_gpu.Seconds());

	fmt.Printf("Pohitritev S = %f\n", elapsed_cpu.Seconds() / elapsed_gpu.Seconds())
}

func processImageOnGPU(inputImageStr, outputImageStr string) {
	// Initialize CUDA API on OS thread
	var err error
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// Open image file
	inputFile, err := os.Open(inputImageStr + ".png")
	if err != nil {
		panic(err)
	}
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
	dimBlock := cuda.Dim3{X: 16, Y: 16, Z: 1}
	dimGrid := cuda.Dim3{X: uint32(imgSize.X) + dimBlock.X-1 / dimBlock.X, Y: uint32(imgSize.Y) + dimBlock.Y / dimBlock.Y, Z: 1}

	// Call the kernel to execute on the device
	start := time.Now()
	err = cudaMedianFilter.MedianFilter(dimGrid, dimBlock, imgInDevice.Ptr, imgOutDevice.Ptr, int32(imgSize.X), int32(imgSize.Y))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elapsed only method on GPU: %f\n", time.Since(start).Seconds())

	// Copy image back to host
	imgOutHost := make([]byte, size)
	imgOutDevice.MemcpyFromDevice(unsafe.Pointer(&imgOutHost[0]), size)

	// Save image to file
	outputFile, err := os.Create(outputImageStr + "-GPU.png")
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
}

func processImageOnCPU(inputImageStr, outputImageStr string) {
	inputFile, err := os.Open(inputImageStr + ".png")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	inputImage, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	inputImageGray := rgbaToGray(inputImage)
	outputImage := medianFilter(inputImageGray)

	outputFile, err := os.Create(outputImageStr + "-CPU.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outputImage)
	if err != nil {
		panic(err)
	}
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
