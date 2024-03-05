package tar

import (
	"archive/tar"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
)

const bootstrapentry = "__bootstrap__"
const bootstrapoffset = 512
const indexentry = "__index__"

var IN_MEMORY_TAR_FOLDERS = os.Getenv("DYNATRACE_IN_MEMORY_TAR_FOLDERS") == "true"

var folderCache = FolderCache{
	mu:      sync.Mutex{},
	folders: map[string]*Folder{},
}

type FolderCache struct {
	mu      sync.Mutex
	folders map[string]*Folder
}

type Folder struct {
	mu          sync.Mutex
	name        string
	index       tarIndex
	indexOffset int64
	offsetBytes []byte
}

type indexEntry struct {
	api.Stub
	Offset int64
}

type tarIndex map[string]indexEntry

func New(name string) (*Folder, bool, error) {

	tf := &Folder{name: name + ".tar", index: tarIndex{}, offsetBytes: make([]byte, 4)}

	if fileExists(tf.name) {
		return tf, true, tf.initExisting()
	}
	return tf, false, tf.initNew()
}

func NewExisting(name string) (*Folder, bool, error) {

	if IN_MEMORY_TAR_FOLDERS {
		folderCache.mu.Lock()
		folder, ok := folderCache.folders[name]
		if ok {
			folderCache.mu.Unlock()
			return folder, true, nil
		}
		folderCache.mu.Unlock()
	}

	tf := &Folder{name: name + ".tar", index: tarIndex{}, offsetBytes: make([]byte, 4)}

	if fileExists(tf.name) {
		err := tf.initExisting()
		if err != nil {
			return tf, true, err
		}

		if IN_MEMORY_TAR_FOLDERS {
			folderCache.mu.Lock()

			folder, ok := folderCache.folders[name]
			if ok {
				folderCache.mu.Unlock()
				return folder, true, nil
			}

			folderCache.folders[name] = tf
			folderCache.mu.Unlock()
		}

		return tf, true, err
	}
	return nil, false, nil
}

