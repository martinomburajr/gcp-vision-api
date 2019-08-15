package vision

import (
	"cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"os"
)

// [END imports]




// [START vision_text_detection_pdf_gcs]



// DetectAsyncDocumentURI performs Optical Character Recognition (OCR) on a
// PDF file stored in GCS.
func DetectAsyncDocumentURI(ctx context.Context, gcsSourceURI,
	gcsDestinationURI string) (*vision.AsyncBatchAnnotateFilesOperation, error) {
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	request := &visionpb.AsyncBatchAnnotateFilesRequest{
		Requests: []*visionpb.AsyncAnnotateFileRequest{
			{
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
					},
				},
				InputConfig: &visionpb.InputConfig{
					GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
					// Supported MimeTypes are: "application/pdf" and "image/tiff".
					MimeType: "application/pdf",
				},
				OutputConfig: &visionpb.OutputConfig{
					GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
					// How many pages should be grouped into each json output file.
					BatchSize: 100,
				},
			},
		},
	}

	return client.AsyncBatchAnnotateFiles(ctx, request)
}








// [END vision_text_detection_pdf_gcs]
