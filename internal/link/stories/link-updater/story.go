package link_updater

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"

	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/scrape"
)

func New(repository repository, consumer amqpConsumer) *Story {
	return &Story{repository: repository, consumer: consumer}
}

type Story struct {
	repository repository
	consumer   amqpConsumer
}

func (s *Story) Run(ctx context.Context) error {
	ch, err := s.consumer.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to consume from queue")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			if msg == nil {
				continue
			}

			var linkMsg message
			if err := json.Unmarshal(msg.Body, &linkMsg); err != nil {
				log.Printf("failed to decode message: %v", err)
				continue
			}

			linkID := linkMsg.ID

			link, err := s.repository.FindByID(ctx, linkID)
			if err != nil {
				log.Printf("failed to fetch link from database: %v", err)
				continue
			}

			newData, err := scrape.GetLinkData(link.Url)
			if err != nil {
				log.Printf("failed to scrape link data: %v", err)
				continue
			}

			// Update link data with scraped data
			if newData.Title != "" {
				link.Title = newData.Title
			}
			if newData.Images != nil {
				link.Images = newData.Images
			}
			if newData.Tags != nil {
				link.Tags = newData.Tags
			}

			// Save updated link to database
			if err := s.repository.Update(ctx, link); err != nil {
				log.Printf("failed to update link in database: %v", err)
				continue
			}

			log.Printf("link with ID %s updated successfully", linkID)

			msg.Ack(false) // Acknowledge message processing
		}
	}
}