func (tf *Folder) initNew() error {
	file, err := os.OpenFile(tf.name, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	writer := tar.NewWriter(file)
	defer file.Close()

	tf.indexOffset = 1024
	binary.LittleEndian.PutUint32(tf.offsetBytes, uint32(tf.indexOffset))
	if err := tf.write(writer, bootstrapentry, tf.offsetBytes); err != nil {
		return err
	}
	if err := tf.writeIndex(writer); err != nil {
		return err
	}
	return nil
}

func (me *Folder) initExisting() error {
	file, err := os.OpenFile(me.name, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := tar.NewReader(file)
	header, err := reader.Next()
	if err != nil {
		return err
	}
	if header.Name != bootstrapentry {
		return errors.New("not an indexed tar file")
	}
	if header.Size != 4 {
		return errors.New("bootstrap size is not 4. file is corrupt")
	}
	if _, err := reader.Read(me.offsetBytes); err != nil && err != io.EOF {
		return err
	}
	me.indexOffset = int64(binary.LittleEndian.Uint32(me.offsetBytes))
	if _, err = file.Seek(me.indexOffset, io.SeekStart); err != nil {
		return err
	}
	data, err := me.read(file)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &me.index); err != nil {
		return err
	}
	return nil
}

func (me *Folder) Get(id string) (*api.Stub, []byte, error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	idxEntry, found := me.index[id]
	if !found {
		return nil, nil, nil
	}
	if idxEntry.Offset == 0 {
		stub := idxEntry.Stub
		return &stub, nil, nil
	}
	file, err := os.OpenFile(me.name, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	if _, err := file.Seek(idxEntry.Offset, 0); err != nil {
		return nil, nil, err
	}
	data, err := me.read(file)
	if err != nil {
		return nil, nil, err
	}
	stub := idxEntry.Stub
	return &stub, data, nil
}

func (me *Folder) List() (api.Stubs, error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	stubs := api.Stubs{}
	for _, v := range me.index {
		st := v.Stub
		stubs = append(stubs, &st)
	}
	return stubs, nil
}

func (me *Folder) Save(stub api.Stub, data []byte) error {
	me.mu.Lock()
	defer me.mu.Unlock()
	file, err := os.OpenFile(me.name, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Seek(me.indexOffset, io.SeekStart); err != nil {
		return err
	}
	writer := tar.NewWriter(file)
	if data != nil {
		me.index[stub.ID] = indexEntry{Stub: stub, Offset: me.indexOffset}
		if err = me.write(writer, stub.ID, data); err != nil {
			return err
		}
		if me.indexOffset, err = file.Seek(0, io.SeekCurrent); err != nil {
			return err
		}
		if err = me.writeIndex(writer); err != nil {
			return err
		}
		if _, err = file.Seek(bootstrapoffset, 0); err != nil {
			return err
		}
		binary.LittleEndian.PutUint32(me.offsetBytes, uint32(me.indexOffset))
		if _, err = file.Write(me.offsetBytes); err != nil {
			return nil
		}
	} else {
		me.index[stub.ID] = indexEntry{Stub: stub, Offset: 0}
		if me.indexOffset, err = file.Seek(0, io.SeekCurrent); err != nil {
			return err
		}
		if err = me.writeIndex(writer); err != nil {
			return err
		}
	}
	return nil
}

type GetStubDataFunc func(idx int) (api.Stub, []byte, bool, error)

func (me *Folder) SaveAllCallback(numItems int, getStubDataFunc GetStubDataFunc) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	file, err := os.OpenFile(me.name, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Seek(me.indexOffset, io.SeekStart); err != nil {
		return err
	}
	writer := tar.NewWriter(file)

	for idx := 0; idx < numItems; idx++ {
		stub, data, isSkipped, err := getStubDataFunc(idx)
		if err != nil {
			return err
		}

		if isSkipped {
			continue
		}

		if data != nil {

			me.index[stub.ID] = indexEntry{Stub: stub, Offset: me.indexOffset}
			if err = me.write(writer, stub.ID, data); err != nil {
				return err
			}
			if me.indexOffset, err = file.Seek(0, io.SeekCurrent); err != nil {
				return err
			}
		} else {
			me.index[stub.ID] = indexEntry{Stub: stub, Offset: 0}
		}
	}

	if _, err = file.Seek(me.indexOffset, io.SeekStart); err != nil {
		return err
	}
	if err = me.writeIndex(writer); err != nil {
		return err
	}
	if _, err = file.Seek(bootstrapoffset, 0); err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(me.offsetBytes, uint32(me.indexOffset))
	if _, err = file.Write(me.offsetBytes); err != nil {
		return nil
	}

	return nil
}

func (tf *Folder) write(writer *tar.Writer, name string, data []byte) error {
	if name != bootstrapentry {
		name = name + ".json"
	}
	if err := writer.WriteHeader(
		&tar.Header{Name: name, Size: int64(len(data))}); err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return writer.Flush()
}

func (me *Folder) writeIndex(writer *tar.Writer) error {
	defer writer.Close()
	data, err := json.MarshalIndent(me.index, "", "\t")
	if err != nil {
		panic(err)
	}
	return me.write(writer, indexentry, data)
}

func (me *Folder) read(file *os.File) ([]byte, error) {
	return me.readWith(tar.NewReader(file))
}

func (me *Folder) readWith(reader *tar.Reader) ([]byte, error) {
	header, err := reader.Next()
	if err != nil {
		return nil, err
	}
	data := make([]byte, header.Size)
	if _, err := reader.Read(data); err != nil && err != io.EOF {
		return nil, err
	}
	return data, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func (me *Folder) Delete(id string) error {
	me.mu.Lock()
	defer me.mu.Unlock()
	if _, found := me.index[id]; found {
		delete(me.index, id)
		file, err := os.OpenFile(me.name, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err = file.Seek(me.indexOffset, io.SeekStart); err != nil {
			return err
		}
		writer := tar.NewWriter(file)
		if err = me.writeIndex(writer); err != nil {
			return err
		}
		if _, err = file.Seek(bootstrapoffset, 0); err != nil {
			return err
		}
		binary.LittleEndian.PutUint32(me.offsetBytes, uint32(me.indexOffset))
		if _, err = file.Write(me.offsetBytes); err != nil {
			return nil
		}
		return nil
	}
	return nil
}
