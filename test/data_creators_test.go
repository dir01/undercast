package server_test

import "context"

func (s *ServerSuite) dropDb() {
	db, err := s.getDatabase("test")
	s.Require().NoError(err)
	ctx := context.Background()
	err = db.Drop(ctx)
	s.Require().NoError(err)
}

type downloadOpts struct {
	Source string `bson:"source"`
}

func (s *ServerSuite) insertDownload(opts *downloadOpts) {
	db, err := s.getDatabase("test")
	s.Require().NoError(err)
	ctx := context.Background()
	_, err = db.Collection("downloads").InsertOne(ctx, map[string]string{"source": opts.Source})
	s.Require().NoError(err)
}
