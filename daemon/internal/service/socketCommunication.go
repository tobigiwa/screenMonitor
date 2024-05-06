package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func treatMessage(ctx context.Context, c net.Conn) {
	for {
		select {
		case <-ctx.Done():
			c.Close()
			fmt.Println("we closed connection successfully")

		default:
			var (
				msg Message
				err error
			)

			if err = gob.NewDecoder(c).Decode(&msg); err != nil {
				log.Println("error reading message:", err)
				continue
			}

			switch msg.Endpoint {
			case "socketAlive":
				fmt.Printf("\n%+v\n\n", msg)
				c.Write([]byte("hello from the daemon"))
				continue

			case "weekStat":
				data, err := ServiceInstance.service.GetWeeklyScreenStats(repository.Active, msg.Body)
				if err != nil {
					log.Println("error:", err)
					continue

				}
				for _, value := range data {
					fmt.Println(value.Key, value.Value)
				}
			}

		}
	}
}

type Message struct {
	Endpoint string
	Body     string
}

// func (m *Message) encode() ([]byte, error) {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(m); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func (m *Message) decode(data []byte) error {
// 	buf := bytes.NewBuffer(data)
// 	if err := gob.NewDecoder(buf).Decode(m); err != nil {
// 		return err
// 	}
// 	return nil
// }
