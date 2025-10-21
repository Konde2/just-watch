package handlers

import (
	"context"
	"metadata-service/pb"
	"metadata-service/repository"
	"time"
)

type MetadataService struct {
	pb.UnimplementedMetadataServiceServer
	repo *repository.VideoRepository
}

func NewMetadataService(repo *repository.VideoRepository) *MetadataService {
	return &MetadataService{repo: repo}
}

func (s *MetadataService) CreateVideo(ctx context.Context, req *pb.VideoRequest) (*pb.VideoResponse, error) {
	video := &repository.Video{
		Title:    req.Title,
		OwnerID:  req.OwnerId,
		MinIOKey: req.MinioKey,
		Status:   req.Status,
	}
	err := s.repo.Create(video)
	if err != nil {
		return nil, err
	}
	return &pb.VideoResponse{
		Id:        video.ID.Hex(),
		Title:     video.Title,
		OwnerId:   video.OwnerID,
		MinioKey:  video.MinIOKey,
		Status:    video.Status,
		CreatedAt: video.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MetadataService) UpdateVideoStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.VideoResponse, error) {
	err := s.repo.UpdateStatus(req.VideoId, req.Status)
	if err != nil {
		return nil, err
	}
	video, err := s.repo.FindByID(req.VideoId)
	if err != nil {
		return nil, err
	}
	return &pb.VideoResponse{
		Id:        video.ID.Hex(),
		Title:     video.Title,
		OwnerId:   video.OwnerID,
		MinioKey:  video.MinIOKey,
		Status:    video.Status,
		CreatedAt: video.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MetadataService) GetVideo(ctx context.Context, req *pb.VideoIdRequest) (*pb.VideoResponse, error) {
	video, err := s.repo.FindByID(req.VideoId)
	if err != nil {
		return nil, err
	}
	return &pb.VideoResponse{
		Id:        video.ID.Hex(),
		Title:     video.Title,
		OwnerId:   video.OwnerID,
		MinioKey:  video.MinIOKey,
		Status:    video.Status,
		CreatedAt: video.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MetadataService) ListVideos(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	videos, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	var res []*pb.VideoResponse
	for _, v := range videos {
		res = append(res, &pb.VideoResponse{
			Id:        v.ID.Hex(),
			Title:     v.Title,
			OwnerId:   v.OwnerID,
			MinioKey:  v.MinIOKey,
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}
	return &pb.ListResponse{Videos: res}, nil
}
