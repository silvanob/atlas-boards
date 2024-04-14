package main

import (
	"slices"
	"strings"
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

type CardStorageSlice struct {
	cards []card
}

func (c card) String() string {
	return strings.Join([]string{"{", "Title:", c.title, "Content:", c.content, "}"}, " ")
}

func (s *CardStorageSlice) Modify(cardInput card) card {
	return card{title: "", content: ""}
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

func (s *CardStorageSlice) Add(cardInput card) card {
	s.cards = append(s.cards, cardInput)
	return cardInput
}

func (s *CardStorageSlice) Get(title string) card {
	for _, x := range s.cards {
		if x.title == title {
			return x
		}
	}
	return card{}
}

func (s *CardStorageSlice) List() []card {
	return s.cards
}

func NewCardStorageSlice() *CardStorageSlice {
	return &CardStorageSlice{cards: make([]card, 0)}
}
