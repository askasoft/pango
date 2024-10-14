package freshdesk

// ---------------------------------------------------
// Product

type ListProductsOption = PageOption

func (fd *Freshdesk) GetProduct(id int64) (*Product, error) {
	url := fd.endpoint("/products/%d", id)
	product := &Product{}
	err := fd.doGet(url, product)
	return product, err
}

func (fd *Freshdesk) ListProducts(lpo *ListProductsOption) ([]*Product, bool, error) {
	url := fd.endpoint("/products")
	products := []*Product{}
	next, err := fd.doList(url, lpo, &products)
	return products, next, err
}

func (fd *Freshdesk) IterProducts(lpo *ListProductsOption, ipf func(*Product) error) error {
	if lpo == nil {
		lpo = &ListProductsOption{}
	}
	if lpo.Page < 1 {
		lpo.Page = 1
	}
	if lpo.PerPage < 1 {
		lpo.PerPage = 100
	}

	for {
		ps, next, err := fd.ListProducts(lpo)
		if err != nil {
			return err
		}
		for _, ag := range ps {
			if err = ipf(ag); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lpo.Page++
	}
	return nil
}
