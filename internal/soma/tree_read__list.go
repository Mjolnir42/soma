/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

// bucketsInRepository returns the list of buckets in a repository
func (r *TreeRead) bucketsInRepository(id string) ([]string, error) {
	rows, err := r.stmtListRepositoryMemberBuckets.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		bID := ``
		if err = rows.Scan(&bID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, bID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// groupsInBucket returns the list of groups in a bucket
func (r *TreeRead) groupsInBucket(id string) ([]string, error) {
	rows, err := r.stmtListBucketMemberGroups.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		gID := ``
		if err = rows.Scan(&gID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, gID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// clustersInBucket returns the list of clusters in a bucket
func (r *TreeRead) clustersInBucket(id string) ([]string, error) {
	rows, err := r.stmtListBucketMemberClusters.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		cID := ``
		if err = rows.Scan(&cID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, cID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// nodesInBucket returns the list of nodes in a bucket
func (r *TreeRead) nodesInBucket(id string) ([]string, error) {
	rows, err := r.stmtListBucketMemberNodes.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		nID := ``
		if err = rows.Scan(&nID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, nID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// groupsInGroup returns the list of groups in a group
func (r *TreeRead) groupsInGroup(id string) ([]string, error) {
	rows, err := r.stmtListGroupMemberGroups.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		gID := ``
		if err = rows.Scan(&gID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, gID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// clustersInGroup returns the list of clusters in a group
func (r *TreeRead) clustersInGroup(id string) ([]string, error) {
	rows, err := r.stmtListGroupMemberClusters.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		cID := ``
		if err = rows.Scan(&cID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, cID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// nodesInGroup returns the list of nodes in a group
func (r *TreeRead) nodesInGroup(id string) ([]string, error) {
	rows, err := r.stmtListGroupMemberNodes.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		nID := ``
		if err = rows.Scan(&nID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, nID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// nodesInCluster returns the list of nodes in a cluster
func (r *TreeRead) nodesInCluster(id string) ([]string, error) {
	rows, err := r.stmtListClusterMemberNodes.Query(id)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		nID := ``
		if err = rows.Scan(&nID); err != nil {
			rows.Close()
			return []string{}, err
		}
		res = append(res, nID)
	}
	if err = rows.Err(); err != nil {
		return []string{}, err
	}
	return res, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
