package router

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/martinomburajr/gcp-vision-api/app"
	"github.com/martinomburajr/gcp-vision-api/app/api/gcp_storage"
	"github.com/martinomburajr/gcp-vision-api/app/api/vision"
	"log"
	"net/http"
	"strings"
	"sync"
)

// OCRHandler handles POST requests that represent a basic request to perform OCR.
// A GCS file path and output path must be provided
func OCRHandler(gCPVisionAPIServer *GCPVisionAPIServer) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *AppError {

		outputParam := r.URL.Query().Get("output")
		if outputParam == "json" || outputParam == ""|| outputParam == "text" {
			err := fmt.Errorf("invalid url param")
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("invalid param. output = 'json', 'text | %v", err), err)
		}

		ctx := context.Background()
		var bodyInfo FileInfoJSON

		// DECODE HTTP REQUEST BODY
		err := json.NewDecoder(r.Body).Decode(&bodyInfo)
		if err != nil {
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("error reading request | %v", err), err)
		}

		// GET SOURCE AND DESTINATION BUCKET NAMES
		sourceBucketName, err := GetBucketName(bodyInfo.InputURI)
		if err != nil {
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("error parsing bucket name | %v", err), err)
		}

		outputBucketName, err := GetBucketName(bodyInfo.OutputURI)
		if outputBucketName == "" {
			outputBucketName = sourceBucketName
		}
		if err != nil {
			log.Printf("annotateFilesOperation bucket not set - defaulting to input bucket name")
		}
		log.Printf("Using Input Storage Bucket: %s", sourceBucketName)

		inputPath, err := ExtractGCSPath(bodyInfo.InputURI)
		outputPath := ""
		if bodyInfo.OutputPath == "" {
			outputPath = "visionapi/output/" + inputPath
		} else {
			outputPath = bodyInfo.OutputPath + "/" + inputPath
		}

		outputFullURI := fmt.Sprintf("gs://%s/%s", outputBucketName, outputPath)
		sourceFullURI := fmt.Sprintf("%s", bodyInfo.InputURI)

		//Performs the document scanning
		annotateFilesOperation, err := vision.DetectAsyncDocumentURI(ctx, sourceFullURI, outputFullURI)
		if err != nil {
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("error performing document detection | %v", err), err)
		}

		_, err = annotateFilesOperation.Wait(ctx)
		if err != nil {
			return AppErrorf(http.StatusInternalServerError, fmt.Sprintf("error waiting on files operation | %v", err),
				err)
		}
		if annotateFilesOperation.Done() {
			log.Println("completed annotation of files in VisionAPI")
		}

		// OUTPUT TO CLIENT
		outputMsg := fmt.Sprintf("Successfully wrote to %s", outputFullURI)
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(outputMsg))
		if err != nil {
			log.Printf("error writing to client %s", err.Error())
		}
		return nil
	}

}

// OCRBucketDirHandler performs the vision API operation on the entire bucket or dir provided
func OCRBucketDirHandler(gCPVisionAPIServer *GCPVisionAPIServer) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *AppError {

		outputParam := r.URL.Query().Get("format")
		if outputParam != "json" && outputParam != "" && outputParam != "text" {
			err := fmt.Errorf("invalid url param")
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("invalid param. output = 'json', 'text | %v", err), err)
		}

		ctx := context.Background()
		var bodyInfo FileInfoJSON

		// DECODE HTTP REQUEST BODY
		err := json.NewDecoder(r.Body).Decode(&bodyInfo)
		if err != nil {
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("error reading request | %v", err), err)
		}

		// GET SOURCE AND DESTINATION BUCKET NAMES
		sourceBucketName, err := GetBucketName(bodyInfo.InputURI)
		if err != nil {
			return AppErrorf(http.StatusBadRequest, fmt.Sprintf("error parsing bucket name | %v", err), err)
		}

		outputBucketName, err := GetBucketName(bodyInfo.OutputURI)
		if outputBucketName == "" {
			outputBucketName = sourceBucketName
		}
		if err != nil {
			log.Printf("annotateFilesOperation bucket not set - defaulting to input bucket name")
		}
		log.Printf("Using Input Storage Bucket: %s", sourceBucketName)

		inputPath, err := ExtractGCSPath(bodyInfo.InputURI)
		outputPath := ""
		if bodyInfo.OutputPath == "" {
			outputPath = "visionapi/output"
		} else {
			outputPath = bodyInfo.OutputPath
		}

		sourceBucketPath, err := ExtractGCSPath(bodyInfo.InputURI)
		if err != nil {
			return AppErrorf(http.StatusInternalServerError, "", err)
		}

		fileDestinationURI := fmt.Sprintf("gs://%s/%s/%s", outputBucketName, outputPath, inputPath)
		sourceFullURI := fmt.Sprintf("%s", bodyInfo.InputURI)

		///////////////////////////////////////// READ FILES FROM BUCKET ///////////////////////////////////////////////

		// Get Files In Bucket
		objectsHandles, err := gcp_storage.GetObjectHandlesFromBucket(ctx, sourceBucketName, sourceBucketPath,
			gCPVisionAPIServer.StorageClient)
		if err != nil {
			return AppErrorf(http.StatusInternalServerError, fmt.Sprintf("error getting files from buckets names | %v",
				err), err)
		}
		var wg sync.WaitGroup
		var errChan chan *AppError
		errChan = make(chan *AppError)
		var fileDoneChan = make(chan string)
		var completeChan = make(chan bool)

		// Run Vision API on Files
		for _, object := range objectsHandles {
			wg.Add(1)

			go func(object *storage.ObjectAttrs, outputBucketName, outputPath string, wg *sync.WaitGroup) {
				defer wg.Done()

				fileName := object.Name
				fileDestinationURI = fmt.Sprintf("gs://%s/%s/%s", outputBucketName, outputPath, fileName)
				fileSourceURI := fmt.Sprintf("%s/%s", sourceFullURI, fileName)
				/////////////////////////////////////////// CALL VISION API
				annotateFilesOperation, err := vision.DetectAsyncDocumentURI(ctx, fileSourceURI, fileDestinationURI)
				if err != nil {
					errChan <- AppErrorf(http.StatusBadRequest,
						fmt.Sprintf("error performing document detection | %v",
							err), err)
				}

				_, err = annotateFilesOperation.Wait(ctx)
				if err != nil {
					errChan <- AppErrorf(http.StatusInternalServerError, fmt.Sprintf("error waiting on files operation | %v", err),
						err)
				}

				if annotateFilesOperation.Done() {
					msg := fmt.Sprintf("Completed: File %s | Output: %s", fileName, fileDestinationURI)
					fileDoneChan <- msg
					log.Print("LOG: " + msg)
				}
			}(object, outputBucketName, outputPath, &wg)

		}
		wg.Wait()
		log.Print(" Waiting! ")
		completeChan <- true

		sb := strings.Builder{}
		//for {
			log.Print(" In for! ")
			select {
			case str := <-fileDoneChan:
				sb.WriteString(str + "\n")
				log.Print(str)
			case err := <-errChan:
				errStr := fmt.Sprintf("error: %d | %s | %s", err.Code, err.Message, err.Error.Error())
				sb.WriteString(errStr)
			case <-completeChan:
				log.Print(" COMPLETE INVOKED! ")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(sb.String()))
				break
			default:
				fmt.Print("Waiting on Channels")
			}
		//}
		return nil
	}
}

