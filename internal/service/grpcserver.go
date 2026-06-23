package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	// читаем метаданные и извлекаем токен
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no metadata found")
	}
	// проверяем токен
	if values := md["authorization"]; len(values) > 0 {
		token := values[0]
		isExist := repository.GetStorage().IsCookieExist(token)
		if !isExist {
			logger.Log.Infoln("not authorizated user")
			return nil, status.Error(codes.InvalidArgument, "not authorizated user")
		}
	}
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
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	// проверяем токен
	var token string
	// читаем метаданные и извлекаем токен
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no metadata found")
	}
	if values := md["authorization"]; len(values) > 0 {
		token = values[0]
		isExist := repository.GetStorage().IsCookieExist(token)
		if !isExist {
			logger.Log.Infoln("not authorizated user")
			return nil, status.Error(codes.InvalidArgument, "not authorizated user")
		}
	}
	originalUrl := in.GetUrl()
	//сокращаем url
	hash := sha256.Sum256([]byte(originalUrl))
	shortstr := hex.EncodeToString(hash[:4])
	ShortURL := baseResultAdress + "/" + shortstr // 4 байта хеша = 8 символов в hex
	out.SetResult(ShortURL)
	// пишем в хранилище
	err = repository.GetStorage().SetData(shortstr, originalUrl, token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ulr is exist")
	}

	return out, nil
}
