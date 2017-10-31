package lib

import (
	"strings"
	"errors"
)

func GetCephFilesystems(path string) ([]Filesystem, error) {
	// Check if ceph filesystem already exists
	out, err := ShWithDefaultTimeout("ceph", "fs", "ls")
	if(err != nil) {
		return nil, errors.New(REQUEST_LIST_ERROR + err.Error())
	}

	var existingFs []Filesystem
	var index			int
	filesystems := strings.Split(out, "\n")
	for _, element := range filesystems {
		properties := strings.Split(element, ", ")
		if(len(properties) != 3) {
			return nil, InternalError(errors.New(PROCESSING_LIST_ERROR))
		}

		existingFs = append(existingFs, Filesystem{})
		index = len(existingFs)-1
		existingFs[index].Path = path
		for _, property := range properties {
			value := strings.Split(property, ": ")
			if(len(value) != 2) {
				return nil, InternalError(errors.New(PROCESSING_LIST_ERROR))
			}

			switch (value[0]) {
			case "name":
				existingFs[index].Name = value[1]
			case "metapool":
				existingFs[index].MetaPool = value[1]
			case "data pools":
				existingFs[index].DataPool = value[1]
			}
		}
	}

	return existingFs, nil
}

func GetCephPools() ([]string, error) {
	out, err := ShWithDefaultTimeout("ceph", "osd", "pool", "ls")
	if(err != nil) {
		err = errors.New(REQUEST_POOLS_ERROR+err.Error())
		return nil, err
	}

	pools := strings.Split(out, "\n")
	if(len(pools) == 0) {
		err = errors.New(PROCESSING_POOLS_ERROR)
		return nil, err
	}

	return pools, nil
}

func ExistsCephPools(names... string) (bool, error) {
	pools, err := GetCephPools()
	if(err != nil) {
		return false, err
	}

	for _, elem := range names {
		if(!existsCephPool(pools, elem)) {
			return false, nil
		}
	}

	return true, nil
}

func existsCephPool(pools []string, name string) bool {

	for _, pool := range pools {
		if(pool == name) {
			return true
		}
	}

	return false
}

