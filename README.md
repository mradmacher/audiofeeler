# audiofeeler

## Development examples

Running jekyll server:
```
jekyll serve -s ../afdata/example.com/ -d ../afdata/example.com/_site/
```

Building site with jekyll:
```
jekyll build -s ../afdata/example.com/ -d ../afdata/example.com/_site/
```

Deploying it to ftp server:
```
npx ftp-deploy --server example.com --username *** --password *** --local-dir ../afdata/example.com/_site/ --server-dir ./
```
