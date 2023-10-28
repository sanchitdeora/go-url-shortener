# go-url-shortener
A simple URL shortener using Go and MySQL

## Description
Note: This URL shortener was created for learning and practice purpose. 

The Go URL shortener uses Gin to render a basic HTML form that accepts a full length URL. 
Upon submission, it generates a random 10 character key to build a shortened URL with the set domain name in the service (i.e. 'localhost' now).
The complete URL and shortened URL is stored in a simple mysql table.

Additionally, created a small unit test for the service package using GoMock, Mockgen and testify.