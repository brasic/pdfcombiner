package job

import (
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	awsRegion              = aws.USEast
	uploadedFilePermission = s3.Private
)

// Get retrieves the requested document, either from S3 or by decoding the
// embedded `Data` attribute of the Document.
func (j *Job) Get(doc Document) (docContent []byte, err error) {
	if testmode.IsEnabled() {
		return
	}
	switch doc.Data {
	case "":
		return j.download(doc)
	default:
		return decodeEmbeddedData(doc.Data)
	}
}

// UploadCombinedFile sends a file to the job's CombinedKey on S3.
func (j *Job) UploadCombinedFile(localPath string) (err error) {
	if testmode.IsEnabled() {
		return
	}
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		j.AddError(j.CombinedKey, err)
		return
	}
	err = j.bucket.Put(j.CombinedKey, content, "application/pdf", uploadedFilePermission)
	if err != nil {
		j.AddError(j.CombinedKey, err)
	}
	j.uploadComplete = true
	return
}

// Connect to AWS and add the handle to the Job object.
func (j *Job) connect() (err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	s := s3.New(auth, awsRegion)
	j.bucket = s.Bucket(j.BucketName)
	return
}

func (j *Job) download(doc Document) (content []byte, err error) {
	remotePath := doc.s3Path()
	return j.bucket.Get(remotePath)
}