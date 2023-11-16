package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/processImage", processImageHandler)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
	fmt.Println("Server listening at port :4000")
}

func processImageHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32) // maxMemory 32MB
	var format string

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the photo key - First Approach
	file, h, err := r.FormFile("image")
	message := r.FormValue("message")
	format = strings.Split(h.Filename, ".")[1]
	fmt.Println(format)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error at image decoding: ", err)
		return
	}

	img_message := modifyIMG(&img, message)

	file_new, err := os.Create("new_" + h.Filename)
	if err != nil {
		fmt.Println("Error at file: ", err)
		return
	}
	defer file_new.Close()

	if format == "png" {
		err = png.Encode(file_new, img_message)
	} else {
		err = jpeg.Encode(file_new, img_message, nil)
	}

	if err != nil {
		fmt.Println("Error at image encoding:", err)
		return
	}

	img_send, err := os.Open("new_" + h.Filename)
	defer img_send.Close()

	w.Header().Set("Content-Type", "image/"+format)
	io.Copy(w, img_send)
	w.WriteHeader(200)
	return
}

func modifyIMG(img *image.Image, message string) image.Image {
	bounds := (*img).Bounds()
	modifiedImage := image.NewRGBA(bounds)
	byteArray := []byte(message)
	var binaryArray []string
	var binaryIntArray []int

	for _, b := range byteArray {
		binaryRepresentation := strconv.FormatInt(int64(b), 2)
		paddedBinary := fmt.Sprintf("%08s", binaryRepresentation)
		binaryArray = append(binaryArray, paddedBinary)
	}

	for _, character := range binaryArray {
		for _, b := range character {
			intb, _ := strconv.ParseInt(string(b), 2, 8)
			binaryIntArray = append(binaryIntArray, int(intb))
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 1 {
			currentColor := (*img).At(x, y)
			red, green, blue, alpha := currentColor.RGBA()

			index := y*bounds.Max.Y + x*3
			if index <= len(binaryIntArray)-2 {
				r_value := setValue(red, binaryIntArray[index])
				g_value := setValue(green, binaryIntArray[index+1])
				b_value := uint32(0)
				if index < len(binaryIntArray)-4 && index%6 != 0 {
					b_value = setValue(blue, binaryIntArray[index+2])
				} else if index < len(binaryIntArray)-4 && index%6 == 0 {
					b_value = setValue(blue, 0)
				} else {
					b_value = setValue(blue, 1)
				}

				fmt.Println("Original colors: ", red, green, blue)
				fmt.Println("New colors: ", r_value, g_value, b_value)
				modifiedImage.Set(x, y, color.RGBA{
					uint8(r_value >> 8),
					uint8(g_value >> 8),
					uint8(b_value >> 8),
					uint8(alpha >> 8)})
			} else {
				modifiedImage.Set(x, y,
					color.RGBA{
						uint8(red >> 8),
						uint8(green >> 8),
						uint8(blue >> 8),
						uint8(alpha >> 8)})
			}
		}
	}

	return modifiedImage
}

func setValue(pixelValue uint32, bitValue int) uint32 {
	if bitValue == 0 && pixelValue%2 != 0 || bitValue == 1 && pixelValue%2 == 0 {
		if pixelValue == 255 {
			pixelValue--
		} else {
			pixelValue++
		}
	}
	return pixelValue
}
