package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"bytes"
	"encoding/gob"
	"log"
)

type Message struct {
	Endpoint   string
	SliceData  []repository.KeyValuePair
	StringData string
	IntData    int
	IsError    bool
	Error      error
}

func (m *Message) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(m); err != nil {
		return err
	}
	return nil
}

type Service struct {
	store repository.IRepository
}

func (s *Service) weekStat(msg Message) Message {
	data, err := s.store.GetWeeklyScreenStats(repository.Active, msg.StringData)
	if err != nil {
		log.Println("error weekStat:", err)
		return Message{
			IsError: true,
			Error:   err,
		}
	}

	return Message{
		SliceData: data,
	}
}
