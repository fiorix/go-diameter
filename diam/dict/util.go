// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser, helper functions.  Part of go-diameter.

package dict

import "fmt"

// Apps return a list of all applications loaded in the Dict instance.
func (p Parser) Apps() []*App {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var apps []*App
	for _, f := range p.file {
		for _, app := range f.App {
			apps = append(apps, app)
		}
	}
	return apps
}

// FindAVP is a helper function that returns a pre-loaded AVP from the Dict.
// If the AVP code is not found in the given appid it tries with appid=0
// before returning an error.
// @code can be either the AVP code (int, uint32) or name (string).
func (p Parser) FindAVP(appid uint32, code interface{}) (*AVP, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var (
		avp *AVP
		ok  bool
		err error
	)
	switch code.(type) {
	case string:
		avp, ok = p.avpname[nameIdx{appid, code.(string)}]
		if !ok {
			err = fmt.Errorf("Could not find AVP %s", code.(string))
		}
	case uint32:
		avp, ok = p.avpcode[codeIdx{appid, code.(uint32)}]
		if !ok {
			err = fmt.Errorf("Could not find AVP %d", code.(uint32))
		}
	case int:
		return p.FindAVP(appid, uint32(code.(int)))
	default:
		return nil, fmt.Errorf("Unsupported AVP code type %#v", code)
	}
	if ok {
		return avp, nil
	} else if appid != 0 {
		return p.FindAVP(0, code)
	}
	return nil, err
}

// ScanAVP is a helper function that returns a pre-loaded AVP from the Dict.
// It's similar to FindAPI except that it scans the list of available AVPs
// instead of looking into one specific appid.
// @code can be either the AVP code (uint32) or name (string).
func (p Parser) ScanAVP(code interface{}) (*AVP, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	switch code.(type) {
	case string:
		for idx, avp := range p.avpname {
			if idx.Name == code.(string) {
				return avp, nil
			}
		}
		return nil, fmt.Errorf("Could not find AVP %s", code.(string))
	case uint32:
		for idx, avp := range p.avpcode {
			if idx.Code == code.(uint32) {
				return avp, nil
			}
		}
		return nil, fmt.Errorf("Could not find AVP code %d", code.(uint32))
	case int:
		return p.ScanAVP(uint32(code.(int)))
	}
	return nil, fmt.Errorf("Unsupported AVP code type %#v", code)
}

// FindCmd is a helper function that returns a pre-loaded Cmd from the Dict.
func (p Parser) FindCmd(appid, code uint32) (*Cmd, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if cmd, ok := p.cmd[codeIdx{appid, code}]; ok {
		return cmd, nil
	} else if cmd, ok = p.cmd[codeIdx{0, code}]; ok {
		// Always fall back to base dict.
		return cmd, nil
	}
	return nil, fmt.Errorf("Could not find preloaded Cmd with code %d", code)
}

// Enum is a helper function that returns a pre-loaded Enum item for the
// given AVP appid, code and n. (n is the enum code in the dictionary)
func (p *Parser) Enum(appid, code uint32, n uint8) (*Enum, error) {
	avp, err := p.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	if avp.Data.Type != "Enumerated" {
		return nil, fmt.Errorf(
			"Data of AVP %s (%d) data is not Enumerated.",
			avp.Name, avp.Code)
	}
	for _, item := range avp.Data.Enum {
		if item.Code == n {
			return item, nil
		}
	}
	return nil, fmt.Errorf(
		"Could not find preload Enum %d for AVP %s (%d)",
		n, avp.Name, avp.Code)
}

// Rule is a helper function that returns a pre-loaded Rule item for the
// given AVP code and name.
func (p *Parser) Rule(appid, code uint32, n string) (*Rule, error) {
	avp, err := p.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	if avp.Data.Type != "Grouped" {
		return nil, fmt.Errorf(
			"Data of AVP %s (%d) data is not Grouped.",
			avp.Name, avp.Code)
	}
	for _, item := range avp.Data.Rule {
		if item.AVP == n {
			return item, nil
		}
	}
	return nil, fmt.Errorf(
		"Could not find preload Rule for %s for AVP %s (%d)",
		n, avp.Name, avp.Code)
}

// PrettyPrint prints the Dict in a human readable form.
func (p *Parser) PrettyPrint() {
	for _, f := range p.file {
		for _, app := range f.App {
			fmt.Printf("Application Id: %d\n", app.Id)
			fmt.Printf("  Vendors:\n")
			for _, vendor := range app.Vendor {
				fmt.Printf("    Id=%d Name=%s\n", vendor.Id, vendor.Name)
			}
			fmt.Printf("  Commands:\n")
			for _, cmd := range app.Cmd {
				printCmd(cmd)
			}
			fmt.Printf("  AVPs:\n")
			for _, avp := range app.AVP {
				printAVP(avp)
			}
		}
	}
}

func printCmd(cmd *Cmd) {
	fmt.Printf("    % -4d %s\n", cmd.Code, cmd.Name)
	fmt.Printf("      %sR:\n", cmd.Short)
	for _, rule := range cmd.Request.Rule {
		if rule.Required && rule.Min == 0 {
			rule.Min = 1
		}
		fmt.Printf("        % -40s required=% -5t min=%d max=%d\n",
			rule.AVP, rule.Required, rule.Min, rule.Max)
	}
	fmt.Printf("      %sA:\n", cmd.Short)
	for _, rule := range cmd.Answer.Rule {
		if rule.Required && rule.Min == 0 {
			rule.Min = 1
		}
		fmt.Printf("        % -40s required=% -5t min=%d max=%d\n",
			rule.AVP, rule.Required, rule.Min, rule.Max)
	}
}

func printAVP(avp *AVP) {
	fmt.Printf("   % -4d %s: %s\n",
		avp.Code, avp.Name, avp.Data.Type)
	// Enumerated
	if len(avp.Data.Enum) > 0 {
		fmt.Printf("    Items:\n")
		for _, item := range avp.Data.Enum {
			fmt.Printf("      % -2d %s\n", item.Code, item.Name)
		}
	}
	// Grouped AVPs
	if len(avp.Data.Rule) > 0 {
		fmt.Printf("    Rules:\n")
		for _, rule := range avp.Data.Rule {
			fmt.Printf("      % -40s required=% -5t min=%d max=%d\n",
				rule.AVP, rule.Required, rule.Min, rule.Max)
		}
	}
}
