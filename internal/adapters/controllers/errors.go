package controllers

import "errors"

var (
	ErrGettingRecordingGroup = errors.New("error getting recording group")
	ErrInvalidID             = errors.New("invalid ID")
	ErrStartingRecording     = errors.New("error starting recording")
	ErrStopRecording         = errors.New("error stopping recording")
	ErrInvalidDuration       = errors.New("duration must be informed")
	ErrclipRecording         = errors.New("error clipping recording")
	ErrFileNotFound          = errors.New("file not found")
	ErrGettingCameras        = errors.New("error getting available cameras")
	ErrGettingHomeDir        = errors.New("could`n get home dir")
)
