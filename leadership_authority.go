package election

type LeadershipAuthority struct {
}

func (authority *LeadershipAuthority) Resignation() {
	consulClient.Session().Destroy(session, nil)
}
