package main

import (
	lib "./lib"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"errors"
	"os"
)


type cephFSDriver struct { volume.Driver
	defaultPath	string
	volumes		map[string]*lib.Volume
	monitor 	string
	user 		string
	secretfile	string
}

/**

 */
func newCephFSDriver( defaultPath  string, monitor string, user string, secretfile string) (cephFSDriver, error) {
	return cephFSDriver{
		defaultPath: defaultPath,
		volumes:	nil,
		monitor: 	monitor,
		user: 		user,
		secretfile: secretfile,
	}, nil
}

func (d cephFSDriver ) Create( r *volume.CreateRequest ) error {
	logrus.Info("--- Create Called ", r.Name, " ", r.Options)
	defer logrus.Info("--- Create End")

	cvol := &lib.Volume{
		Name:		r.Name,
		Subpath:	nil,
	}

	logrus.Info("Processing options ...")
	// Process Options
	for key, val := range r.Options {
		switch key {
			case "datapool":
				cvol.Filesystem.DataPool = val
			case "metapool":
				cvol.Filesystem.MetaPool = val
			case "fsname":
				cvol.Filesystem.Name = val
			case "path":
				cvol.Filesystem.Path = val
			case "subpath":
				cvol.Subpath = val
		}
	}

	// Validate required options
	if(len(cvol.Filesystem.Name) == 0) {
		//Required options must be set
		return errors.New(lib.REQUIRED_OPTIONS)
	}

	// Process empty options
	if(len(cvol.Filesystem.Path) == 0) {
		cvol.Filesystem.Path = d.defaultPath
	}
	if(len(cvol.Subpath) == 0) {
		cvol.Subpath = cvol.Name
	}

	logrus.Info("Checking filesystem ...")
	exists, err := cvol.Filesystem.Exists()
	if(err != nil) {
		logrus.Error(err.Error())
		return err
	} else if (!exists) {
		logrus.Info("Creating new filesystem ...")
		// Create new filesystem if it doesn't exist
		// Validate if all options are set to create a new filesystem
		if (len(cvol.Filesystem.DataPool) == 0 || len(cvol.Filesystem.MetaPool) == 0) {
			err := errors.New(lib.MISSING_POOL_OPTION)
			logrus.Error(err.Error())
			return err
		}

		_, err = lib.NewFilesystem(cvol.Filesystem.Name,
										cvol.Filesystem.Path,
										cvol.Filesystem.DataPool,
										cvol.Filesystem.MetaPool)

		if(err != nil) {
			logrus.Error(err.Error())
			return err
		}
	}

	logrus.Info("Mounting filesystem ...")
	// Mount filesystem
	fsvol := lib.Volume{
		Name: "root",
		Subpath: "/",
		Filesystem: cvol.Filesystem,
	}
	fsvol.Mount(d.monitor, d.user, d.secretfile)

	logrus.Info("Checking volume ...")
	// Check if volume already exists
	// Create new volume if it doesn't exist
	if(!lib.IsDirectory(cvol.Filesystem.Path+cvol.Subpath)) {
		logrus.Info("Creating new volume ...")
		err = os.MkdirAll(cvol.Filesystem.Path+cvol.Subpath, os.ModePerm)
		if(err != nil) {
			err = errors.New(lib.UNABLE_CREATE_DIR+err.Error())
			logrus.Error(err.Error())
			return err
		}
	}

	logrus.Info("Unmounting filesystem ...")
	// Unmount Filesystem
	err = fsvol.Unmount()
	if(err != nil) {
		logrus.Error(err.Error())
		return err
	}

	logrus.Info("Mounting volume ...")
	// Mount Volume
	err = cvol.Mount(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return err
	}

	return nil
}

func( d cephFSDriver ) List() (*volume.ListResponse, error) {
	logrus.Info("List Called ")
	defer logrus.Info("List End")

	logrus.Info("Getting volume list ....")
	// Get volumes
	vols, err := lib.GetVolumes(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return nil, err
	}

	logrus.Info("Converting volume list ...")
	var vvols []*volume.Volume
	// Convert volumes
	for _, vol := range vols {
		vvols = append(vvols, &volume.Volume{
									Name: vol.Name,
									Mountpoint: vol.Filesystem.Path+vol.Subpath,
									Status: nil,
								})
	}

	return &volume.ListResponse {
		Volumes: vvols,
	}, nil
}


func( d cephFSDriver ) Get( r *volume.GetRequest ) (*volume.GetResponse, error) {
	logrus.Info("Get Called ", r.Name)
	defer logrus.Info("Get End")

	logrus.Info("Getting volume by name ...")
	vols, err := lib.GetVolumes(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return nil, err
	}

	vol := vols.ByName(r.Name)
	if(vol == nil) {
		err = errors.New(lib.UNABLE_FIND_VOLUME+r.Name)
		logrus.Error(err.Error())
		return nil, err
	}

	return &volume.GetResponse{Volume: &volume.Volume{
		Name:       vol.Name,
		Mountpoint: vol.Filesystem.Path+vol.Subpath,
		Status:     make(map[string]interface{}),
	}}, nil
}

func( d cephFSDriver ) Remove( r *volume.RemoveRequest ) error {
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

func( d cephFSDriver ) Path( r *volume.PathRequest ) (*volume.PathResponse, error) {
	logrus.Info("Path Called ", r.Name)
	defer logrus.Info("Path End")

	logrus.Info("Getting volume by name ...")
	vols, err := lib.GetVolumes(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return nil, err
	}

	vol := vols.ByName(r.Name)
	if(vol == nil) {
		err = errors.New(lib.UNABLE_FIND_VOLUME+r.Name)
		logrus.Error(err.Error())
		return nil, err
	}
	
	return &volume.PathResponse{
		Mountpoint: vol.Filesystem.Path+vol.Subpath,
	}, nil
}


func (d cephFSDriver ) Mount( r *volume.MountRequest ) (*volume.MountResponse, error) {
	logrus.Info("Mount Called ",r.ID," ", r.Name)
	defer logrus.Info("Mount End")

	logrus.Info("Getting volume by name ...")
	vols, err := lib.GetVolumes(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return nil, err
	}

	vol := vols.ByName(r.Name)
	if(vol == nil) {
		err = errors.New(lib.UNABLE_FIND_VOLUME+r.Name)
		logrus.Error(err.Error())
		return nil, err
	}

	logrus.Info("Mounting ceph volume ...")
	// Mount volume
	err = vol.Mount(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return nil, err
	}

	return &volume.MountResponse{ Mountpoint: vol.Filesystem.Path+vol.Subpath}, nil
}

func (d cephFSDriver ) Unmount( r *volume.UnmountRequest ) error {
	logrus.Info("Unmount Called ", r.ID, " ", r.Name)
	defer logrus.Info("Unmount End")

	// Get volume by name
	logrus.Info("Getting volume by name ...")
	vols, err := lib.GetVolumes(d.monitor, d.user, d.secretfile)
	if(err != nil) {
		logrus.Error(err.Error())
		return err
	}

	vol := vols.ByName(r.Name)
	if(vol == nil) {
		err = errors.New(lib.UNABLE_FIND_VOLUME+r.Name)
		logrus.Error(err.Error())
		return err
	}

	logrus.Info("Unmount volume ...")
	// Unmount volume
	err = vol.Unmount()
	if (err != nil) {
		logrus.Error(err.Error())
		return err
	}

	return nil
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