/*
   Copyright Naohiro Heya

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */


package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"net/http"
	"io"
)

type Tweets []struct {
	Source   string `json:"source"`
	Entities struct {
		UserMentions []interface{} `json:"user_mentions"`
		Media        []struct {
			ExpandedURL   string `json:"expanded_url"`
			Indices       []int  `json:"indices"`
			URL           string `json:"url"`
			MediaURL      string `json:"media_url"`
			IDStr         string `json:"id_str"`
			ID            int64  `json:"id"`
			MediaURLHTTPS string `json:"media_url_https"`
			Sizes         []struct {
				H      int    `json:"h"`
				Resize string `json:"resize"`
				W      int    `json:"w"`
			} `json:"sizes"`
			MediaAlt   string `json:"media_alt"`
			DisplayURL string `json:"display_url"`
		} `json:"media"`
		Hashtags []interface{} `json:"hashtags"`
		Urls     []interface{} `json:"urls"`
	} `json:"entities"`
	Geo struct {
	} `json:"geo"`
	IDStr           string `json:"id_str"`
	Text            string `json:"text"`
	RetweetedStatus struct {
		Source   string `json:"source"`
		Entities struct {
			UserMentions []interface{} `json:"user_mentions"`
			Media        []struct {
				ExpandedURL   string `json:"expanded_url"`
				Indices       []int  `json:"indices"`
				URL           string `json:"url"`
				MediaURL      string `json:"media_url"`
				IDStr         string `json:"id_str"`
				ID            int64  `json:"id"`
				MediaURLHTTPS string `json:"media_url_https"`
				Sizes         []struct {
					H      int    `json:"h"`
					Resize string `json:"resize"`
					W      int    `json:"w"`
				} `json:"sizes"`
				MediaAlt   string `json:"media_alt"`
				DisplayURL string `json:"display_url"`
			} `json:"media"`
			Hashtags []interface{} `json:"hashtags"`
			Urls     []interface{} `json:"urls"`
		} `json:"entities"`
		Geo struct {
		} `json:"geo"`
		IDStr     string `json:"id_str"`
		Text      string `json:"text"`
		ID        int64  `json:"id"`
		CreatedAt string `json:"created_at"`
		User      struct {
			Name                 string `json:"name"`
			ScreenName           string `json:"screen_name"`
			Protected            bool   `json:"protected"`
			IDStr                string `json:"id_str"`
			ProfileImageURLHTTPS string `json:"profile_image_url_https"`
			ID                   int64  `json:"id"`
			Verified             bool   `json:"verified"`
		} `json:"user"`
	} `json:"retweeted_status"`
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	User      struct {
		Name                 string `json:"name"`
		ScreenName           string `json:"screen_name"`
		Protected            bool   `json:"protected"`
		IDStr                string `json:"id_str"`
		ProfileImageURLHTTPS string `json:"profile_image_url_https"`
		ID                   int    `json:"id"`
		Verified             bool   `json:"verified"`
	} `json:"user"`
}


func Download(url, downloadpath string) error {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP Error")
		return err
	}
	defer res.Body.Close()
	file, err := os.Create(downloadpath)
	if err != nil {
		fmt.Println("create Error")
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("io copy Error")
		return err
	}
	return nil
}


func PicDownloader(path, downloadpath string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Open Error", path)
		return err
	}

	/* remove first line */
	for i, v := range raw {
		if v == '\n' {
			raw = raw[i+1:]
			break
		}
	}
	data := new(Tweets)
	if err := json.Unmarshal(raw, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return err
	}
	for _, v := range *data {
		timestamp := v.CreatedAt
		for _, m := range v.Entities.Media {
			if v.RetweetedStatus.User.ID == 0 {
				url := m.MediaURL
				tmp := strings.Split(url, "/")
				filename := timestamp + "_" + tmp[len(tmp)-1]
				filefullpath := filepath.Join(downloadpath, filename)
				err = Download(url, filefullpath)
				if err != nil {
					fmt.Println("Error", url ,filefullpath, path)
				} else {
					fmt.Println(filefullpath, " is created!")
				}

			}
		}
		
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("TwiPicDownloader [tweet directory path] [download directory path]")
		os.Exit(1)
	}
	p := filepath.Join(args[0], "/data/js/tweets")

	fi, err := os.Stat(p)
	if err != nil || !fi.IsDir() {
		fmt.Println("error: ", args[0], " is not Direcotry path or not Exist")
		os.Exit(1)
	}
	fi, err = os.Stat(args[1])
	if err != nil || !fi.IsDir() {
		fmt.Println("error: ", args[1], " is not Direcotry path")
		os.Exit(1)
	}
	err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err := PicDownloader(path, args[1]); err != nil {
				fmt.Println("error: SOME DOWNLOAD FAILED ", path)
			}
		}
		return nil
	})
}
