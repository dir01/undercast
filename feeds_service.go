package undercast

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jbub/podcasts"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"time"
)

type Episode struct {
	ID       string
	MediaId  string
	MediaURL string
	Date     time.Time
}

type FeedsRepository interface {
	InsertEpisode(ctx context.Context, episode *Episode) error
	ListEpisodes(ctx context.Context) (episodes []Episode, err error)
}

type FeedsMediaService interface {
	GetMedia(ctx context.Context, mediaId string) (*Media, error)
	OnMediaUploaded(onMediaUploaded func(media *Media))
}

type FeedsStorage interface {
	StoreData(ctx context.Context, data io.ReadSeeker, filename string) (url string, err error)
}

type FeedsService struct {
	mediaService FeedsMediaService
	repository   FeedsRepository
	storage      FeedsStorage
}

func (service *FeedsService) Run() {
	service.mediaService.OnMediaUploaded(func(media *Media) {
		episode := &Episode{ID: uuid.NewV4().String(), MediaId: media.ID, MediaURL: media.Url, Date: time.Now()}
		if err := service.repository.InsertEpisode(context.TODO(), episode); err != nil {
			log.Printf("Error while creating episode by media: %s", err)
		}
		go service.publishFeed(context.TODO())
	})
}

type CreateEpisodeRequest struct {
	MediaId string `json:"mediaId"`
}

func (service *FeedsService) CreateEpisode(ctx context.Context, req CreateEpisodeRequest) error {
	media, err := service.mediaService.GetMedia(ctx, req.MediaId)
	if err != nil {
		return err
	}
	episode := &Episode{ID: uuid.NewV4().String(), MediaId: media.ID, MediaURL: media.Url, Date: time.Now()}
	if err = service.repository.InsertEpisode(ctx, episode); err != nil {
		return err
	}
	if err = service.publishFeed(ctx); err != nil {
		return err
	}
	return nil
}

func (service *FeedsService) publishFeed(ctx context.Context) error {
	episodes, err := service.repository.ListEpisodes(ctx)
	if err != nil {
		return err
	}
	p := &podcasts.Podcast{
		Title: "Audiobooks",
		//Description: "This is my very simple podcast.",
		//Language:    "EN",
		//Link:        "http://www.example-podcast.com/my-podcast",
		//Copyright:   "2015 My podcast copyright",
	}

	for _, e := range episodes {
		p.AddItem(&podcasts.Item{
			Title:    "Episode " + e.ID,
			GUID:     e.ID,
			PubDate:  podcasts.NewPubDate(e.Date),
			Duration: podcasts.NewDuration(time.Second * 230),
			Enclosure: &podcasts.Enclosure{
				URL:    e.MediaURL,
				Length: "12312",
				Type:   "MP3",
			},
		})
	}

	feed, err := p.Feed(
	//podcasts.Author("Author Name"),
	//podcasts.Block,
	//podcasts.Explicit,
	//podcasts.Complete,
	//podcasts.NewFeedURL("http://www.example-podcast.com/new-feed-url"),
	//podcasts.Subtitle("This is my very simple podcast subtitle."),
	//podcasts.Summary("This is my very simple podcast summary."),
	//podcasts.Owner("Podcast Owner", "owner@example-podcast.com"),
	//podcasts.Image("http://www.example-podcast.com/my-podcast.jpg"),
	)
	b := &bytes.Buffer{}
	feed.Write(b)
	str := b.String()
	fmt.Println(str)
	url, err := service.storage.StoreData(ctx, bytes.NewReader(b.Bytes()), "feed.xml")
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}
