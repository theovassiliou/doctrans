package serviceimplementation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	pb "github.com/theovassiliou/doctrans/dtaservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DtaService struct {
	pb.IDocTransServer
	pb.GenDocTransServer
	pb.UnimplementedDTAServerServer
}

// CountResults describes the results of the transformation
type CountResults struct {
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
func Work(s *DtaService, input []byte, options []string) (string, []string, error) {

	b := len(input)
	l, err := counter(bytes.NewReader(input), []byte{'\n'})
	w := len(re.FindAllString(string(input), -1))

	res := &CountResults{
		Bytes: b,
		Lines: l,
		Words: w,
	}
	resB, _ := json.MarshalIndent(res, "", "  ")
	log.WithFields(log.Fields{"Service": "Work"}).Infof("The result %s\n", resB)

	return string(resB), []string{}, err
}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.DocumentRequest) (*pb.TransformDocumentResponse, error) {
	l, sOut, sErr := Work(s, in.GetDocument(), in.GetOptions())
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

func (s *DtaService) ListServices(ctx context.Context, req *pb.ListServiceRequest) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil

}

func (*DtaService) TransformPipe(context.Context, *pb.TransformPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}
func (*DtaService) Options(context.Context, *pb.OptionsRequest) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

func (s *DtaService) GetDocTransServer() pb.GenDocTransServer {
	return s.GenDocTransServer
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

func check(e error) {
	if e != nil {
		log.Errorln(e)
	}
}

func ExecuteWorkerLocally(s DtaService, fileName string) {
	if fileName == "" {
		log.Errorln("No fileName on local executing provided. Aborting.")
		return
	}

	dat, err := ioutil.ReadFile(fileName)
	check(err)

	transDoc, _, err := Work(&s, dat, nil)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(1)
	}

	fmt.Println(transDoc)
}
