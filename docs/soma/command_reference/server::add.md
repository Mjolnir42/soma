# DESCRIPTION

This command is used to add physical servers to SOMA. Their main use is
to anchor the nodes that reference them to a physical datacenter location.

Server names must not be formatted as UUIDs.

# SYNOPSIS

```
soma server add ${name} assetid ${assetID} datacenter ${locode} location ${loc} [online ${isOnline}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the server | | no
assetID | integer | Numeric asset ID of the server | | no
locode | string | UN/Locode of the datacenter the server is located in | | no
loc | string | Sublocation of the server within the datacenter | | no
isOnline | boolean | Status of the server | true | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | server | add | yes | no

# EXAMPLES

```
soma server add example-server-a \
    assetid 42 \
    datacenter de.fra \
    location 'Row A, Rack 2, Unit 5'
soma server add example-server-b \
    assetid 23 \
    datacenter de.fra \
    location 'Row A, Rack 2, Unit 6'
    online false
```
