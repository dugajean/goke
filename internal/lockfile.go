package internal

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path"
)

type singleProjectJSON map[string]int64
type lockFileJSON map[string]singleProjectJSON

type Lockfile struct {
	Files []string
	JSON  lockFileJSON
}

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

func (l *Lockfile) GetCurrentProject() singleProjectJSON {
	cwd, _ := os.Getwd()
	return l.JSON[cwd]
}

func (l *Lockfile) UpdateTimestampsForFiles(files []string) {
	lockfileMap := l.prepareMap(files)
	cwd, _ := os.Getwd()
	l.JSON[cwd] = lockfileMap

	for f := range l.JSON[cwd] {
		l.JSON[cwd][f] = lockfileMap[f]
	}

	l.generateLockfile(false)
}

func (l *Lockfile) generateLockfile(initialLockfile bool) {
	contents := l.JSON
	if initialLockfile {
		lockfileMap := l.prepareMap(l.Files)
		cwd, _ := os.Getwd()
		contents = lockFileJSON{cwd: lockfileMap}
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

func (l *Lockfile) prepareMap(files []string) singleProjectJSON {
	lockfileMapCh := make(chan singleProjectJSON)
	go l.getFileModifiedMap(files, lockfileMapCh)
	return <-lockfileMapCh
}

func (l *Lockfile) getFileModifiedMap(files []string, ch chan singleProjectJSON) {
	lockfileMap := make(singleProjectJSON)

	for _, f := range files {
		fo, err := os.Stat(f)

		if err != nil {
			log.Fatalln(err)
		}

		lockfileMap[f] = fo.ModTime().Unix()
	}

	ch <- lockfileMap
}

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

func (l *Lockfile) getLockfilePath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".goke"), nil
}
