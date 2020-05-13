# Transcribe Service

# Installation

`go run ./services/qda_transcribe/transcribe.go --help` shows a list of options

```
  Local Execution options:
  --local-execution, -x        If set, execute the service locally once.
  --local-additional-info, -a  Additional information on local execution. Otherwise ignored.
  --local-file-name, -f        media file name if executed locally, Otherwise ignored.
```

The services uses [Amazon AWS transcribe](https://aws.amazon.com/de/transcribe/) services. You need your credentials properly installed


## API Documentation

## DocumentRequest
- `file_name` *might* contain a file name
- `document` *must* contain the audio file
- `service_name` *must* contain "`transcribe`"
- `options` 

    *must* contain 

    * the media format in form as mime-type, for example :"`Content-Type:audio/mpeg`" or "`Content-Type:audio/wav`". Currently on this two formats are supported.  All other options are ignored. Parameters to the MIMETypes are currently ignored

    *can* contain

    * `Transcript-Details: true | false` to request additional information on the transcription. Additional information will be returned in `TransOutput` (or stdout if called from command line)