package freshdesk

import "context"

// ---------------------------------------------------
// Automation

type ListAutomationRulesOption = PageOption

func (fd *Freshdesk) ListAutomationRules(ctx context.Context, aType AutomationType, laro *ListAutomationRulesOption) ([]*AutomationRule, bool, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	rules := []*AutomationRule{}
	next, err := fd.doList(ctx, url, laro, &rules)
	return rules, next, err
}

func (fd *Freshdesk) IterAutomationRules(ctx context.Context, aType AutomationType, laro *ListAutomationRulesOption, iarf func(*AutomationRule) error) error {
	if laro == nil {
		laro = &ListAutomationRulesOption{}
	}
	if laro.Page < 1 {
		laro.Page = 1
	}
	if laro.PerPage < 1 {
		laro.PerPage = 100
	}

	for {
		ars, next, err := fd.ListAutomationRules(ctx, aType, laro)
		if err != nil {
			return err
		}
		for _, ar := range ars {
			if err = iarf(ar); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		laro.Page++
	}
	return nil
}

func (fd *Freshdesk) GetAutomationRule(ctx context.Context, aType AutomationType, rid int64) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	rule := &AutomationRule{}
	err := fd.doGet(ctx, url, rule)
	return rule, err
}

func (fd *Freshdesk) DeleteAutomationRule(ctx context.Context, aType AutomationType, rid int64) error {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	return fd.doDelete(ctx, url)
}

func (fd *Freshdesk) CreateAutomationRule(ctx context.Context, aType AutomationType, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules", aType)
	result := &AutomationRule{}
	if err := fd.doPost(ctx, url, rule, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateAutomationRule(ctx context.Context, aType AutomationType, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fd.endpoint("/automations/%d/rules/%d", aType, rid)
	result := &AutomationRule{}
	if err := fd.doPut(ctx, url, rule, result); err != nil {
		return nil, err
	}
	return result, nil
}
