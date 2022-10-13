# awsebs

## overview
a simple go command to query aws ebs volumes, snapshot, and kubernetes persistent volume.
get unencrypted persistent volumes from kubernetes, and run commands to encrypt them.

## usage
go run main.go 
```sh
AWS_PROFILE=something
AWS_REGION=us-east-1
# encrypt volume
go run main.go encryptv -v aws://{us-east-1a}/vol-xxxxxxxxxxx
# get volume info
go run main.go volumes get -i <volume id>
# create snapshot for volume
go run main.go snapshots create -i <volume id>  
# get snapshot 
go run main.go snapshots get -i <snapshot id>  
# copy snapshot
go run main.go snapshots copy -i <snapshot id>
# create volume from snapshot
go run main.go volumes create -i <snapshot id> -a us-east-1a
# get persisitent volumes, if multiple context, can add arg -c contextname
go run main.go pvs get
```
