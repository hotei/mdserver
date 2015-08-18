// images.go

package main

import (
	// go 1.4.2 std lib
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	//
	"github.com/hotei/bmp"
)

// handle any images that md files call for
func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("imageHandler: r.URL.Path %q\n", r.URL.Path)
	workDir := serverRoot + imageURL[1:]
	fmt.Printf("workDir(%s)\n", workDir)
	fmt.Printf("r.URL.Path(%s)\n", r.URL.Path)
	imageName := workDir + r.URL.Path[len(imageURL):]
	fmt.Printf("imageHandler: imageName = %s\n", imageName)
	ext := strings.ToLower(filepath.Ext(imageName))
	//fmt.Printf("ext = %s\n",ext)
	if ext == ".bmp" {
		bmpWriteOut(imageName, w)
		return
	}
	if ext == ".png" {
		pngWriteOut(imageName, w)
		return
	}
	if ext == ".jpeg" {
		jpegWriteOut(imageName, w)
		return
	}
	if ext == ".jpg" {
		jpegWriteOut(imageName, w)
		return
	}
	if ext == ".svg" {
		rawWriteOut(imageName, w)
		return
	}
}

// just hand the raw file over to the browser
func rawWriteOut(fileName string, w http.ResponseWriter) {
	rawBuf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("rawWriteOut: cant open file %s\n", fileName)
		return
	}
	w.Write(rawBuf)
}

func jpegWriteOut(imageName string, w http.ResponseWriter) {
	bf, err := os.Open(imageName)
	if err != nil {
		fmt.Printf("jpegWriteOut: cant open image %s\n", imageName)
		return
	}
	img, err := jpeg.Decode(bf)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("jpegWriteOut: Decode failed for %s of error:%v\n", imageName, err)))
		return
	}
	b := make([]byte, 0, 1000*1000) // b will expand as needed
	wo := bytes.NewBuffer(b)
	err = png.Encode(wo, img)
	if err != nil {
		fmt.Printf("imageHandler: png encode failed for %s\n", imageName)
		return
	}
	w.Write(wo.Bytes())
}

// testing
func pngWriteOut(imageName string, w http.ResponseWriter) {
	fmt.Printf("pngWriteOut: imageName = %s\n", imageName)
	bf, err := os.Open(imageName)
	if err != nil {
		fmt.Printf("pngWriteOut: cant open image %s\n", imageName)
		return
	}
	img, err := png.Decode(bf)
	if err != nil {
		fmt.Printf("pngWriteOut: image decode failed for %s png\n", imageName)
		w.Write([]byte(fmt.Sprintf("image Decode failed for %s png error:%v\n", imageName, err)))
		return
	}
	b := make([]byte, 0, 1000*1000) // b will expand as needed
	wo := bytes.NewBuffer(b)
	err = png.Encode(wo, img)
	if err != nil {
		fmt.Printf("pngWriteOut: png encode failed for %s\n", imageName)
		return
	}
	w.Write(wo.Bytes())
}

func bmpWriteOut(imageName string, w http.ResponseWriter) {
	fmt.Printf("bmpWriteOut: imageName = %s\n", imageName)
	bf, err := os.Open(imageName)
	if err != nil {
		fmt.Printf("bmpWriteOut: cant open bmp %s\n", imageName)
		return
	}
	img, err := bmp.Decode(bf)
	if err != nil {
		fmt.Printf("bmpWriteOut: bmp decode failed for %s\n", imageName)
		w.Write([]byte(fmt.Sprintf("Decode failed for %s\n", imageName)))
		return
	}
	b := make([]byte, 0, 10000)
	wo := bytes.NewBuffer(b)
	err = png.Encode(wo, img)
	if err != nil {
		fmt.Printf("bmpWriteOut: png encode failed for %s\n", imageName)
		return
	}
	w.Write(wo.Bytes())
}

func bmpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("bmpHandler: r.URL.Path %q\n", r.URL.Path)
	fileName := serverRoot + r.URL.Path[1:]
	fmt.Printf("bmpHandler: fname = %s\n", fileName)
	ext := strings.ToLower(filepath.Ext(fileName))
	//fmt.Printf("ext = %s\n",ext)
	if ext == ".bmp" {
		bmpWriteOut(fileName, w)
		return
	}
}
