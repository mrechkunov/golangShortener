package service

import (
	"errors"

	"github.com/mrechkunov/golangShortener.git/internal/repository"
	pb "github.com/mrechkunov/golangShortener.git/proto"
)

type gRPCServer struct {
	pb.UnimplementedShortenerServiceServer
}

// ListUserURLs return to user all urls where user is creator
func (g gRPCServer) ListUserURLs() (rez pb.UserURLsResponse) {

	return
}

// ExplandURLS return original url if short url is exist in DB and not set as deleted
func (g gRPCServer) ExpandURL(in pb.URLExpandRequest) (rez pb.URLExpandResponse, err error) {
	shortURL := in.GetId()
	shortURL = shortURL[1:]
	longURL, isFound := repository.GetStorage().GetData(shortURL)
	if repository.GetStorage().IsDeleted(shortURL) {
		err = errors.New("short URL is deleted")
		rez.SetResult("")
		return
	}
	if !isFound {
		err = errors.New("short URL not found")
		return
	}
	rez.SetResult(longURL)
	return
}

// ShortenURL shorting url from json and insert it in DB
func (g gRPCServer) ShortenURL(in pb.URLShortenRequest) (rez pb.URLShortenResponse) {
	in.GetUrl()

	return
}
