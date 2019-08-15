package gcp_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
)

// UploadToGCS
func UploadTextToGCS(ctx context.Context, text string, storageBucket, bucketPath string, client *storage.Client) {
	wc := client.Bucket(storageBucket).Object(bucketPath).NewWriter(ctx)
	wc.ContentType = "text/plain"

	if _, err := wc.Write([]byte(text)); err != nil {
		log.Printf("createFile: unable to write data to bucket %q, file %q: %v", storageBucket, bucketPath, err)
		return
	}
	if err := wc.Close(); err != nil {
		log.Printf("createFile: unable to close bucket %q, file %q: %v", storageBucket, bucketPath, err)
		return
	}
}

// GetFileNameFromBucket gets a file name from a bucket
func GetFileNameFromBucket(ctx context.Context, gcsBucket, path string, client *storage.Client) (string, error) {
	prefix := path
	delim := "/"

	it := client.Bucket(gcsBucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "",fmt.Errorf("error reading from bucket e.g.prefix | %v", err)
		}
		return attrs.Name,nil
	}
	return "", nil
}

func GetFileNamesFromBucket(ctx context.Context, gcsBucket, path string, client *storage.Client) ([]string, error) {
	prefix := path
	delim := "/"
	filenames := make([]string, 0)

	it := client.Bucket(gcsBucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil ,fmt.Errorf("error reading from bucket e.g.prefix | %v", err)
		}
		filenames = append(filenames, attrs.Name)
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("no items in file location")
	}
	return filenames,nil
}

// GetObjectHandlesFromBucket returns a list of object handles from a given bucket and path
func GetObjectHandlesFromBucket(ctx context.Context, gcsBucket, path string, client *storage.Client) ([]*storage.ObjectAttrs, error) {
	prefix := path
	delim := "/"

	files := make([]*storage.ObjectAttrs, 0)

	it := client.Bucket(gcsBucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil ,fmt.Errorf("error reading from bucket e.g.prefix | %v", err)
		}
		files = append(files, attrs)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no items in file location")
	}
	return files,nil
}

//Download retrieves the item from cloud storage. The MIME type in Cloud storage must be application/json otherwise
//The  application will crash
func Download(ctx context.Context, gcsBucket, path string, client *storage.Client) (*DocumentDetectResponse, error) {
	reader, err := client.Bucket(gcsBucket).Object(path).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving file from gcsbucket | %v", err)
	}

	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading file from gcsbucket | %v", err)
	}

	var response DocumentDetectResponse
	err1 := json.Unmarshal(bytes, &response)
	if err1 != nil {
		return nil, fmt.Errorf("error unmarshaling json | %v", err)
	}

	return &response, nil
}

//RetrieveText parses the JSON and retrieves only the relevant text elements
func RetrieveText(response *DocumentDetectResponse) (string, error) {
	top := "Source File: " + response.InputConfig.GcsSource.URI
	str := "\n----NEW-PAGE----\n"
	for _, item := range response.Responses {
		str = item.FullTextAnnotation.Text + str
	}
	return top + "\n\n" + str, nil
}


//DocumentDetectResponse represents the output upon analyzing a document. The fulltext annotations are omitted to save space.
type DocumentDetectResponse struct {
	InputConfig struct {
		GcsSource struct {
			URI string `json:"uri"`
		} `json:"gcsSource"`
		MimeType string `json:"mimeType"`
	} `json:"inputConfig"`
	Responses []struct {
		Context struct {
			PageNumber int    `json:"pageNumber"`
			URI        string `json:"uri"`
		} `json:"context"`
		FullTextAnnotation struct {
			Pages []struct {
				Blocks []struct {
					BlockType   string `json:"blockType"`
					BoundingBox struct {
						NormalizedVertices []struct {
							X float64 `json:"x"`
							Y float64 `json:"y"`
						} `json:"normalizedVertices"`
					} `json:"boundingBox"`
					Confidence float64 `json:"confidence"`
					Paragraphs []struct {
						BoundingBox struct {
							NormalizedVertices []struct {
								X float64 `json:"x"`
								Y float64 `json:"y"`
							} `json:"normalizedVertices"`
						} `json:"boundingBox"`
						Confidence float64 `json:"confidence"`
						Words      []struct {
							BoundingBox struct {
								NormalizedVertices []struct {
									X float64 `json:"x"`
									Y float64 `json:"y"`
								} `json:"normalizedVertices"`
							} `json:"boundingBox"`
							Confidence float64 `json:"confidence"`
							Property   struct {
								DetectedLanguages []struct {
									LanguageCode string `json:"languageCode"`
								} `json:"detectedLanguages"`
							} `json:"property"`
							Symbols []struct {
								Confidence float64 `json:"confidence"`
								Property   struct {
									DetectedLanguages []struct {
										LanguageCode string `json:"languageCode"`
									} `json:"detectedLanguages"`
								} `json:"property"`
								Text string `json:"text"`
							} `json:"symbols"`
						} `json:"words"`
					} `json:"paragraphs"`
				} `json:"blocks"`
				Height   int `json:"height"`
				Property struct {
					DetectedLanguages []struct {
						Confidence   float64 `json:"confidence"`
						LanguageCode string  `json:"languageCode"`
					} `json:"detectedLanguages"`
				} `json:"property"`
				Width int `json:"width"`
			} `json:"-"`
			Text string `json:"text"`
		} `json:"fullTextAnnotation"`
	} `json:"responses"`
}
