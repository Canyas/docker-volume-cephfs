package utils

import(
	"fmt"
 	"errors"
)

const (
	REQUEST_LIST_ERROR = "Unable to request ceph volumes: "
	PROCESSING_LIST_ERROR = "Unable to convert output from command \"ceph fs ls\"."
)

func InternalError(err error) error {
	return errors.New(fmt.Sprintf("Internal error(maybe ceph version is not compatible): %s", err.Error()))
}