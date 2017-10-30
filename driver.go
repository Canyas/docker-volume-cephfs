package main

import (
	utils 	"./utils"
	cephfs	"./cephfs"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"fmt"
	"os"
	"errors"
)


type cephFSDriver struct { volume.Driver
	defaultPath	string
	volumes		map[string]*cephfs.Volume
}

/**

 */
func newCephFSDriver( defaultPath  string  ) (cephFSDriver, error) {
	return cephFSDriver{
		defaultPath: defaultPath,
		volumes:	nil,
	}, nil
}

func (d cephFSDriver ) Create( r volume.CreateRequest ) error {
	logrus.Info("Create Called ", r.Name, " ", r.Options)
	defer logrus.Info("Create End")

	cvol := &cephfs.Volume{
		Name: 		r.Name,
		Path: 		nil,
		Subpath:	nil,
		DataPool: 	nil,
		MetaPool: 	nil,
	}

	// Process Options
	for key, val := range r.Options {
		switch key {
			case "datapool":
				cvol.DataPool = val
			case "metapool":
				cvol.MetaPool = val
			case "name":
				cvol.Name = val
			case "path":
				cvol.Path = val
			case "subpath":
				cvol.Subpath = val
		}
	}

	// Validate required options
	if(len(cvol.Name) == 0) {
		//Required options must be set
		return errors.New("You have to specify all required options. (Required options: name)")
	}

	// Process empty options
	if(len(cvol.Path) == 0) {
		cvol.Path = d.defaultPath
	}
	if(len(cvol.Subpath) == 0) {
		cvol.Subpath = cvol.Name
	}

	volumes, err := utils.GetCephFsVolumes()
	if(err != nil) {
		logrus.Error(err.Error())
		return err
	}

	//TODO: Check if filesystem already exists

	//TODO: Create new filesystem if it doesn't exist

	//TODO: Mount filesystem

	//TODO: Check if volume already exists

	//TODO: Create new volume if it doesn't exist

	//TODO: Mount Volume

	return nil
}

func( d cephFSDriver ) List() (volume.ListResponse, error) {
	logrus.Info("List Called ")
	defer logrus.Info("List End")

	//TODO: Convert volumes

	return volume.ListResponse {}, nil
}


func( d cephFSDriver ) Get( r volume.GetRequest ) (*volume.GetResponse, error) {
	logrus.Info("Get Called ", r.Name)
	defer logrus.Info("Get End")

	//TODO: Get ceph volume by name

	//TODO: Gen mountpoint


	return &volume.GetResponse{Volume: &volume.Volume{
		Name:       r.Name,
		Mountpoint: nil,
		Status:     make(map[string]interface{}),
	}}, nil
}

func( d cephFSDriver ) Remove( r volume.RemoveRequest ) error {
	logrus.Info("Remove Called ", r.Name)
	defer logrus.Info("Remove End")

	//TODO: Update ceph volumes

	//TODO: Get ceph volume by name

	//TODO: Mount filesystem

	//TODO: Delete volume directory

	//TODO: Remove volume from array

	//TODO: Unmount filesystem

	return errors.New("error Remove NIJ")
}

func( d cephFSDriver ) Path( r volume.PathRequest ) (volume.PathResponse, error) {
	logrus.Info("Path Called ", r.Name)
	defer logrus.Info("Path End")

	//TODO: Get ceph volume by name

	//TODO: Gen mountpoint
	
	return volume.PathResponse{}, nil
}


func (d cephFSDriver ) Mount( r volume.MountRequest ) (*volume.MountResponse, error) {
	logrus.Info("Mount Called ",r.ID," ", r.Name)
	defer logrus.Info("Mount End")

	//TODO: Get volume by name

	//TODO: Mount volume

	m := fmt.Sprintf("%s/%s",d.defaultPath, r.Name)
	if( ! IsDirectory(m) ) {
		return nil, errors.New(fmt.Sprintf(" %s is not a directory ", m))
	};
	return &volume.MountResponse{ Mountpoint: m}, nil
}

func (d cephFSDriver ) Unmount( r volume.UnmountRequest ) error {
	logrus.Info("Unmount Called ", r.ID, " ", r.Name)
	defer logrus.Info("Unmount End")

	//TODO: Get volume by name

	//TODO: Unmount volume

	return errors.New("error NIJ")
}
func (d cephFSDriver ) Capabilities() *volume.CapabilitiesResponse {
	logrus.Info("Capabilities Called")
	defer logrus.Info("Capabilities End")
	
	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "global",
		},
	}
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path);
	if  err != nil {
		return false
	}
	return fileInfo.IsDir()
}