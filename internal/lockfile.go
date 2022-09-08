package internal

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path"
)

type singleProjectJson map[string]int64
type lockFileJson map[string]singleProjectJson

type Lockfile struct {
	Files []string
	JSON  lockFileJson
}

// Loads existing lock information generates it for the first time.
func (l *Lockfile) Bootstrap() {
	lockfilePath, err := l.getLockfilePath()
	if err != nil {
		log.Fatalln(err)
	}

	if !FileExists(lockfilePath) {
		l.generateLockfile(true)
	}

	currentLockFile, err := os.ReadFile(lockfilePath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(currentLockFile, &l.JSON)
	if err != nil {
		log.Fatalln(err)
	}
}

// Returns the lock information for the current project.
func (l *Lockfile) GetCurrentProject() singleProjectJson {
	cwd, _ := os.Getwd()
	return l.JSON[cwd]
}

// Update timestamps for files in current project.
func (l *Lockfile) UpdateTimestampsForFiles(files []string) {
	lockfileMap := l.prepareMap(files)
	cwd, _ := os.Getwd()
	l.JSON[cwd] = lockfileMap

	for f := range l.JSON[cwd] {
		l.JSON[cwd][f] = lockfileMap[f]
	}

	l.generateLockfile(false)
}

// Generate the lockfile file, or update it with new contents.
func (l *Lockfile) generateLockfile(initialLockfile bool) {
	contents := l.JSON
	if initialLockfile {
		lockfileMap := l.prepareMap(l.Files)
		cwd, _ := os.Getwd()
		contents = lockFileJson{cwd: lockfileMap}
	}

	jsonString, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	writeCh := make(chan error)
	go l.writeLockfile(jsonString, writeCh)

	if err = <-writeCh; err != nil {
		log.Fatalln(err)
	}
}

// Prepares the map used to populate individual project files.
func (l *Lockfile) prepareMap(files []string) singleProjectJson {
	lockfileMapCh := make(chan singleProjectJson)
	go l.getFileModifiedMap(files, lockfileMapCh)
	return <-lockfileMapCh
}

// Go routine used to dispatch file mtime checks in the background.
func (l *Lockfile) getFileModifiedMap(files []string, ch chan singleProjectJson) {
	lockfileMap := make(singleProjectJson)

	for _, f := range files {
		fo, err := os.Stat(f)

		if err != nil {
			log.Fatalln(err)
		}

		lockfileMap[f] = fo.ModTime().Unix()
	}

	ch <- lockfileMap
}

// Writes the lockfile into the filesystem.
func (l *Lockfile) writeLockfile(contents []byte, ch chan error) {
	gokePath, err := l.getLockfilePath()

	if err != nil {
		ch <- err
		return
	}

	err = os.WriteFile(gokePath, contents, 0644)

	if err != nil {
		ch <- err
		return
	}

	ch <- nil
}

// Returns the location of the lockfile in the system.
func (l *Lockfile) getLockfilePath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".goke"), nil
}
