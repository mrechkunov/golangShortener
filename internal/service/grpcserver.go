package service

import (
	"context"
	"errors"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	pb "github.com/mrechkunov/golangShortener.git/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerServer struct {
	pb.UnimplementedShortenerServer
}

// ListUserURLs return to user all urls where user is creator
func (g *ShortenerServer) ListUserURLs(ctx context.Context, e *emptypb.Empty) (out *pb.UserURLsResponse, err error) {
	// читаем метаданные и извлекаем токен
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no metadata found")
	}
	// берем uid из токена
	var uid uint32
	if values := md["authorization"]; len(values) > 0 {
		token := values[0]
		uid, err = cryptoauth.GetIDFromCookie(token)
		if err != nil {
			logger.Log.Infoln("no ID in metadata")
			return nil, status.Error(codes.InvalidArgument, "error while get ID in metadata")
		}
	}
	// Выбрать из хранилища все записи с uid
	responseBatch := repository.GetStorage().GetDataByUID(uid)
	// добавляем всем префикс базового адреса
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	var result []*pb.URLData
	for _, rb := range responseBatch {
		urlData := &pb.URLData{}
		urlData.SetOriginalUrl(rb.OriginalURL)
		urlData.SetShortUrl(baseResultAdress + "/" + rb.ShortURL)
		result = append(result, urlData)
	}
	out.SetUrl(result)
	return out, nil
}

// ExplandURLS return original url if short url is exist in DB and not set as deleted
func (g *ShortenerServer) ExpandURL(ctx context.Context, in *pb.URLExpandRequest) (out *pb.URLExpandResponse, err error) {
	shortURL := in.GetId()
	shortURL = shortURL[1:]
	longURL, isFound := repository.GetStorage().GetData(shortURL)
	if repository.GetStorage().IsDeleted(shortURL) {
		err = errors.New("short URL is deleted")
		return nil, err
	}
	if !isFound {
		err = errors.New("short URL not found")
		return nil, err
	}
	out.SetResult(longURL)
	return out, nil
}

// ShortenURL shorting url from json and insert it in DB
func (g *ShortenerServer) ShortenURL(ctx context.Context, in *pb.URLShortenRequest) (out *pb.URLShortenResponse, err error) {

	return
}
