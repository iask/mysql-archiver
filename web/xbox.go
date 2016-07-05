package web

import (
	. "archiver/utils"
	"encoding/json"
	"fmt"
)

type XboxTree struct {
	Tree []XboxNode `json:"tree"`
}

type XboxNode struct {
	Data     string     `json:"data"`
	Open     bool       `json:"open"`
	Path     string     `json:"path"`
	Type     string     `json:"type"`
	Children []XboxNode `json:"children"`
}

type Tags struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func NewXboxTree() *XboxTree {
	return &XboxTree{}
}

func (x *XboxTree) Get(username string) ([]Tags, error) {
	var t Tags
	var tags []Tags

	if len(username) <= 0 {
		return tags, fmt.Errorf("user name is empty")
	}
	url := fmt.Sprintf("%s/%s/tree", WEB.XboxUrl, username)
	retJson, err := HttpGet(url, nil, "", "")
	//fmt.Printf("%s\n", retJson)
	err = json.Unmarshal(retJson, &x)
	if err != nil {
		return tags, err
	}
	//fmt.Println(x)
	for _, xm := range x.Tree {
		for _, dba := range xm.Children {
			for _, pro := range dba.Children {
				for _, v := range pro.Children {
					//fmt.Println(t.Data, t.Path)
					t.Name = v.Data
					t.Tag = v.Path
					tags = append(tags, t)
				}
			}
		}
	}

	return tags, nil
}
