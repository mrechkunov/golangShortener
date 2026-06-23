package service

import (
	"context"
	"errors"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	pb "github.com/mrechkunov/golangShortener.git/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServiceServer
}

// ListUserURLs return to user all urls where user is creator
func (g GRPCServer) ListUserURLs(ctx context.Context) (out *pb.UserURLsResponse, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no metadata found")
	}
	var uid uint32
	if values := md["authorization"]; len(values) > 0 {
		token := values[0]
		uid, err = cryptoauth.GetIDFromCookie(token)
		if err != nil {
			logger.Log.Infoln("no ID in metadata")
			return nil, status.Error(codes.InvalidArgument, "error while get ID in metadata")
		}
	}
	// TODO:как задать в метаданные сгенерированный новый токен ?

	// Выбрать из хранилища все записи с uid
	responseBatch := repository.GetStorage().GetDataByUID(uid)
	// добавляем всем префикс базового адреса
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	var result []model.ResponseDataBatchByCookie
	for _, rb := range responseBatch {
		rb.ShortURL = baseResultAdress + "/" + rb.ShortURL

		result = append(result, rb)
	}
	return
}

// ExplandURLS return original url if short url is exist in DB and not set as deleted
func (g GRPCServer) ExpandURL(ctx context.Context, in *pb.URLExpandRequest) (out *pb.URLExpandResponse, err error) {
	shortURL := in.GetId()
	shortURL = shortURL[1:]
	longURL, isFound := repository.GetStorage().GetData(shortURL)
	if repository.GetStorage().IsDeleted(shortURL) {
		err = errors.New("short URL is deleted")
		out.SetResult("")
		return
	}
	if !isFound {
		err = errors.New("short URL not found")
		return
	}
	out.SetResult(longURL)
	return
}

// ShortenURL shorting url from json and insert it in DB
func (g GRPCServer) ShortenURL(ctx context.Context, in *pb.URLShortenRequest) (out *pb.URLShortenResponse) {

	return
}
