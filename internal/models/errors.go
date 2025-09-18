package models

import (
	"errors"
)

// define a new error handler for our application
// Why to use a custom error handler instead of using the default error handler provided by go / a standard library / or a third-party library ?
// To abstract the process error handling from main program . This keeps application related data separated from the data used for error handling
//⭐ For instance , this handler is used for sql errors , if no such entry was found in the database . In that case , we could have returned ErrNoRows error provided by the sql package . BUt there could be a possibility that sql / DB query related data could accidentally get exposed into HTTP layer , which we don't want to happen . This is where custom error handler(i.e ErrNoRecord) comes into the picture
var ErrNoRecord = errors.New("models: no matching record found")