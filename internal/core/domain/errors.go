package domain

import "errors"

var (
	// ErrInternalError is the error we return when something has gone wrong our end.
	ErrInternalError = errors.New("an internal error occurred")
	// ErrForbidden is an error that is returned when iconik returns a 403.
	ErrForbidden = errors.New("please check your app id and auth token are correct")
	// Err401Search is an error that is returned when user doesn't have correct permissions to search.
	Err401Search = errors.New("you do not have the correct permissions to access that collection")
	// Err401GetAsset is an error that is returned when user doesn't have correct permissions to get the asset.
	Err401GetAsset = errors.New("you do not have the correct permissions to get that asset")
	// Err401UpdateAsset is an error that is returned when user doesn't have correct permissions to update the asset.
	Err401UpdateAsset = errors.New("you do not have the correct permissions to update that asset")
	// Err401UpdateMetadataAsset is an error that is returned when user doesn't have correct permissions to update the metadata in the asset.
	Err401UpdateMetadataAsset = errors.New("you do not have the correct permissions to update the metadata in that asset")
	// Err401GetMetadataView is an error that is returned when user doesn't have correct permissions to get the metadata view.
	Err401GetMetadataView = errors.New("you do not have the correct permissions to get that metadata view")
	// Err401Collection is an error that is returned when user doesn't have correct permissions to access collection.
	Err401Collection = errors.New("you do not have the correct permissions to access that collection")
	// Err401CollectionContents is an error that is returned when user doesn't have correct permissions to access collection contents.
	Err401CollectionContents = errors.New("you do not have the correct permissions to access the collection contents")
)
