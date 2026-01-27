package skiplist

import (
	"bytes"
	//"bytes"
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

// SkipList — In-Memory движок для HLR.
// Обеспечивает O(log N) на чтение/запись и упорядоченный доступ.
//
// В практической реализации вам нужно хранить:
// - ключи/значения как []byte
// - уровни (forward pointers)
// - генератор уровней с фиксируемым seed (для детерминизма тестов)
type SkipList struct {
	head     *Node
	level    int
	maxlevel int
	p        float64
	rnd      *rand.Rand
}

type Node struct {
	key   []byte
	value []byte
	next  []*Node // указатель на все уровни
}

func (s *SkipList) randomlevel() int {
	level := 1
	for s.rnd.Float64() < s.p && level < s.maxlevel {
		level++
	}
	return level
}

func New(seed int64) *SkipList {
	s := &SkipList{level: 1, maxlevel: 100, p: 0.5, rnd: rand.New(rand.NewSource(seed))}
	s.head = &Node{next: make([]*Node, s.maxlevel)}
	return s
}

func (s *SkipList) Put(key, value []byte) {
	curr := s.head
	update := make([]*Node, s.maxlevel)

	for i := s.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && bytes.Compare(curr.next[i].key, key) < 0 { // идем вправо до тех пор пока ключ следущего значегия не больше текущего
			curr = curr.next[i]
		}
		update[i] = curr
	}

	target := curr.next[0]
	//т.к. мы дошли и встали слева от нуной позиции

	if target != nil && bytes.Compare(target.key, key) == 0 {
		target.value = value
		return
	}

	newLevel := s.randomlevel()

	if newLevel > s.level {
		for i := s.level; i < newLevel; i++ {
			update[i] = s.head
		}
		s.level = newLevel
	}

	newNode := &Node{
		key:   key,
		value: value,
		next:  make([]*Node, newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

func (s *SkipList) Get(key []byte) ([]byte, error) {
	curr := s.head

	for i := s.level - 1; i >= 0; i-- {
		for curr.next != nil && bytes.Compare(curr.next[0].key, key) < 0 {
			curr = curr.next[i]
		}
	}

	curr = curr.next[0]

	if curr != nil && bytes.Compare(curr.key, key) == 0 {
		return curr.value, nil
	}

	return nil, ErrNotFound
}

func (s *SkipList) Delete(key []byte) {
	update := make([]*Node, s.maxlevel)
	curr := s.head
	for i := s.level - 1; i >= 0; i-- {
		for curr.next != nil && bytes.Compare(curr.next[0].key, key) < 0 {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	curr = curr.next[0]
}

// Scan возвращает итератор по диапазону [start, end).
// Если start == nil, считается -∞ (начало списка).
// Если end == nil, считается +∞ (конец списка).
func (s *SkipList) Scan(start, end []byte) (Iterator, error) {
	_ = s
	_ = start
	_ = end
	return nil, ErrNotImplemented
}
