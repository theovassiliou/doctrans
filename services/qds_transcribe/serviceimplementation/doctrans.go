package serviceimplementation

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransformDocument
func (s *DtaService) TransformDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	l, sOut, sErr := Work(s, req.GetDocument(), req.GetFileName())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Received document: %s and echoing", string(req.GetDocument()))

	return &pb.TransformDocumentResponse{
		Document: []byte(l),
		Output:   sOut,
		Error:    errorS,
	}, nil
}

func (s *DtaService) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil
}

func (*DtaService) TransformPipe(ctx context.Context, req *pb.TransformPipeRequest) (*pb.TransformPipeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}

func (*DtaService) Options(context.Context, *pb.OptionsRequest) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

func (s *DtaService) GetDocTransServer() pb.GenDocTransServer {
	return s.GenDocTransServer
}

func (s *DtaService) ApplicationName() string {
	return s.AppName
}
