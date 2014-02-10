package parser

import (
	"common"
	"time"
)

type QuerySpec struct {
	query                  *Query
	database               string
	isRegex                bool
	names                  []string
	user                   common.User
	startTime              time.Time
	endTime                time.Time
	seriesValuesAndColumns map[*Value][]string
}

func NewQuerySpec(user common.User, database string, query *Query) *QuerySpec {
	return &QuerySpec{user: user, query: query, database: database}
}

func (self *QuerySpec) GetStartTime() time.Time {
	if self.query.SelectQuery != nil {
		return self.query.SelectQuery.GetStartTime()
	} else if self.query.DeleteQuery != nil {
		return self.query.DeleteQuery.GetStartTime()
	}
	return time.Now()
}

func (self *QuerySpec) GetEndTime() time.Time {
	if self.query.SelectQuery != nil {
		return self.query.SelectQuery.GetEndTime()
	} else if self.query.DeleteQuery != nil {
		return self.query.DeleteQuery.GetEndTime()
	}
	return time.Now()
}

func (self *QuerySpec) Database() string {
	return self.database
}

func (self *QuerySpec) User() common.User {
	return self.user
}

func (self *QuerySpec) TableNames() []string {
	if self.names != nil {
		return self.names
	}
	if self.query.SelectQuery == nil {
		self.names = []string{}
		return self.names
	}

	namesAndColumns := self.query.SelectQuery.GetReferencedColumns()

	names := make([]string, 0, len(namesAndColumns))
	for name, _ := range namesAndColumns {
		if _, isRegex := name.GetCompiledRegex(); isRegex {
			self.isRegex = true
		} else {
			names = append(names, name.Name)
		}
	}
	return names
}

func (self *QuerySpec) SeriesValuesAndColumns() map[*Value][]string {
	if self.seriesValuesAndColumns != nil {
		return self.seriesValuesAndColumns
	}
	if self.query.SelectQuery != nil {
		self.seriesValuesAndColumns = self.query.SelectQuery.GetReferencedColumns()
	} else {
		self.seriesValuesAndColumns = make(map[*Value][]string)
	}
	return self.seriesValuesAndColumns
}

func (self *QuerySpec) GetGroupByInterval() *time.Duration {
	if self.query.SelectQuery == nil {
		return nil
	}
	duration, _ := self.query.SelectQuery.GetGroupByClause().GetGroupByTime()
	return duration
}

func (self *QuerySpec) IsRegex() bool {
	self.TableNames()
	return self.isRegex
}

func (self *QuerySpec) HasReadAccess(name string) bool {
	return self.user.HasReadAccess(name)
}

func (self *QuerySpec) IsSinglePointQuery() bool {
	if self.query.SelectQuery != nil {
		return self.query.SelectQuery.IsSinglePointQuery()
	}
	return false
}

func (self *QuerySpec) SelectQuery() *SelectQuery {
	return self.query.SelectQuery
}