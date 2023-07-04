package freshdesk

// ---------------------------------------------------
// Automation

func (fd *Freshdesk) ListAutomationRules(aType AutomationType) ([]*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	rules := []*AutomationRule{}
	_, err := fd.doList(url, nil, &rules)
	return rules, err
}

func (fd *Freshdesk) GetAutomationRule(aType AutomationType, rid int64) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	rule := &AutomationRule{}
	err := fd.doGet(url, rule)
	return rule, err
}

func (fd *Freshdesk) DeleteAutomationRule(aType AutomationType, rid int64) error {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateAutomationRule(aType AutomationType, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	result := &AutomationRule{}
	err := fd.doPost(url, rule, result)
	return result, err
}

func (fd *Freshdesk) UpdateAutomationRule(aType AutomationType, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	result := &AutomationRule{}
	err := fd.doPut(url, rule, result)
	return result, err
}
