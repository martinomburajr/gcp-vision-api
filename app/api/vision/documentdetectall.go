package vision

// [END imports]
//
//var _ = context.Background()
//var _ = vision.ImageAnnotatorClient{}
//var _ = os.Open
//
//// [END imports]
//const xthreads = 50 // Total number of threads to use, excluding the main() thread
//
////DetectAllFiles - This should be used with caution, it will run through all available files in the storage bucket and begin using VISION API, may be costly
//// If used irresponsibly
////var resc, errc = make(chan string), make(chan error)
//func DetectAllFiles(folderLocations []string, BayportGCSBucket string) error {
//	//limit := 50
//	allPaths := make([]string,1)
//	finalPath := "processed"
//
//	var wg0 sync.WaitGroup
//	var wg2 sync.WaitGroup
//	var ch = make(chan int, 80)
//
//	//init cloud storage
//	gcsClient, err := storage.NewClient(ctx)
//	if err != nil {
//		log.Fatalf("error creating gcsClient | %v", err)
//	}
//
//	imageAnnotatorClient, err := vision.NewImageAnnotatorClient(ctx)
//	if err != nil {
//		return err
//	}
//
//	//Get File Names
//
//	for x := range  folderLocations {
//		wg2.Add(1)
//		go func(i int) {
//			defer wg2.Done()
//			//1. GetAllFileNamesInFolder
//			log.Printf("iteration folders: %d | %s", i+1, folderLocations[i])
//			fileNamesFromBucket, err2 := GetFileNamesFromBucket(BayportGCSBucket, folderLocations[i], *gcsClient)
//			if err2 != nil  {
//				log.Printf("failed to get file names | %v")
//			}
//			//2. AppendToFinalList
//			allPaths = append(allPaths, fileNamesFromBucket...)
//		}(x)
//	}
//
//	wg2.Wait()
//	log.Printf("\nParsing %d files", len(allPaths))
//
//	inputURI := fmt.Sprintf("gs://%s",BayportGCSBucket)
//	outputURI := fmt.Sprintf( "gs://%s/%s", BayportGCSBucket, finalPath)
//
//	allPaths = allPaths[1:]
//
//	//for path := range allPaths {
//	//	go DetectAsyncDocumentURIV2(os.Stdout, fmt.Sprintf(inputURI+"/%s", path), fmt.Sprintf(outputURI+"/%s/", path), imageAnnotatorClient, gcsClient, &wg)
//	//}
//
//	wg0.Add(xthreads)
//	for i:=0; i<xthreads; i++ {
//		go func() {
//			for {
//				_, ok := <-ch
//				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
//					wg0.Done()
//					return
//				}
//				for k:=0; k < len(allPaths); k++ {
//					func(w io.Writer, gcsSourceURI, gcsDestinationURI string, client *vision.ImageAnnotatorClient, gcsClient *storage.Client){
//						request := &visionpb.AsyncBatchAnnotateFilesRequest{
//							Requests: []*visionpb.AsyncAnnotateFileRequest{
//								{
//									Features: []*visionpb.Feature{
//										{
//											Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
//										},
//									},
//									InputConfig: &visionpb.InputConfig{
//										GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
//										// Supported MimeTypes are: "application/pdf" and "image/tiff".
//										MimeType: "application/pdf",
//									},
//									OutputConfig: &visionpb.OutputConfig{
//										GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
//										// How many pages should be grouped into each json output file.
//										BatchSize: 20,
//									},
//								},
//							},
//						}
//
//						operation, err := client.AsyncBatchAnnotateFiles(ctx, request)
//						if err != nil {
//							log.Printf("ensure pdf or tiff | %v", err)
//							//errc <- err
//						}
//						resp, err := operation.Wait(ctx)
//						if err != nil {
//							log.Printf("%v", err)
//							//errc <- err
//						}
//
//						if operation.Done() {
//							log.Println("Loading file from Cloud Storage")
//
//							uri := strings.Replace(gcsDestinationURI, fmt.Sprintf("gs://%s/",BayportGCSBucket), "", -1)
//							fullPath, err1 := GetFileNameFromBucket(BayportGCSBucket, uri, *gcsClient)
//							//name := strings.Replace(fullPath, uri,"",-1)
//							if err1 != nil {
//								log.Printf("error getting file from bucket | %v", err1)
//								//errc <- err1
//							}
//
//							response, err2 := Load(BayportGCSBucket, fullPath, gcsClient)
//							if err2 != nil {
//								//errc <- err2
//								log.Printf("error loading file | %v", err2)
//							}
//
//							parsedText, err3 := RetrieveText(response)
//							if err3 != nil {
//								//errc <- err3
//								log.Printf("%v", err3)
//							}
//							gcp_storage.UploadTextToGCS() SendTextToGCS(BayportGCSBucket, parsedText, fullPath, gcsClient)
//						}
//
//						fmt.Fprintf(w, "%v", resp)
//					}(os.Stdout, fmt.Sprintf(inputURI+"/%s", allPaths[k]), fmt.Sprintf(outputURI+"/%s/", allPaths[k]), imageAnnotatorClient, gcsClient)
//				}
//			}
//		}()
//	}
//
//
//	//for i :=0; i < limit; i++ {
//	//	select {
//	//	case err := <-errc:
//	//		log.Println(err)
//	//	}
//	//}
//
//
//
//	// Now the jobs can be added to the channel, which is used as a queue
//	for i:=0; i<50; i++ {
//		ch <- i // add i to the queue
//	}
//
//	close(ch) // This tells the goroutines there's nothing else to do
//	wg0.Wait() // Wait for the threads to finish
//
//	// Close the imageAnnotatorClient when finished.
//	if err := gcsClient.Close(); err != nil {
//		log.Printf("error closing storage imageAnnotatorClient | %v", err)
//		// TODO: handle error.
//	}
//
//	log.Printf("\nPROCESSING COMPLETE!\nApprox Files Processed: %d", len(allPaths))
//	return nil
//}
//
//
//// detectAsyncDocument performs Optical Character Recognition (OCR) on a
//// PDF file stored in GCS.
//func DetectAsyncDocumentURIV2(w io.Writer, gcsSourceURI, gcsDestinationURI string, client *vision.ImageAnnotatorClient, gcsClient *storage.Client, wg *sync.WaitGroup) {
//
//
//
//}
