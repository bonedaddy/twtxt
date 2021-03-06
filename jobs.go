package twtxt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var Jobs map[string]JobFactory

func init() {
	Jobs = map[string]JobFactory{
		"@every 5m": NewUpdateFeedsJob,
		"@every 1h": NewFixUserAccountsJob,
	}
}

type JobFactory func(conf *Config, store Store) cron.Job

type UpdateFeedsJob struct {
	conf *Config
	db   Store
}

func NewUpdateFeedsJob(conf *Config, db Store) cron.Job {
	return &UpdateFeedsJob{conf: conf, db: db}
}

func (job *UpdateFeedsJob) Run() {
	users, err := job.db.GetAllUsers()
	if err != nil {
		log.WithError(err).Warn("unable to get all users from database")
		return
	}

	log.Infof("updating feeds for %d users", len(users))

	sources := make(map[string]string)

	for _, user := range users {
		for u, n := range user.sources {
			sources[n] = u
		}
	}

	log.Infof("updating %d sources", len(sources))

	cache, err := LoadCache(job.conf.Data)
	if err != nil {
		log.WithError(err).Warn("error loading feed cache")
		return
	}

	cache.FetchTweets(sources)

	if err := cache.Store(job.conf.Data); err != nil {
		log.WithError(err).Warn("error saving feed cache")
	} else {
		log.Info("updated feed cache")
	}
}

type FixUserAccountsJob struct {
	conf *Config
	db   Store
}

func NewFixUserAccountsJob(conf *Config, db Store) cron.Job {
	return &FixUserAccountsJob{conf: conf, db: db}
}

func (job *FixUserAccountsJob) Run() {
	users, err := job.db.GetAllUsers()
	if err != nil {
		log.WithError(err).Warn("unable to get all users from database")
		return
	}

	for _, user := range users {
		normalizedUsername := NormalizeUsername(user.Username)

		if normalizedUsername != user.Username {
			log.Infof("migrating user account %s -> %s", user.Username, normalizedUsername)

			if err := job.db.DelUser(user.Username); err != nil {
				log.WithError(err).Errorf("error deleting old user %s", user.Username)
				return
			}

			p := filepath.Join(filepath.Join(job.conf.Data, feedsDir))

			if err := os.Rename(filepath.Join(p, user.Username), filepath.Join(p, fmt.Sprintf("%s.tmp", user.Username))); err != nil {
				log.WithError(err).Errorf("error renaming old feed for %s", user.Username)
				return
			}

			if err := os.Rename(filepath.Join(p, fmt.Sprintf("%s.tmp", user.Username)), filepath.Join(p, normalizedUsername)); err != nil {
				log.WithError(err).Errorf("error renaming new feed for %s", user.Username)
				return
			}

			// Fix Username
			user.Username = normalizedUsername

			// Fix URL
			user.URL = fmt.Sprintf(
				"%s/u/%s",
				strings.TrimSuffix(job.conf.BaseURL, "/"),
				normalizedUsername,
			)

			if err := job.db.SetUser(normalizedUsername, user); err != nil {
				log.WithError(err).Errorf("error migrating user %s", normalizedUsername)
				return
			}

			log.Infof("successfully migrated user account %s", normalizedUsername)
		}

		if user.URL == "" {
			log.Infof("fixing URL for user %s", user.Username)
			// Fix URL
			user.URL = fmt.Sprintf(
				"%s/u/%s",
				strings.TrimSuffix(job.conf.BaseURL, "/"),
				normalizedUsername,
			)

			if err := job.db.SetUser(normalizedUsername, user); err != nil {
				log.WithError(err).Errorf("error migrating user %s", normalizedUsername)
				return
			}

			log.Infof("successfully fixed URL for user %s", user.Username)
		}
	}
}
