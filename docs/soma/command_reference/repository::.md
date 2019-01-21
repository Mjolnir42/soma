# repository management

Repository management contains the actions for managing configuration repositories.

# SYNOPSIS OVERVIEW

```
soma repository create ${repository} team ${team}
soma repository destroy ${repository} [from ${team}]
soma repository rename ${repository} to ${newName} [from ${team}]
soma repository repossess ${repository} to ${newTeam} [from ${team}]
soma repository list
soma repository show ${repository} [from ${team}]
soma repository search [id ${uuid}] [name ${repository}] [team ${team}] [deleted ${isDeleted}] [active ${isActive}]
soma repository dumptree ${repository}
soma repository property create system ${system} on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create custom ${custom} on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create service ${service} on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create oncall ${oncall} on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property destroy system ${system} on ${repository} view ${view}
soma repository property destroy custom ${custom} on ${repository} view ${view}
soma repository property destroy service ${service} on ${repository} view ${view}
soma repository property destroy oncall ${oncall} on ${repository} view ${view}
```

See `soma repository help ${command}` or `soma repository property help ${command}` for detailed help.
