package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kr/beanstalk"
	"gopkg.in/mgo.v2"
)

func MongoDB(mongoAddr string) gin.HandlerFunc {
	conn, err := mgo.Dial(mongoAddr)

	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("dbConn", conn)
		c.Next()
	}
}

func Beanstalkd(beanstalkAddr string) gin.HandlerFunc {
	conn, err := beanstalk.Dial("tcp", beanstalkAddr)

	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("queueConn", conn)
		c.Next()
	}
}
