// MIT License
//
// Copyright (c) 2018 Endre Simo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package function

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/esimov/caire"
)

// Options struct contains the resize options defined in the json object.
type Options struct {
	Input          string `json:"input"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Percentage     bool   `json:"perc,string"`
	Square         bool   `json:"square,string"`
	Scale          bool   `json:"scale,string"`
	Debug          bool   `json:"debug,string"`
	Blur           int    `json:"blur"`
	SobelThreshold int    `json:"sobel"`
	UseFace        bool   `json:"face,string"`
}

// Handle a serverless request
func Handle(req []byte) string {
	var (
		options Options
		data    []byte
		image   []byte
	)

	// Decode json
	json.Unmarshal(req, &options)

	if val, exists := os.LookupEnv("input_mode"); exists && val == "url" {
		inputURL := strings.TrimSpace(options.Input)

		// Retrieve the url and decode the response body.
		res, err := http.Get(inputURL)
		if err != nil {
			return fmt.Sprintf("Unable to download image file from URI: %s, status %v", inputURL, res.Status)
		}
		defer res.Body.Close()

		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Sprintf("Unable to read response body: %s", err)
		}
	} else {
		var decodeError error
		data, decodeError = base64.StdEncoding.DecodeString(string(req))
		if decodeError != nil {
			data = req
		}

		contentType := http.DetectContentType(req)
		if contentType != "image/jpeg" && contentType != "image/png" {
			return fmt.Sprintf("Only jpeg or png images, either raw uncompressed bytes or base64 encoded are acceptable inputs, you uploaded: %s", contentType)
		}
	}
	tmpfile, err := ioutil.TempFile("/tmp", "image")
	if err != nil {
		return fmt.Sprintf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Copy the image binary data into the temporary file.
	_, err = io.Copy(tmpfile, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Sprintf("Unable to copy the source URI to the destination file")
	}

	var output string
	query, err := url.ParseQuery(os.Getenv("Http_Query"))
	if err == nil {
		output = query.Get("output")
	}

	if val, exists := os.LookupEnv("output_mode"); exists {
		output = val
	}

	if output == "image" || output == "json_image" {
		p := &caire.Processor{
			BlurRadius:     options.Blur,
			SobelThreshold: options.SobelThreshold,
			NewWidth:       options.Width,
			NewHeight:      options.Height,
			Percentage:     options.Percentage,
			Square:         options.Square,
			Debug:          options.Debug,
			Scale:          options.Scale,
			FaceDetect:     options.UseFace,
			Classifier:     "./data/facefinder",
		}

		input, err := os.Open(tmpfile.Name())
		if err != nil {
			return fmt.Sprintf("Unable to open the temporary image file: %v", err)
		}

		filename := fmt.Sprintf("/tmp/%d.jpg", time.Now().UnixNano())

		output, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return fmt.Sprintf("Unable to open output file: %v", err)
		}
		defer os.Remove(filename)

		// A single line to rescale the image.
		err = p.Process(input, output)
		if err != nil {
			return fmt.Sprintf("Error on resize process: %v", err)
		}

		// Retrieve the resized image.
		image, err = ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Sprintf("Unable to read the resized image: %v", err)
		}
	}
	return string(image)
}
