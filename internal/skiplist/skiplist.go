package skiplist

import (
	"bytes"
	"errors"
	"math/rand"
)

// ErrNotFound означает отсутствие ключа (IMSI).
var ErrNotFound = errors.New("skiplist: ключ не найден")

// ErrNotImplemented используется в заготовке практики первого дня.
var ErrNotImplemented = errors.New("skiplist: функция не реализована")

// Iterator — упорядоченная итерация по диапазону ключей (Range Scan).
// В HLR используется для выгрузки абонентов по префиксу IMSI.
type Iterator interface {
	Next() (key, value []byte, ok bool, err error)
	Close() error
}

type scanIter struct {
	cur *Node
	end []byte
}

func (it *scanIter) Next() (key, value []byte, ok bool, err error) {
	if it.cur == nil {
		return nil, nil, false, nil
	}

	if it.end != nil && bytes.Compare(it.cur.key, it.end) >= 0 {
		it.cur = nil
		return nil, nil, false, nil
	}

	key = append([]byte(nil), it.cur.key...)
	value = append([]byte(nil), it.cur.value...)
	it.cur = it.cur.next[0]

	return key, value, true, nil
}

func (it *scanIter) Close() error {
	return nil
}

// SkipList — In-Memory движок для HLR.
// Обеспечивает O(log N) на чтение/запись и упорядоченный доступ.
//
// В практической реализации вам нужно хранить:
// - ключи/значения как []byte
// - уровни (forward pointers)
// - генератор уровней с фиксируемым seed (для детерминизма тестов)

type Node struct {
	value []byte
	key   []byte
	next  []*Node
}

type SkipList struct {
	head     *Node
	level    int
	maxlevel int
	p        float64
	rnd      *rand.Rand
}

// New создаёт SkipList. seed требуется для детерминируемых тестов (воспроизводимость поведения при ошибках).
func New(seed int64) *SkipList {
	s := &SkipList{level: 1, maxlevel: 100, p: 0.5, rnd: rand.New(rand.NewSource(seed))}
	s.head = &Node{next: make([]*Node, s.maxlevel)}
	return s
}

func (s *SkipList) randomlevel() int {
	level := 1
	for s.rnd.Float64() < s.p && level < s.maxlevel {
		level++
	}
	return level
}

func (s *SkipList) Put(key, value []byte) error {
	update := make([]*Node, s.maxlevel) // путь к месту вставки
	curr := s.head
	for i := s.maxlevel - 1; i >= 0; i-- {
		for curr.next[i] != nil && bytes.Compare(curr.next[i].key, key) < 0 {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	target := curr.next[0]

	if target != nil && bytes.Compare(curr.next[0].key, key) == 0 {
		copyVal := append([]byte(nil), value...)
		target.value = copyVal
		return nil
	}

	newLevel := s.randomlevel()

	if newLevel > s.level {
		for i := s.level; i < newLevel; i++ {
			update[i] = s.head
		}
		s.level = newLevel
	}

	node := &Node{
		key:   append([]byte(nil), key...),
		value: append([]byte(nil), value...),
		next:  make([]*Node, newLevel),
	}

	for i := 0; i < newLevel; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}

	return nil
}

func (s *SkipList) Get(key []byte) ([]byte, error) {
	curr := s.head

	for i := s.maxlevel - 1; i >= 0; i-- {
		for curr.next[i] != nil && bytes.Compare(curr.next[i].key, key) < 0 {
			curr = curr.next[i]
		}
	}

	curr = curr.next[0]

	if curr != nil && bytes.Compare(curr.key, key) == 0 {

		res := make([]byte, len(curr.value))
		copy(res, curr.value)
		return res, nil
	}

	return nil, ErrNotFound
}

func (s *SkipList) Delete(key []byte) error {
	update := make([]*Node, s.maxlevel)
	curr := s.head

	for i := s.maxlevel - 1; i >= 0; i-- {
		for curr.next[i] != nil && bytes.Compare(curr.next[i].key, key) < 0 {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	curr = curr.next[0]
	if curr == nil || bytes.Compare(curr.key, key) != 0 {
		return ErrNotFound
	}

	for i := 0; i < s.maxlevel; i++ {
		if update[i].next[i] == curr {
			update[i].next[i] = curr.next[i]
		}
	}
	return nil
}

// Scan возвращает итератор по диапазону [start, end).
// Если start == nil, считается -∞ (начало списка).
// Если end == nil, считается +∞ (конец списка).
func (s *SkipList) Scan(start, end []byte) (Iterator, error) {
	curr := s.head
	if start != nil {
		for i := s.maxlevel - 1; i >= 0; i-- {
			for curr.next[i] != nil && bytes.Compare(curr.next[i].key, start) < 0 {
				curr = curr.next[i]
			}
		}
	}

	curr = curr.next[0]

	if end != nil {
		end = append([]byte(nil), end...)
	}

	return &scanIter{
		cur: curr,
		end: end,
	}, nil
}
