package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type RawStorage interface {
	Read(key Key) ([]byte, error)
	Exists(key Key) bool
	Write(key Key, content []byte) error
	Delete(key Key) error
	List(key KindKey) ([]Key, error)
	Checksum(key Key) (string, error)
	Format(key Key) Format
	Dir() string
}

func NewDefaultRawStorage(dir string) RawStorage {
	return &DefaultRawStorage{
		dir: dir,
	}
}

// TODO: This is a test: Make DefaultRawStorage MappedRawStorage compliant
var _ MappedRawStorage = &DefaultRawStorage{}

type DefaultRawStorage struct {
	dir string
}

func (r *DefaultRawStorage) realPath(key AnyKey) string {
	var file string

	switch key.(type) {
	case KindKey:
	// KindKeys get no special treatment
	case Key:
		// Keys get the metadata filename added to the returned path
		file = constants.METADATA
	default:
		panic(fmt.Sprintf("invalid key type received: %T", key))
	}

	return path.Join(r.dir, key.String(), file)
}

func (r *DefaultRawStorage) Read(key Key) ([]byte, error) {
	return ioutil.ReadFile(r.realPath(key))
}

func (r *DefaultRawStorage) Exists(key Key) bool {
	return util.FileExists(r.realPath(key))
}

func (r *DefaultRawStorage) Write(key Key, content []byte) error {
	file := r.realPath(key)

	// Create the underlying directories if they do not exist already
	if !r.Exists(key) {
		if err := os.MkdirAll(path.Dir(file), constants.DATA_DIR_PERM); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(file, content, 0644)
}

func (r *DefaultRawStorage) Delete(key Key) error {
	return os.RemoveAll(path.Dir(r.realPath(key)))
}

func (r *DefaultRawStorage) List(key KindKey) ([]Key, error) {
	entries, err := ioutil.ReadDir(r.realPath(key))
	if err != nil {
		return nil, err
	}

	result := make([]Key, 0, len(entries))
	for _, entry := range entries {
		result = append(result, NewKey(key.Kind, meta.UID(entry.Name())))
	}

	return result, nil
}

// This returns the modification time as a UnixNano string
// If the file doesn't exist, return blank
func (r *DefaultRawStorage) Checksum(key Key) (s string, err error) {
	var fi os.FileInfo

	if r.Exists(key) {
		if fi, err = os.Stat(r.realPath(key)); err == nil {
			s = strconv.FormatInt(fi.ModTime().UnixNano(), 10)
		}
	}

	return
}

func (r *DefaultRawStorage) Format(key Key) Format {
	return FormatJSON // The DefaultRawStorage always uses JSON
}

func (r *DefaultRawStorage) Dir() string {
	return r.dir
}

func (r *DefaultRawStorage) AddMapping(key Key, path string) {} // Stub

func (r *DefaultRawStorage) GetMapping(p string) (key Key, err error) {
	splitDir := strings.Split(filepath.Clean(r.dir), string(os.PathSeparator))
	splitPath := strings.Split(filepath.Clean(p), string(os.PathSeparator))

	if len(splitPath) < len(splitDir)+2 {
		err = fmt.Errorf("path not long enough: %s", p)
		return
	}

	for i := 0; i < len(splitDir); i++ {
		if splitDir[i] != splitPath[i] {
			err = fmt.Errorf("path has wrong base: %s", p)
			return
		}
	}

	return ParseKey(path.Join(splitPath[len(splitDir)], splitPath[len(splitDir)+1]))
}

func (r *DefaultRawStorage) RemoveMapping(key Key) {} // Stub
