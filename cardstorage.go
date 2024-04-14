package main

import "slices"

type card struct {
	title   string
	content string
}

type CardStorage struct {
	cards []card
}

func (s *CardStorage) RemoveCardByTitle(title string) {
	index := slices.IndexFunc(s.cards, func(e card) bool {
		if e.title == title {
			return true
		}
		return false
	})
	s.cards = append(s.cards[:index], s.cards[index+1:]...)
}

func (s *CardStorage) AddCard(title string, content string) {
	s.cards = append(s.cards, card{title: title, content: content})
}

func (s *CardStorage) ListCards() []card {
	return s.cards
}

func NewCardStorage() *CardStorage {
	return &CardStorage{cards: make([]card, 0)}
}
