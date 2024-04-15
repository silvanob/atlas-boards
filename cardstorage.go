package main

import (
	"slices"
	"strings"
	"sync"
)

type StorageDriver interface {
	Get(string) card
	Add(card) card
	List() []card
	Remove(string) card
	Modify(card) card
}

type card struct {
	title   string
	content string
}

type CardStorageMap struct {
	sync.RWMutex
	cardMap map[string]card
}

type CardStorageSlice struct {
	cards []card
}

func (c card) String() string {
	return strings.Join([]string{"{", "Title:", c.title, "Content:", c.content, "}"}, " ")
}

func (s *CardStorageMap) Modify(cardInput card) card {
	return card{title: "", content: ""}
}
func (s *CardStorageSlice) Modify(cardInput card) card {
	return card{title: "", content: ""}
}

func (s *CardStorageMap) Remove(title string) card {
	s.RLock()
	returnCard := s.cardMap[title]
	s.RUnlock()
	s.Lock()
	delete(s.cardMap, title)
	s.Unlock()
	return returnCard
}

func (s *CardStorageSlice) Remove(title string) card {
	index := slices.IndexFunc(s.cards, func(e card) bool {
		if e.title == title {
			return true
		}
		return false
	})
	cardReturn := s.cards[index]
	s.cards = append(s.cards[:index], s.cards[index+1:]...)

	return cardReturn

}

func (s *CardStorageMap) Add(cardInput card) card {
	s.Lock()
	s.cardMap[cardInput.title] = cardInput
	s.Unlock()
	return cardInput
}

func (s *CardStorageSlice) Add(cardInput card) card {
	s.cards = append(s.cards, cardInput)
	return cardInput
}

func (s *CardStorageMap) Get(title string) card {
	s.RLock()
	returnCard := s.cardMap[title]
	s.RUnlock()
	return returnCard
}

func (s *CardStorageSlice) Get(title string) card {
	for _, x := range s.cards {
		if x.title == title {
			return x
		}
	}
	return card{}
}

func (s *CardStorageMap) List() []card {
	s.RLock()
	cardList := make([]card, 0, len(s.cardMap))
	for _, value := range s.cardMap {
		cardList = append(cardList, value)
	}
	s.RUnlock()
	return cardList
}
func (s *CardStorageSlice) List() []card {
	return s.cards
}

func NewCardStorageSlice() *CardStorageSlice {
	return &CardStorageSlice{cards: make([]card, 0)}
}

func NewCardStorageMap() *CardStorageMap {
	return &CardStorageMap{cardMap: make(map[string]card)}
}
