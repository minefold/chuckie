# Chuckie
# 
# (C) 2012 Mütli Corp. By Christopher Lloyd & David Newman.

http = require('http')
url  = require('url')
stream = require('stream')
zlib = require('zlib')

express = require('express')
bugsnag = require('bugsnag')
knox = require('knox')


# --

express.logger.format 'json', (env, req, res) ->
  JSON.stringify(
    event: 'request'
    at: new Date().getTime()
    method: req.method
    url: req.originalUrl
    status: res.statusCode,
    time: (new Date - req._startTime),
    length: parseInt(res.getHeader('Content-Length'), 10)
  )

s3 = knox.createClient
  key: process.env.AWS_ACCESS_KEY,
  secret: process.env.AWS_SECRET_KEY,
  bucket: process.env.BUCKET

# --

app = express()

app.configure ->
  app.set 'port', process.env.PORT
  app.use express.logger('json')
  app.use bugsnag.register(process.env.BUGSNAG)


app.get '/worlds', (req, res) ->
  url = req.query.url
  file = s3.get(url)
  file.on 'response', (f) ->
    
    # Check that the file exists in S3
    if f.statusCode isnt 200
      bugsnag.notify new Error("S3 Error"),
        status: f.statusCode
        headers: f.headers
      
      res.send(404)
      return
    
    res.set(
      'Content-Type': f.headers['content-type']
      'Content-Length': f.headers['content-length']    
      'ETag': f.headers['etag']
      'Last-Modified': f.headers['Last-Modified']
    )
    
    f.pipe(res, end: false)
  
  file.end()
  

# --

http.createServer(app).listen app.get('port')
