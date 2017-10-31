package lib

import(
	"fmt"
 	"errors"
)

const (
	REQUIRED_OPTIONS = "You have to specify all required options. (Required options: fsname)"
	MISSING_POOL_OPTION = "You need to specify a Data-/Metapool to create a new Filesystem."

	MISSING_POOL = "One of the given pools doesn't exist."
	MISSING_FILESYSTEM = "Can't find newly created filesystem."

	UNABLE_CREATE_DIR = "Unable to create volume directory. Error: "
	UNABLE_GET_VOLUMES = "Unable to list all volumes. Error: "
	UNABLE_FIND_VOLUME = "Unable to find the volume. Name: "

	REQUEST_LIST_ERROR = "Unable to request ceph volumes: "
	PROCESSING_LIST_ERROR = "Unable to convert output from command \"ceph fs ls\"."
	REQUEST_POOLS_ERROR = "Unable to request ceph pools: "
	PROCESSING_POOLS_ERROR = "There are no pools."
)

func InternalError(err error) error {
	return errors.New(fmt.Sprintf("Internal error(maybe ceph version is not compatible): %s", err.Error()))
}