// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"bufio"
	"fmt"
	"strings"
	"io/ioutil"
	"flag"
	"encoding/base64"
)

const credentials = ".twitter/drchase-benchmark-bot"

func main() {
	var mediaFiles []string
	var inputFileName string
	var toTweet string
	flag.Var((*repeatedString)(&mediaFiles), "m", "name of file containing twitter appropriate media")
	flag.StringVar(&inputFileName, "i", inputFileName,  "name of file containing tweet text")
	flag.StringVar(&toTweet, "t", toTweet,  "tweet text")
	flag.Usage = func() {
		fmt.Printf(`Usage of %s:
Post a tweet from either standard input or the file specified by -i,
with optional media from -m, to the account whose keys etc are found
in %s
`, os.Args[0], credentials)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Has to be true.
	if len(mediaFiles) > 4 {
		fmt.Printf("Can only upload 4 media files\n")
		os.Exit(1)
	}

	if toTweet == "" {
		// Read tweet text
		inputFile := os.Stdin
		if inputFileName != "" {
			var err error
			inputFile, err = os.Open(inputFileName)
			check(err, "Could not open -i file %s", inputFileName)
		}
		tBytes, err := ioutil.ReadAll(inputFile)
		check(err)
		toTweet = string(tBytes)
	} else {
		if inputFileName != "" {
			panic("Cannot have both -t and -i options")
		}
	}

	// Read media from files
	urlv := make(map[string][]string)
	var mediaStrings []string
	for _, s := range mediaFiles {
		inputFile, err := os.Open(s)
		check(err, "Could not open mediafile %s", s)
		bytes, err := ioutil.ReadAll(inputFile)
		check(err, "Could not read mediafile %s", s)
		mediaStrings = append(mediaStrings, base64.StdEncoding.EncodeToString(bytes))
	}

	// Defer this as long as possible, so that all the errors come out fast.
	api := GetApi(readMapFile(credentials))

	// Upload media if any
	if len(mediaStrings) > 0 {
		var mediaValues = []string{""} // all the media ids are comma-separated in a single string
		for i, s := range mediaStrings {
			fmt.Printf("Uploading %s\n", mediaFiles[i])
			mv, err := api.UploadMedia(s)
			check(err, "Unable to upload media in file %s", mediaFiles[i])
			if i > 0 {
				mediaValues[0] = mediaValues[0] + ","
			}
			mediaValues[0] = mediaValues[0] + mv.MediaIDString
		}
		urlv["media_ids"] = mediaValues
	}

	tweet, err := api.PostTweet(toTweet, urlv)
	check(err, "Was not able to post tweet \n%s\n", toTweet)
	fmt.Printf("Tweet is \n%+v\n", tweet)
}

func check(err error, messages ...interface{}) {
	if err != nil {
		if len(messages) > 0 {
			maybeFmt := messages[0]
			if fm, ok := maybeFmt.(string); ok && len(messages) > 1 {
				fmt.Printf(fm, messages[1:])
				fmt.Println()
			} else {
				fmt.Println(maybeFmt)
			}
		}
		panic(err)
	}
}
func readMapFile(fileName string) map[string]string {
	r := make(map[string]string)
	f, err := os.Open(fileName)
	check(err, "Could not open configuration (map) file %s", fileName)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		colon := strings.Index(line, ":")
		key := strings.ToLower(strings.TrimSpace(line[:colon]))
		value := strings.TrimSpace(line[colon+1:])
		r[key] = value
	}
	check(scanner.Err())

	return r
}

func GetApi(m map[string]string) *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(m["access token"], m["access token secret"],
		m["api key"], m["api secret"])
}

type repeatedString []string

func (c *repeatedString) String() string {
	s := ""
	for i, v := range *c {
		if i > 0 {
			s += ","
		}
		s += v
	}
	return s
}

func (c *repeatedString) Set(s string) error {
	*c = append(*c, s)
	return nil
}

func (c *repeatedString) IsBoolFlag() bool {
	return false
}
