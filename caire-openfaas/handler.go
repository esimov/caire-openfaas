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
	"github.com/esimov/caire"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type options struct {
	input          string `json:"input"`
	width          int    `json:width`
	height         int    `json:height`
	percentage     bool   `json:perc,string`
	square         bool   `json:square,string`
	scale          bool   `json:scale,string`
	debug          bool   `json:debug,string`
	blur           int    `json:blur`
	sobelThreshold int    `json:sobel`
	useFace        bool   `json:face,string`
	classifier     string `json:cascade`
}

// Handle a serverless request
func Handle(req []byte) string {
	var (
		options options
		data    []byte
		image   []byte
	)

	json.Unmarshal(req, &options)

	if val, exists := os.LookupEnv("input_mode"); exists && val == "url" {
		inputURL := strings.TrimSpace(options.input)

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

	_, err = io.Copy(tmpfile, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Sprintf("Unable to copy the source URI to the destionation file")
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
			BlurRadius:     options.blur,
			SobelThreshold: options.sobelThreshold,
			NewWidth:       options.width,
			NewHeight:      options.height,
			Percentage:     options.percentage,
			Square:         options.square,
			Debug:          options.debug,
			Scale:          options.scale,
			FaceDetect:     options.useFace,
			Classifier:     options.classifier,
		}

		filename := fmt.Sprintf("/tmp/%d.jpg", time.Now().UnixNano())

		output, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return fmt.Sprintf("Unable to open output file: %v", err)
		}
		defer os.Remove(filename)

		p.Process(options.input, output)

		image, err = ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Sprintf("Unable to read the resized image: %v", err)
		}
	}

	return string(image)
}
