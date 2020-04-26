package serviceimplementation

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"regexp"

	pb "github.com/theovassiliou/doctrans/dtaservice"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
)

// countResults describes the results of the transformation
type countResults struct {
	Bytes int
	Lines int
	Words int
}

var re *regexp.Regexp = regexp.MustCompile(`[\S]+`)

// Work returns an encoded JSON object containing the
// bytes 	count the number of bytes
// lines	count the numnber of lines
// words		count the number of words
// The Service returns  the number of lines, words, and bytes contained in the input document
func Work(input []byte, options []string) (string, []string, error) {

	b := len(input)
	l, err := counter(bytes.NewReader(input), []byte{'\n'})
	w := len(re.FindAllString(string(input), -1))

	res := &countResults{
		Bytes: b,
		Lines: l,
		Words: w,
	}
	resB, _ := json.MarshalIndent(res, "", "  ")
	log.WithFields(log.Fields{"Service": "Work"}).Infof("The result %s\n", resB)

	return string(resB), []string{}, err
}

func counter(r io.Reader, sep []byte) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], sep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

type Implementation struct {
}

// TransformDocument implements dtaservice.DTAServer
func (Implementation) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	l, sOut, sErr := Work(in.GetDocument(), in.GetOptions())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s and has lines %s", string(in.GetDocument()), l)

	return &pb.TransformDocumentResponse{
		TransDocument: []byte(l),
		TransOutput:   sOut,
		Error:         errorS,
	}, nil
}

// ListServices list available services provided by this implementation
func (Implementation) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil

}

func (Implementation) TransformPipe(ctx context.Context, req *pb.TransformPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}

// ApplicationName returns the name of the service application
func ApplicationName() string {
	// return serviceName
	return "count"
}
