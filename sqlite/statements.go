package sqlite

import "database/sql"

// Statements - generic sqlite prepared Statements
type Statements struct {
	Node         *sql.Stmt
	NodeTags     *sql.Stmt
	Way          *sql.Stmt
	WayTags      *sql.Stmt
	WayNodes     *sql.Stmt
	Relation     *sql.Stmt
	RelationTags *sql.Stmt
	Member       *sql.Stmt
}

// Close - close all prepared Statements
func (s *Statements) Close() {
	s.Node.Close()
	s.NodeTags.Close()
	s.Way.Close()
	s.WayTags.Close()
	s.WayNodes.Close()
	s.Relation.Close()
	s.RelationTags.Close()
	s.Member.Close()
}
