package service

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/vladimirvivien/gowfs"
	"os"
	"path"
	"strconv"
	"time"
)

type WebHdfsClient struct {
	Addr string
	User string
}

var fs *gowfs.FileSystem


func (w* WebHdfsClient) Connect () {
	log.Infof("connect webhdfs namenode -> %s, user -> %s", w.Addr, w.User)
	fileSystem, err := gowfs.NewFileSystem(gowfs.Configuration{Addr: w.Addr, User: w.User})
	if err != nil{
		log.Error("connect hdfs error ...")
	}

	log.Info("connect hdfs success ... ")

	fs = fileSystem

}

func (w * WebHdfsClient) GetFileStatus (p string) (gowfs.FileStatus, error) {
	fileStatus, err := fs.GetFileStatus(gowfs.Path{Name:p})
	if err != nil {
		log.Info("get home directory err", err)
	}
	log.Info("fileStatus -> ", fileStatus)
	return fileStatus, err
}

func (w * WebHdfsClient) Ls( hdfsPath string) {
	stats, err := fs.ListStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Unable to list paths: ", err)
	}
	log.Printf("Found %d file(s) at %s\n", len(stats), hdfsPath)
	for _, stat := range stats {
		fmt.Printf(
			"%-11s %3s %s\t%s\t%11d %20v %s\n",
			FormatFileMode(stat.Permission, stat.Type),
			FormatReplication(stat.Replication, stat.Type),
			stat.Owner,
			stat.Group,
			stat.Length,
			FormatModTime(stat.ModificationTime),
			stat.PathSuffix)
	}
}

func FormatFileMode(webfsPerm string, fileType string) string {
	perm, _ := strconv.ParseInt(webfsPerm, 8, 16)
	fm := os.FileMode(perm)
	if fileType == "DIRECTORY" {
		fm = fm | os.ModeDir
	}
	return fm.String()
}

func FormatReplication(rep int64, fileType string) string {
	repStr := strconv.FormatInt(rep, 8)
	if fileType == "DIRECTORY" {
		repStr = "-"
	}
	return repStr
}

func FormatModTime(modTime int64) string {
	modTimeAdj := time.Unix((modTime / 1000), 0) // adjusted for Java Calendar in millis.
	return modTimeAdj.Format("2006-01-02 15:04:05")
}

func (w * WebHdfsClient) CreateDir(hdfsPath string) {
	path := gowfs.Path{Name: hdfsPath}
	ok, err := fs.MkDirs(path, 0744)
	if err != nil || !ok {
		log.Fatal("Unable to create test directory ", hdfsPath, ":", err)
	}
	log.Info("HDFS Path ", path.Name, " created.")
	w.Ls(path.Name)
}

func (w * WebHdfsClient) UploadFile(local, hdfsPath string) string {
	file, err := os.Open(local)
	if err != nil {
		log.Fatal("Unable to find local test file: ", err)
	}
	stat, _ := file.Stat()
	if stat.Mode().IsDir() {
		log.Fatal("Data file expected, directory found.")
	}
	log.Infof("Test file ", stat.Name(), " found.")

	shell := gowfs.FsShell{FileSystem: fs}
	log.Infof("Sending file ", file.Name(), " to HDFS location ", hdfsPath)
	ok, err := shell.Put(file.Name(), hdfsPath, true)
	if err != nil || !ok {
		log.Fatal("Failed during test file upload: ", err)
	}
	_, fileName := path.Split(file.Name())
	log.Infof("File ", fileName, " Copied OK.")
	remoteFile := hdfsPath + "/" + fileName
	w.Ls(remoteFile)

	return remoteFile
}

func (w * WebHdfsClient) RnameRemoteFile(oldName, newName string) {
	_, err := fs.Rename(gowfs.Path{Name: oldName}, gowfs.Path{Name: newName})
	if err != nil {
		log.Fatal("Unable to rename remote file ", oldName, " to ", newName)
	}
	log.Infof("HDFS file ", oldName, " renamed to ", newName)
	w.Ls(newName)
}

func (w * WebHdfsClient) ChangeOwner( hdfsPath string) {
	shell := gowfs.FsShell{FileSystem: fs}
	_, err := shell.Chown([]string{hdfsPath}, "owner2")
	if err != nil {
		log.Fatal("Chown failed for ", hdfsPath, ": ", err.Error())
	}
	stat, err := fs.GetFileStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Unable to validate chown() operation: ", err.Error())
	}
	if stat.Owner == "owner2" {
		log.Infof("Chown for ", hdfsPath, " OK ")
		w.Ls( hdfsPath)
	} else {
		log.Fatal("Chown() failed.")
	}
}

func (w * WebHdfsClient) ChangeGroup( hdfsPath string) {
	shell := gowfs.FsShell{FileSystem: fs}
	_, err := shell.Chgrp([]string{hdfsPath}, "superduper")
	if err != nil {
		log.Fatal("Chgrp failed for ", hdfsPath, ": ", err.Error())
	}
	stat, err := fs.GetFileStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Unable to validate chgrp() operation: ", err.Error())
	}
	if stat.Group == "superduper" {
		log.Infof("Chgrp for ", hdfsPath, " OK ")
		w.Ls( hdfsPath)
	} else {
		log.Fatal("Chgrp() failed.")
	}
}

func (w * WebHdfsClient) ChangeMod(hdfsPath string) {
	shell := gowfs.FsShell{FileSystem: fs}
	_, err := shell.Chmod([]string{hdfsPath}, 0744)
	if err != nil {
		log.Fatal("Chmod() failed for ", hdfsPath, ": ", err.Error())
	}
	stat, err := fs.GetFileStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Unable to validate Chmod() operation: ", err.Error())
	}
	if stat.Permission == "744" {
		log.Infof("Chmod for ", hdfsPath, " OK ")
		w.Ls( hdfsPath)
	} else {
		log.Fatal("Chmod() failed.")
	}
}

func (w * WebHdfsClient)AppendToRemoteFile(localFile, hdfsPath string) {
	stat, err := fs.GetFileStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Unable to get file info for ", hdfsPath, ":", err.Error())
	}
	shell := gowfs.FsShell{FileSystem: fs}
	_, err = shell.AppendToFile([]string{localFile}, hdfsPath)
	if err != nil {
		log.Fatal("AppendToFile() failed: ", err.Error())
	}

	stat2, err := fs.GetFileStatus(gowfs.Path{Name: hdfsPath})
	if err != nil {
		log.Fatal("Something went wrong, unable to get file info:", err.Error())
	}
	if stat2.Length > stat.Length {
		log.Infof("AppendToFile() for ", hdfsPath, " OK.")
		w.Ls( hdfsPath)
	} else {
		log.Fatal("AppendToFile failed. File size for ", hdfsPath, " expected to be larger.")
	}
}

func (w * WebHdfsClient)MoveRemoteFileLocal(remoteFile string) {
	log.Infof("Moving Remote file!!")
	shell := gowfs.FsShell{FileSystem: fs}
	remotePath, fileName := path.Split(remoteFile)
	_, err := shell.MoveToLocal(remoteFile, fileName)
	if err != nil {
		log.Fatal("MoveToLocal() failed: ", err.Error())
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("MoveToLocal() - local file can't be open. ")
	}
	defer file.Close()
	defer os.Remove(file.Name())

	_, err = fs.GetFileStatus(gowfs.Path{Name: remoteFile})
	if err == nil {
		log.Fatal("Expecing a FileNotFoundException, but file is found. ", remoteFile, ": ", err.Error())
	}
	log.Printf("Remote file %s has been removed Ok", remoteFile)
	w.Ls( remotePath)
}