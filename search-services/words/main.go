package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"

	// "unicode/utf8"

	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"

	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
	words "yadro.com/course/words/words"
)

// const maxPhraseLen = 4096

type server struct {
	wordspb.UnimplementedWordsServer
	log *slog.Logger
}

type Config struct {
	Port string `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"8080"`
}

func (s *server) Ping(_ context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	s.log.Debug("got ping")
	return nil, nil
}

func (s *server) Norm(_ context.Context, in *wordspb.WordsRequest) (*wordspb.WordsReply, error) {
	s.log.Debug("got phrase", "phrease", in.Phrase)

	// size check OBSOLETE
	// if utf8.RuneCountInString(in.Phrase) > maxPhraseLen {
	// 	return nil, status.Errorf(codes.ResourceExhausted, "phrase is too long and exceed 4 KiB limit")
	// }

	// main logic
	stemmed := words.Norm(in.Phrase)

	return &wordspb.WordsReply{
		Words: stemmed,
	}, nil
}

func main() {
	var address, config string
	flag.StringVar(&address, "address", ":8080", "server address")
	flag.StringVar(&config, "config", "config.yaml", "configuration file")
	flag.Parse()

	var cfg Config

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		panic(err)
	}

	log := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	))

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("failed to listen", "error", err)
	}

	s := grpc.NewServer()
	wordspb.RegisterWordsServer(s, &server{log: log})
	reflection.Register(s)

	log.Debug("started grpc server")
	if err := s.Serve(listener); err != nil {
		log.Error("failed to serve", "errror", err)
	}

}
