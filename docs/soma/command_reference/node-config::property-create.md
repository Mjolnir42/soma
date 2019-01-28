# node management

```
soma node add ${node} assetid ${id} team ${team} [server ${server}] [online ${isOnline}]
soma node remove ${node}
soma node update ${nodeUUID} name ${name} assetid ${id} team ${team} server ${server} online ${isOnline} deleted ${isDeleted}
soma node repossess ${node} to ${team}
soma node rename ${node} to ${name}
soma node relocate ${node} to ${server}
soma node list
soma node show ${node}
soma node sync
soma node config ${node}
soma node assign ${node} to ${bucket}
soma node unassign ${node} [from ${bucket}]
soma node dumptree ${node} [in ${bucket}]
soma node property create system  ${system}  on ${node} [in ${bucket}] view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma node property create custom  ${custom}  on ${node} [in ${bucket}] view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma node property create service ${service} on ${node} [in ${bucket}] view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma node property create oncall  ${oncall}  on ${node} [in ${bucket}] view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma node property destroy system  ${system}  on ${node} [in ${bucket}] view ${view}
soma node property destroy custom  ${custom}  on ${node} [in ${bucket}] view ${view}
soma node property destroy service ${service} on ${node} [in ${bucket}] view ${view}
soma node property destroy oncall  ${oncall}  on ${node} [in ${bucket}] view ${view}
```

See `soma node help ${command}` for detailed help.
