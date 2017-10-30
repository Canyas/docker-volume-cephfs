package utils

import (
	cephfs "../cephfs"

	"strings"
	"errors"
)

func GetCephFsVolumes() ([]cephfs.Volume, error) {
	// Check if ceph filesystem already exists
	out, err := ShWithDefaultTimeout("ceph", "fs", "ls")
	if(err != nil) {
		return nil, errors.New(REQUEST_LIST_ERROR + err.Error())
	}

	var existingVolumes []cephfs.Volume
	var index			int
	volumes := strings.Split(out, "\n")
	for _, element := range volumes {
		properties := strings.Split(element, ", ")
		if(len(properties) != 3) {
			return nil, InternalError(errors.New(PROCESSING_LIST_ERROR))
		}

		existingVolumes = append(existingVolumes, cephfs.Volume{})
		index = len(existingVolumes)-1
		for _, property := range properties {
			value := strings.Split(property, ": ")
			if(len(value) != 2) {
				return nil, InternalError(errors.New(PROCESSING_LIST_ERROR))
			}

			switch (value[0]) {
			case "name":
				existingVolumes[index].Name = value[1]
			case "metapool":
				existingVolumes[index].MetaPool = value[1]
			case "data pools":
				existingVolumes[index].DataPool = value[1]
			}
		}
	}

	return existingVolumes, nil
}