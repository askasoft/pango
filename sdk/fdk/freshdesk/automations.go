package freshdesk

import "fmt"

// ---------------------------------------------------
// Automation

func (fd *Freshdesk) ListAutomationRules(automationTypeID int) ([]*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules", fd.Domain, automationTypeID)
	rules := []*AutomationRule{}
	_, err := fd.doList(url, nil, &rules)
	return rules, err
}

func (fd *Freshdesk) GetAutomationRule(automationTypeID int, rid int64) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	rule := &AutomationRule{}
	err := fd.doGet(url, rule)
	return rule, err
}

func (fd *Freshdesk) DeleteAutomationRule(automationTypeID int, rid int64) error {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateAutomationRule(automationTypeID int, rule *AutomationRule) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules", fd.Domain, automationTypeID)
	result := &AutomationRule{}
	err := fd.doPost(url, rule, result)
	return result, err
}

func (fd *Freshdesk) UpdateAutomationRule(automationTypeID int, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	result := &AutomationRule{}
	err := fd.doPut(url, rule, result)
	return result, err
}