// VisionAPI writes back to cloud storage to the given path

//uri, err := ExtractGCSPath(outputFullURI)
//if err != nil {
//	return AppErrorf(http.StatusInternalServerError, fmt.Sprintf("error constructing output path | %v",
//		err), err)
//}
//fullPath, err := gcp_storage.GetFileNameFromBucket(ctx, outputPath, uri, gCPVisionAPIServer.StorageClient)
//if err != nil {
//	return AppErrorf(http.StatusInternalServerError, fmt.Sprintf("error getting file | %v", err), err)
//}

//response, err := gcp_storage.Download(ctx, outputBucketName, outputPath, gCPVisionAPIServer.StorageClient)
//if err != nil {
//	log.Printf("%v", err)
//}
//
//parsedText, err2 := gcp_storage.RetrieveText(response)
//if err2 != nil {
//	log.Printf("%v", err)
//}

//data, err3 := json.Marshal(struct {
//	Text string `json:"text"`
//}{Text: parsedText})
//if err3 != nil {
//	log.Printf("Error marshalling output json | %v",  err3)
//}

//Upload(w,r, parsedText)

//newPath := strings.Replace(bucketPath, "json", "txt", -1)
//log.Printf("Creating file /%s/%s\n", outputPath, newPath)
//gcp_storage.UploadTextToGCS(ctx, parsedText, fullPath, outputFullURI, gCPVisionAPIServer.StorageClient)

// FileInfoJSON represents the POST from the client detailing the URI for the client. 
// Once processed the items will be sent to the OutputURI. Both must be specified
type FileInfoJSON struct {
	// OutputURI represents the GCS URI to drop off the artifacts after processing.
	OutputURI string `json:"outputUri"`
	// OutputPath represents the internal bucket folder to use
	OutputPath string `json:"outputPath"`
	// InputURI represents the GCS URI of the data to be used
	InputURI string `json:"inputUri"`
}

// GetBucketName tries to obtain the bucket name.
// Bucket names must contain the gs:// prefix to distinguish them as being GCP type buckets.
func GetBucketName(uri string) (string, error) {
	if uri == "" {
		return "", errors.New("error getting bucket name from outputuri | probably empty")
	}
	if len(uri) < 7 {
		return "", errors.New("error invalid size uri | probably small in size")
	}
	if !strings.HasPrefix(uri, "gs://") {
		return "", errors.New("error invalid uri, must start with gs://")
	}
	split := strings.Split(uri[5:], "")
	ans := ""
	for _, v := range split {
		if v != "/" {
			ans = ans + v
		} else {
			break
		}
	}

	return strings.Trim(ans, ""), nil
}

// ExtractGCSPath returns the path after the bucket name and gs:// scheme.
func ExtractGCSPath(uri string) (string, error) {
	if uri == "" {
		return "", fmt.Errorf("uri cannot be empty")
	}
	if !strings.Contains(uri, app.GCSScheme) {
		return "", fmt.Errorf("uri must contain GCS")
	}

	bucketName, err := GetBucketName(uri)
	if err != nil {
		return "", err
	}

	stringWithoutGCSScheme := strings.Replace(uri, app.GCSScheme, "", -1)
	ans := strings.Replace(stringWithoutGCSScheme, bucketName, "", -1)

	if ans == "" {
		return "", nil
	}
	if ans == "/" {
		return "", nil
	}
	return ans[1:], nil
}
