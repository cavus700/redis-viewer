// @description: commands
// @file: command.go
// @date: 2022/02/07

package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/saltfishpr/redis-viewer/internal/config"
	"github.com/saltfishpr/redis-viewer/internal/rv"
	"github.com/saltfishpr/redis-viewer/internal/util"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct {
	err error
}

type scanMsg struct {
	items []list.Item
}

func (m model) scanCmd() tea.Cmd {
	return func() tea.Msg {
		cfg := config.GetConfig()
		ctx := context.Background()
		var (
			val   interface{}
			err   error
			items []list.Item
		)

		keys := rv.GetKeys(m.rdb, 0, m.searchValue, cfg.Count)
		for key := range keys {
			kt := m.rdb.Type(ctx, key).Val()
			switch kt {
			case "string":
				val, err = m.rdb.Get(ctx, key).Result()
			case "list":
				val, err = m.rdb.LRange(ctx, key, 0, -1).Result()
			case "set":
				val, err = m.rdb.SMembers(ctx, key).Result()
			case "zset":
				val, err = m.rdb.ZRange(ctx, key, 0, -1).Result()
			case "hash":
				val, err = m.rdb.HGetAll(ctx, key).Result()
			default:
				val = ""
				err = fmt.Errorf("unsupported type: %s", kt)
			}
			if err != nil {
				items = append(items, item{keyType: kt, key: key, val: err.Error(), err: true})
			} else {
				valBts, _ := util.JsonMarshalIndent(val)
				items = append(items, item{keyType: kt, key: key, val: string(valBts)})
			}
		}

		return scanMsg{items: items}
	}
}

type countMsg struct {
	count int
}

func (m model) countCmd() tea.Cmd {
	return func() tea.Msg {
		count, err := rv.CountKeys(m.rdb, m.searchValue)
		if err != nil {
			return errMsg{err: err}
		}

		return countMsg{count: count}
	}
}

type tickMsg struct {
	t string
}

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tickMsg{t: time.Now().Format("2006-01-02 15:04:05")}
	})
}
