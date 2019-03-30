# Setup 
In order to run the code, please be sure that you have installed Goji and mgo using command listed below.
```
go get goji.io
go get gopkg.in/mgo.v2
```
As for the dataset, please download a CSV version of non-geospatial data(https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh) and make sure to place it under test folder.
# Read CSV file and insert into mongoDB
I use default mongoDB port(27017) for this assignment. Execute reader.go under test folder using command
```
go run reader.go
```
and there should be a database called Building_data. Since the dataset is too large, I only insert the first 10000 records to the database.
Check your mongoDB with dbKoda or whatever you like and you will see all the inserted records.
# Test the APIs
Use one command line to run api.go. When it's connected to the database, there are two ways to test the APIs. \n
1.
