package db

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/spf13/viper"
)

type CheatCommand struct {
	Repo     string
	Name     string
	FileName string
	Columns  string
}

type FileDB_I interface {
	Open()
	Close()
	SaveAndClose()
	Trunc()
	TruncForRepo()
	Append()
	AppendString()

	GetAllCommands()
	FindMatchedCheats()
	FindMatchedCheatsForRepo()
}

type FileDB struct {
	Name     string
	Path     string
	_LOCK    bool
	_DB      *os.File
	_CONTENT []byte
}

func GetDatabase(name string, path string) *FileDB {
	baseDir := viper.GetString("db")

	if path == "" {
		path = fmt.Sprintf("%s/%s", baseDir, name)
	}

	return &FileDB{
		Name: name,
		Path: path,
	}
}

func (db *FileDB) Read() error {
	original_DB, err := ioutil.ReadFile(db.Path)
	if err != nil {
		return errors.New("database file could not be found")
	}
	db._CONTENT = original_DB
	db._LOCK = true
	return nil
}

func (db *FileDB) Open() (err error) {
	original_DB, err := ioutil.ReadFile(db.Path)
	if err != nil {
		fmt.Println("Creating new database ...")
	}
	db._CONTENT = original_DB

	err = ioutil.WriteFile(fmt.Sprintf("%s.update", db.Path), original_DB, 0644)
	if err != nil {
		return
	}

	_DB, err := os.OpenFile(fmt.Sprintf("%s.update", db.Path), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	db._DB = _DB
	return
}

func (db *FileDB) Close() error {
	return db._DB.Close()
}

func (db *FileDB) SaveAndClose() error {
	if db._LOCK {
		return errors.New("database is Locked")
	}
	db._DB.Write(db._CONTENT)
	db._DB.Close()

	original := db.Path
	updated := fmt.Sprintf("%s.update", db.Path)

	return os.Rename(updated, original)
}

func (db *FileDB) Trunc() error {
	if db._LOCK {
		return errors.New("database is Locked")
	}

	db._CONTENT = []byte("")
	return db._DB.Truncate(0)
}

func (db *FileDB) TruncForRepo(repoName string) error {
	if db._LOCK {
		return errors.New("database is Locked")
	}
	var re = regexp.MustCompile(fmt.Sprintf(`\[(%s)\]\s(.*)\:(.*\.md)\:(.*)`, repoName))
	newContent := re.ReplaceAll(db._CONTENT, []byte(""))
	newContent = bytes.Trim(newContent, "\n")
	db.Trunc()
	db._DB.Write(newContent)
	return nil
}

func (db *FileDB) Append(line []byte) error {
	if db._LOCK {
		return errors.New("database is Locked")
	}
	db._CONTENT = append(db._CONTENT, line...)
	return nil
}

func (db *FileDB) AppendString(line string) error {
	if db._LOCK {
		return errors.New("database is Locked")
	}
	db.Append([]byte(line))
	return nil
}

func (db *FileDB) GetAllCommands() (commands map[string][]*CheatCommand) {
	re := regexp.MustCompile(`\[(.*)\]\s(.*)\:(.*\.md)\:(.*)`)
	commands = make(map[string][]*CheatCommand)

	matches := re.FindAllStringSubmatch(string(db._CONTENT), -1)
	for _, m := range matches {
		commands[m[1]] = append(commands[m[1]], &CheatCommand{
			Repo:     m[1],
			Name:     m[2],
			FileName: m[3],
			Columns:  m[4],
		})
	}
	return
}

// FindMatchedCheats
func (db *FileDB) FindMatchedCheats(command string) (commands []*CheatCommand) {
	re := regexp.MustCompile(fmt.Sprintf(`\[(.*)\]\s(%s)\:(.*\.md)\:(.*)`, command))

	matches := re.FindAllStringSubmatch(string(db._CONTENT), -1)
	for _, m := range matches {
		commands = append(commands, &CheatCommand{
			Repo:     m[1],
			Name:     m[2],
			FileName: m[3],
			Columns:  m[4],
		})
	}
	return
}

// FindMatchedCheatsForRepo will look for repo specific commands
// listed in the indexed database and return all the matched commands
func (db *FileDB) FindMatchedCheatsForRepo(command string, repoName string) (commands []*CheatCommand) {
	re := regexp.MustCompile(fmt.Sprintf(`\[(%s)\]\s(%s)\:(.*\.md)\:(.*)`, command, repoName))

	matches := re.FindAllStringSubmatch(string(db._CONTENT), -1)
	for _, m := range matches {
		commands = append(commands, &CheatCommand{
			Repo:     m[1],
			Name:     m[2],
			FileName: m[3],
			Columns:  m[4],
		})
	}
	return
}
