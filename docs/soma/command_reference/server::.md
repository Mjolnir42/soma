# server management

Server management are the functions for maintaining the physical server
resources within SOMA that everything is running on.

# SYNOPSIS OVERVIEW

```
soma server add ${name} assetid ${assetID} datacenter ${locode} location ${loc} [online ${isOnline}]
soma server update ${serverID} name ${name} assetid ${assetID} datacenter ${locode} location ${loc} [online ${isOnline}] [deleted ${isDeleted}]
soma server remove ${name}
soma server purge ${name}
soma server show ${name}
soma server list
soma server sync
soma server null datacenter ${locode}
```

See `soma server help ${command}` for detailed help.
