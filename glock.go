package glock

import (
	zk "github.com/burke/gozk"
	"path"
	"sort"
)

type Glock struct {
	conn   *zk.Conn
	root   string
	myPath string
}

func New(conn *zk.Conn, root string) *Glock {
	return &Glock{conn: conn, root: root}
}

func (g *Glock) Lock() (err error) {
	g.myPath, err = g.conn.Create(g.root+"/lock", "", zk.EPHEMERAL|zk.SEQUENCE, zk.WorldACL(zk.PERM_ALL))
	if err != nil {
		return err
	}

	var (
		children []string
		w        <-chan zk.Event
	)
	for {
		children, _, w, err = g.conn.ChildrenW(g.root)
		sort.Strings(children)
		if children[0] == path.Base(g.myPath) {
			return nil
		}
		<-w
	}

	return nil
}

func (g *Glock) Unlock() {
	g.conn.Delete(g.myPath, -1)
}
