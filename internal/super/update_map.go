package super // import "github.com/mjolnir42/soma/internal/super"

import "github.com/mjolnir42/soma/internal/msg"

func (s *Supervisor) updateMap(q *msg.Request) {

	switch q.Super.Object {
	case `team`:
		switch q.Action {
		case `add`:
			s.mapTeamID.insert(q.Super.Team.Id, q.Super.Team.Name)
		case `update`:
			s.mapTeamID.insert(q.Super.Team.Id, q.Super.Team.Name)
		case `delete`:
			s.mapTeamID.remove(q.Super.Team.Id)
		}
	case `user`:
		switch q.Action {
		case `add`:
			s.mapUserID.insert(q.Super.User.Id, q.Super.User.UserName)
			s.mapUserIDReverse.insert(q.Super.User.UserName, q.Super.User.Id)
			s.mapUserTeamID.insert(q.Super.User.Id, q.Super.User.TeamId)
		case `update`:
			oldname, _ := s.mapUserID.get(q.Super.User.Id)
			if oldname != q.Super.User.UserName {
				s.mapUserIDReverse.remove(oldname)
			}
			s.mapUserID.insert(q.Super.User.Id, q.Super.User.UserName)
			s.mapUserIDReverse.insert(q.Super.User.UserName, q.Super.User.Id)
			s.mapUserTeamID.insert(q.Super.User.Id, q.Super.User.TeamId)
		case `delete`:
			if name, ok := s.mapUserID.get(q.Super.User.Id); ok {
				s.mapUserIDReverse.remove(name)
			}
			s.mapUserID.remove(q.Super.User.Id)
			s.mapUserTeamID.remove(q.Super.User.Id)
		}
	}

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
