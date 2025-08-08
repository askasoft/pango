package freshdesk

import "context"

// ---------------------------------------------------
// Product

type ListProductsOption = PageOption

func (fd *Freshdesk) GetProduct(ctx context.Context, id int64) (*Product, error) {
	url := fd.Endpoint("/products/%d", id)
	product := &Product{}
	err := fd.DoGet(ctx, url, product)
	return product, err
}

func (fd *Freshdesk) ListProducts(ctx context.Context, lpo *ListProductsOption) ([]*Product, bool, error) {
	url := fd.Endpoint("/products")
	products := []*Product{}
	next, err := fd.DoList(ctx, url, lpo, &products)
	return products, next, err
}

func (fd *Freshdesk) IterProducts(ctx context.Context, lpo *ListProductsOption, ipf func(*Product) error) error {
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
		ps, next, err := fd.ListProducts(ctx, lpo)
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
