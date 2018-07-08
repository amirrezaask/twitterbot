package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
)

var (
	consumerKey       = getEnv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getEnv("TWITTER_CONSUMER_SECRET")
	accessToken       = getEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getEnv("TWITTER_ACCESS_TOKEN_SECRET")
	tracks            = getTracks("tracks")
)

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalln("env variable not exists :%v", name)
	}
	return v
}

func getTracks(name string) []string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalln("env variable not exists :%v" + name)
	}
	return strings.Split(v, ",")
}
func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := &logger{logrus.New()}
	api.SetLogger(log)

	s := api.PublicStreamFilter(url.Values{
		"track": tracks,
	})
	for t := range s.C {
		twitt, ok := t.(anaconda.Tweet)
		if !ok {
			log.Error("t is not twitt")
		}
		if twitt.RetweetedStatus == nil {
			_, err := api.Retweet(twitt.Id, false)
			if err != nil {
				log.Errorf("could not retweet :%v", err)
			}
			log.Infof("retweeted %v", twitt.Id)
		}
	}
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
