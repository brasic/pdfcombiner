pdfcombiner
===========

This is an HTTP endpoint that downloads a list of PDFs from Amazon S3,
combines them using `cpdf`, uploads the combined file, and POSTs the job
status to a provided callback URL.  The format of the request to `/` should
look like:

```json
{
  "bucket_name": "somebucket",
  "employer_id": 123,
  "doc_list": [
    "1.pdf",
    "2.pdf"
  ],
  "callback": "http://mycallbackurl.com/combination_result/12345"
}
```

The server will immediately respond either with:

    HTTP/1.1 200 OK
    {"response":"ok"}

or

    HTTP/1.1 400 Bad Request
    {"response":"invalid params"}

and begin processing the file.  When work is complete, the provided
callback URL will recieve a POST with a JSON body similar to:

```json
{
  "success": true,
  "combined_file": "path/to/combined/file.pdf",
  "job": {
    "bucket_name": "somebucket"
    "employer_id": 123,
    "doc_list": [
      "realfile.pdf",
      "nonexistent_file"
    ],
    "downloaded": [
      "realfile.pdf"
    ],
    "callback": "http://mycallbackurl.com/combination_result/12345"
    "errors": {
      "nonexistent_file": "The specified key does not exist."
    }
  }
}
```

`"success"` is true if at least one file downloaded successfully.
`"combined_file"` may be `null` if `success` is false.
