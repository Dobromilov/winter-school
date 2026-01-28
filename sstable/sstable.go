package sstable

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"kvschool/internal/skiplist/ "
)

// ErrNotImplemented используется в заготовке практики второго дня.
var ErrNotImplemented = errors.New("sstable: функция не реализована")

type indexEntry struct {
	key    []byte
	offset int64
}

// Writer пишет отсортированные пары key/value (CDR) в файл.
// Формат файла должен позволять чтение без загрузки всего файла в память.
// Обычно это: [Data Block 1] [Data Block 2] ... [Sparse Index] [Footer].
type Writer struct {
	writer      io.Writer
	indexEntrys []indexEntry
	block       []byte
	offset      int64
}

func NewWriter(dest io.Writer) *Writer {
	return &Writer{
		writer:      bufio.NewWriter(dest),
		indexEntrys: make([]indexEntry, 0),
		block:       make([]byte, 0, 4096),
		offset:      0,
	}
}

// Add добавляет пару. Ключи должны быть строго возрастающими.
func (w *Writer) Add(key, value []byte) error {
	if len(w.block) == 0 {
		w.indexEntrys = append(w.indexEntrys, indexEntry{key: key, offset: w.offset})
	}

	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf, uint32(len(key)))
	w.block = append(w.block, buf...)
	w.block = append(w.block, key...)

	binary.BigEndian.PutUint32(buf, uint32(len(value)))
	w.block = append(w.block, buf...)
	w.block = append(w.block, value...)

	if len(w.block) > 4096 {
		return w.flushBlock()
	}

	return ErrNotImplemented
}

func (w *Writer) Close() error { return ErrNotImplemented }

// Reader читает SSTable с диска.
// Использует RandomAccess (io.ReaderAt) для прыжков по индексу.
type Reader struct{}

func NewReader(_ io.ReaderAt, _ int64) *Reader { return &Reader{} }

// Iterator возвращает упорядоченную итерацию по диапазону [start, end).
// Использует Sparse Index, чтобы найти нужный блок данных.
func (r *Reader) Iterator(_ []byte, _ []byte) (*Iter, error) { return nil, ErrNotImplemented }

type Iter struct {
}

func (it *Iter) Next() (key, value []byte, ok bool, err error) {

	return nil, nil, false, ErrNotImplemented
}

func (it *Iter) Close() error {
	return nil
}
