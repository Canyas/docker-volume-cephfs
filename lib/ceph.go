package lib

import (
	"github.com/Sirupsen/logrus"

	"fmt"
	"errors"
	"strings"
)

type Volume struct {
	Name 		string
	Subpath		string
	Filesystem	Filesystem
}

type VolumeList []Volume

type Filesystem struct {
	Name 		string
	Path 		string
	DataPool 	string
	MetaPool 	string
}

func NewFilesystem(name 		string,
					path 		string,
					dataPool 	string,
					metaPool	string) (*Filesystem, error) {
	fs := Filesystem{
		Name:     name,
		Path:     path,
		DataPool: dataPool,
		MetaPool: metaPool,
	}

	exists, err := ExistsCephPools(fs.MetaPool, fs.DataPool)
	if(err != nil) {
		return nil, err
	} else if(!exists) {
		return nil, errors.New(MISSING_POOL)
	}

	out, err := ShWithDefaultTimeout("ceph", "fs", "new", fs.Name, fs.MetaPool, fs.DataPool)
	if(err != nil) {
		err = InternalError(errors.New(out))
		return nil, err
	}
	logrus.Debug(out)

	exists, err = fs.Exists()
	if(err != nil) {
		return nil, err
	}

	if(!exists) {
		return nil, InternalError(errors.New(MISSING_FILESYSTEM))
	}

	return &fs, nil
}

func ( v Volume) GetAbsolutePathForVolume() string {
	return fmt.Sprintf("%s/%s",v.Filesystem.Path, v.Subpath)
}

func (v Volume) Mount(monitor string, user string, secretfile string) error {
	out, err := ShWithDefaultTimeout("mount", "-t",
																"ceph-fuse",
																monitor+":"+v.Subpath,
																v.Filesystem.Path,
																"-o",
																"name="+user+",secretfile="+secretfile)
	if(err != nil) {
		err = InternalError(errors.New(out))
		return err
	}

	return nil
}

func (v Volume) Unmount() error {
	out, err := ShWithDefaultTimeout("umount", v.Filesystem.Path)
	if(err != nil) {
		err = InternalError(errors.New(out))
		return err
	}
	logrus.Debug(out)
	return nil
}

func (fs Filesystem) Exists() (bool, error) {
	fss, err := GetCephFilesystems("")
	if(err != nil) {
		logrus.Error(err.Error())
		return false, err
	}

	//Check if filesystem already exists
	for _, element := range fss {
		if(element.Name == fs.Name) {
			return true, nil
		}
	}

	return false, nil
}

func GetVolumes(monitor string, user string, secretfile string, path string) (VolumeList, error) {
	var vols []Volume

	fss, err := GetCephFilesystems(path)
	if(err != nil) {
		return nil, err
	}
	logrus.Debug(fss)

	for _, fs := range fss {
		vols_part, err := fs.GetVolumes(monitor, user, secretfile)
		if(err != nil) {
			return nil, err
		}

		for _, vol := range vols_part {
			vols = append(vols, vol)
		}
	}

	return vols, nil
}

func (fs Filesystem) GetVolumes(monitor string, user string, secretfile string) (VolumeList, error) {
	var vols []Volume

	vol := Volume{
		Name: "root",
		Subpath: "/",
		Filesystem: fs,
	}
	err := vol.Mount(monitor, user, secretfile)
	if(err != nil) {
		return nil, err
	}

	out, err := ShWithDefaultTimeout("ls", "-1", fs.Path)
	if(err != nil) {
		err = InternalError(errors.New(UNABLE_GET_VOLUMES+out))
		return nil, err
	}
	logrus.Debug(out)

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if(IsDirectory(fs.Path+"/"+line)) {
			vols = append(vols, Volume{
				Name: line,
				Subpath: "/"+line,
				Filesystem: fs,
			})
		}
	}
	logrus.Debug(lines)

	err = vol.Unmount()
	if(err != nil) {
		return nil, err
	}

	return vols, nil
}

func (vols VolumeList) ByName(name string) *Volume {
	for _, vol := range vols {
		logrus.Debug(vol)
		if(vol.Name == name) {
			return &vol
		}
	}
	return nil
}