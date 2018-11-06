# DESCRIPTION

This command is used to update existing server in SOMA. It is a full
replace, requiring the full server record to be specified in the
command.

It is possible to mark servers as deleted using the update command
by including the optional `deleted true` keyword. It is not possible to
resurrect a deleted server by specifying `deleted false` on a deleted
server. Other fields on deleted servers can be updated.

# SYNOPSIS

```
soma server update ${serverID} name ${name} assetid ${assetID} datacenter ${locode} location ${loc} [online ${isOnline}] [deleted ${isDeleted}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
serverID | string | UUID of the server | | no
name | string | Name of the server | | no
assetID | integer | Numeric asset ID of the server | | no
locode | string | UN/Locode of the datacenter the server is located in | | no
loc | string | Sublocation of the server within the datacenter | | no
isOnline | boolean | Status of the server | true | yes
isDeleted | boolean | Deletion of the server | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | server | update | yes | no

# EXAMPLES

```
soma server update 43a7c39e-20ee-44cb-be46-72fae033911a \
     name example-server-a \
     assetid 43 \
     datacenter de.fra \
     location 'Row A, Rack 2, Unit 5' \
     deleted true
```
