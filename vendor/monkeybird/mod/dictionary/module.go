// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package dictionary provides a custom dictionary.
// It allows users to define
package dictionary

import (
	"compress/gzip"
	"encoding/json"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type module struct {
	m        sync.RWMutex
	file     string
	commands *cmd.Set
	table    map[string]string
}

// New returns a new dictionary module.
func New() mod.Module {
	return &module{
		table: make(map[string]string),
	}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(
		prof.CommandPrefix(),
		func(r *irc.Request) bool {
			return prof.IsWhitelisted(r.SenderMask)
		},
	)

	m.commands.Bind(tr.DefineName, tr.DefineDesc, false, m.cmdDefine).
		Add(tr.DefineTermName, tr.DefineTermDesc, true, cmd.RegAny)

	m.commands.Bind(tr.AddDefineName, tr.AddDefineDesc, true, m.cmdAddDefine).
		Add(tr.AddDefineTermName, tr.AddDefineTermDesc, true, cmd.RegAny).
		Add(tr.AddDefineDefinitionName, tr.AddDefineDefinitionDesc, true, cmd.RegAny)

	m.commands.Bind(tr.RemoveDefineName, tr.RemoveDefineDesc, true, m.cmdRemoveDefine).
		Add(tr.RemoveDefineTermName, tr.RemoveDefineTermDesc, true, cmd.RegAny)

	m.file = filepath.Join(prof.Root(), "dictionary.dat")
	m.load()
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.commands.Clear()
	pb.Unbind("PRIVMSG", m.onPrivMsg)
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

// onPrivMsg ensures custom commands are executed.
func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// cmdAddDefine allows a user to add a new definition.
func (m *module) cmdAddDefine(w irc.ResponseWriter, r *cmd.Request) {
	m.m.Lock()
	defer m.m.Unlock()

	key := strings.ToLower(r.String(0))
	if _, ok := m.table[key]; ok {
		proto.PrivMsg(w, r.SenderName, tr.AddDefineAllreadyUsed, r.String(0))
		return
	}

	m.table[key] = r.Remainder(2)
	m.save()

	proto.PrivMsg(w, r.SenderName, tr.AddDefineDisplayText, r.String(0))
}

// cmdRemoveDefine allows a user to remove an existing definition.
func (m *module) cmdRemoveDefine(w irc.ResponseWriter, r *cmd.Request) {
	m.m.Lock()
	defer m.m.Unlock()

	key := strings.ToLower(r.String(0))
	if _, ok := m.table[key]; !ok {
		proto.PrivMsg(w, r.SenderName, tr.RemoveDefineNotFound, r.String(0))
		return
	}

	delete(m.table, key)
	m.save()

	proto.PrivMsg(w, r.SenderName, tr.RemoveDefineDisplayText, r.String(0))
}

// cmdDefine yields the definition of a given term, if found.
func (m *module) cmdDefine(w irc.ResponseWriter, r *cmd.Request) {
	m.m.RLock()
	defer m.m.RUnlock()

	key := strings.ToLower(r.String(0))
	v, ok := m.table[key]
	if !ok {
		proto.PrivMsg(w, r.Target, tr.DefineNotFound, r.SenderName, r.String(0))
		return
	}

	proto.PrivMsg(w, r.Target, tr.DefineDisplayText, r.SenderName, v)
}

// load reads dictionary data from a file.
func (m *module) load() error {
	fd, err := os.Open(m.file)
	if err != nil {
		return err
	}

	defer fd.Close()

	gz, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}

	defer gz.Close()

	return json.NewDecoder(gz).Decode(&m.table)
}

// save writes dictionary data to a file.
func (m *module) save() error {
	fd, err := os.Create(m.file)
	if err != nil {
		return err
	}

	defer fd.Close()

	gz := gzip.NewWriter(fd)
	defer gz.Close()

	return json.NewEncoder(gz).Encode(m.table)
}
