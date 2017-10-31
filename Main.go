package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"

	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"os"
)

const (
	cephfsId      = "_cephfs"
	socketAddress = "/run/docker/plugins/cephfs.sock"
	logfile		  = "/var/log/docker-volume-cephfs.log"
)

var (
	defaultPath = filepath.Join(volume.DefaultDockerRootDirectory, cephfsId)
	monitor = ""
	user = "admin.clinet"
	secretfile = "/etc/ceph/admin.secretfile"
)

func main() {

	var Usage = func() {
		fmt.Println("LATIN CAPITAL LETTER AA Ꜳ ꜳ")
		fmt.Println("   LAO VOWEL SIGN AA າ ຳ")
	}


	file, err := setupLogging()
	if(err != nil) {
		fmt.Println("Logging not possible.")
	} else {
		defer shutdownLogging(file)
	}
	EnvironmentConfiguration()

	var setup = func() {
		fmt.Printf("Path %s\n", defaultPath)
	}

	Usage()
	setup()

	fstype := LookupFileSystemType(defaultPath)
	if !strings.Contains(fstype, "ceph") {
		log.Print("Warning CePH filesystem not found at ", defaultPath, " found ", fstype)
	}

	driver, err := newCephFSDriver(defaultPath, monitor, user, secretfile)
	if err != nil {
		return
	}
	h := volume.NewHandler(driver)

	fmt.Printf("Listening on %s\n", socketAddress)
	fmt.Println(h.ServeUnix(socketAddress, 1))
}

func LookupFileSystemType(path string) string {
	out, err := exec.Command("df", "--no-sync", "--output=fstype", path).Output()

	if err != nil {
		log.Fatal("Unable to read df output", err)
	}

	fstype := strings.Split(string(out), "\n")[1]
	return fstype
}

func EnvironmentConfiguration() {
	path := os.Getenv("DEFAULT_PATH")
	monitor = os.Getenv("DEFAULT_MONITOR")
	user = os.Getenv("CEPH_USER")
	secretfile = os.Getenv("CEPH_SECRETFILE")

	if(len(defaultPath) > 0 ) {
		defaultPath = path
	}

	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "3":
		logrus.SetLevel(logrus.DebugLevel)
		break;
	case "2":
		logrus.SetLevel(logrus.InfoLevel)
		break;
	case "1":
		logrus.SetLevel(logrus.WarnLevel)
		break;
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
}

// setupLogging attempts to log to a file, otherwise stderr
func setupLogging() (*os.File, error) {
	// use date, time and filename for log output
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("docker-volume-cephfs")

	// setup logfile - path is set from logfileDir and pluginName
	logFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	// check if we can write to directory - otherwise just log to stderr?
		if os.IsPermission(err) {
			log.Printf("WARN: logging fallback to STDERR: %v", err)
		} else {
			// some other, more extreme system error
			return nil, err
		}
	} else {
		logrus.Info(fmt.Sprintf("INFO: setting log file: %s", logfile))
		log.SetOutput(logFile)
		return logFile, nil
	}
	return nil, nil
}

func shutdownLogging(logFile *os.File) {
	// flush and close the file
	if logFile != nil {
		log.Println("INFO: closing log file")
		logFile.Sync()
		logFile.Close()
	}
}