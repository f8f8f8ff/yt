package dl

import "errors"

var BadMedium = errors.New("unknown medium")
var BadUrl = errors.New("couldn't get youtube id from url")
var BadPath = errors.New("couldn't get youtube id from path")
