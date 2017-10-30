package cephfs

import "fmt"

type Volume struct {
	Name		string
	Path		string
	Subpath		string
	DataPool	string
	MetaPool	string
}

func ( v Volume) GetAbsolutePathForVolume() string {
	return fmt.Sprintf("%s/%s",v.Path, v.Subpath)
}