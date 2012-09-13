# Chuckie

![Chuckie Finster](http://upload.wikimedia.org/wikipedia/en/c/ce/Charles_crandall_chuckie_finster_jr.jpg)

The Party Cloud keeps world archives in S3 for games with persistant storage. Minecraft is the canonical example of this. This app is a high-performance part of the API that fetches these archives and serves them up as Zip files no matter what format they were stored as. Currently supports `.tar.gz` and `.tar` archives.

## Usage
    
    npm install
    foreman start
