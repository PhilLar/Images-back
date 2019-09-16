Images-back
=====

`images-back` is the webservice to store and process images.

Link to the web-service
------------
https://images-back2.herokuapp.com/

Installation
------------

Install `images-back` from sources, by running:

```sh
git clone https://github.com/PhilLar/Images-back.git
cd images-back
go install ./cmd/images-back
```
Database
--------
Connect to database using `psql` and run this script to create db and user:
```sql
CREATE DATABASE imagesapp;
CREATE USER images WITH PASSWORD 'secret';
GRANT ALL PRIVILEGES ON DATABASE "imagesapp" to images;
ALTER DATABASE imagesapp OWNER TO images;
```

Usage
-----
You can run it:
```sh
export DATABASE_URL="postgres://images:secret@localhost/imagesapp?sslmode=disable"
images-back
```
to add image to db:
```sh
curl -F 'file=@pic.jpg' -F 'title=mytitle' https://images-back2.herokuapp.com/files
```
to list all images in db:
```sh
curl https://images-back2.herokuapp.com/images
```

Contribute
----------
- Issue Tracker: https://github.com/PhilLar/images-back/issues
- Source Code: https://github.com/PhilLar/images-back

License
--------
[WTFPL 2.0](https://wtfpl2.com/)