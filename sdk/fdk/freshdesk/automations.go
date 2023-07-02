package freshdesk

// ---------------------------------------------------
// Automation

func (fd *Freshdesk) ListAutomationRules(automationTypeID int) ([]*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", automationTypeID)
	rules := []*AutomationRule{}
	_, err := fd.doList(url, nil, &rules)
	return rules, err
}

func (fd *Freshdesk) GetAutomationRule(automationTypeID int, rid int64) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", automationTypeID, rid)
	rule := &AutomationRule{}
	err := fd.doGet(url, rule)
	return rule, err
}

func (fd *Freshdesk) DeleteAutomationRule(automationTypeID int, rid int64) error {
	url := fd.endpoint("/automations/%d/rules/%d", automationTypeID, rid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateAutomationRule(automationTypeID int, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", automationTypeID)
	result := &AutomationRule{}
	err := fd.doPost(url, rule, result)
	return result, err
}

func (fd *Freshdesk) UpdateAutomationRule(automationTypeID int, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", automationTypeID, rid)
	result := &AutomationRule{}
	err := fd.doPut(url, rule, result)
	return result, err
}
